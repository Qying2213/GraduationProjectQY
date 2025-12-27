package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"evaluator-service/internal/config"
	"evaluator-service/internal/logging"
	"evaluator-service/internal/models"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/xuri/excelize/v2"
)

type ExportService struct {
	cfg *config.Config
	log *logging.Logger
}

func NewExportService(cfg *config.Config, log *logging.Logger) *ExportService {
	return &ExportService{cfg: cfg, log: log}
}

// PDFfromHTML renders HTML into PDF via headless Chrome (chromedp).
func (e *ExportService) PDFfromHTML(html string) ([]byte, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	var pdf []byte
	tasks := chromedp.Tasks{
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx2 context.Context) error {
			script := `document.open();document.write(` + "`" + html + "`" + `);document.close();`
			return chromedp.Evaluate(script, nil).Do(ctx2)
		}),
		chromedp.Sleep(200 * time.Millisecond),
		chromedp.ActionFunc(func(ctx2 context.Context) error {
			buf, _, err := page.PrintToPDF().WithPrintBackground(true).Do(ctx2)
			if err != nil {
				return err
			}
			pdf = buf
			return nil
		}),
	}
	if err := chromedp.Run(ctx, tasks); err != nil {
		return nil, err
	}
	return pdf, nil
}

func (e *ExportService) ExcelFromBatchJSON(jsonStr string) ([]byte, error) {
	var results []map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &results); err != nil {
		return nil, err
	}
	f := excelize.NewFile()
	sh := f.GetSheetName(0)
	headers := []string{"排名", "姓名", "文件名", "总分", "评级", "JD匹配", "年龄", "经验", "学历", "公司", "技术", "项目", "推荐结果"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		_ = f.SetCellValue(sh, cell, h)
	}
	for rIdx, r := range results {
		row := []interface{}{
			valOr(r["rank"], rIdx+1),
			valOr(r["candidate_name"], ""),
			valOr(r["filename"], ""),
			valOr(r["total_score"], 0),
			valOr(r["grade"], ""),
			valOr(r["jd_match"], 0),
			valOr(r["age_score"], 0),
			valOr(r["experience_score"], 0),
			valOr(r["education_score"], 0),
			valOr(r["company_score"], 0),
			valOr(r["tech_score"], 0),
			valOr(r["project_score"], 0),
			valOr(r["recommendation"], ""),
		}
		for cIdx, v := range row {
			cell, _ := excelize.CoordinatesToCellName(cIdx+1, rIdx+2)
			_ = f.SetCellValue(sh, cell, v)
		}
	}
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (e *ExportService) ExcelCompare(cands []models.Candidate) ([]byte, error) {
	f := excelize.NewFile()
	sh := f.GetSheetName(0)
	head := []string{"维度", "满分"}
	for _, c := range cands {
		head = append(head, c.Name)
	}
	for i, h := range head {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		_ = f.SetCellValue(sh, cell, h)
	}
	dims := []struct {
		label string
		max   int
		get   func(models.Candidate) interface{}
	}{
		{"总分", 100, func(c models.Candidate) interface{} { return c.TotalScore }},
		{"JD匹配", 100, func(c models.Candidate) interface{} { return c.JDMatch }},
		{"年龄", 10, func(c models.Candidate) interface{} { return c.AgeScore }},
		{"经验", 25, func(c models.Candidate) interface{} { return c.ExperienceScore }},
		{"学历", 20, func(c models.Candidate) interface{} { return c.EducationScore }},
		{"公司", 15, func(c models.Candidate) interface{} { return c.CompanyScore }},
		{"技术", 25, func(c models.Candidate) interface{} { return c.TechScore }},
		{"项目", 15, func(c models.Candidate) interface{} { return c.ProjectScore }},
	}
	for r, d := range dims {
		cell, _ := excelize.CoordinatesToCellName(1, r+2)
		_ = f.SetCellValue(sh, cell, d.label)
		cell, _ = excelize.CoordinatesToCellName(2, r+2)
		_ = f.SetCellValue(sh, cell, d.max)
		for i, c := range cands {
			cell, _ := excelize.CoordinatesToCellName(i+3, r+2)
			_ = f.SetCellValue(sh, cell, d.get(c))
		}
	}
	cell, _ := excelize.CoordinatesToCellName(1, len(dims)+3)
	_ = f.SetCellValue(sh, cell, "评级")
	cell, _ = excelize.CoordinatesToCellName(2, len(dims)+3)
	_ = f.SetCellValue(sh, cell, "-")
	cell, _ = excelize.CoordinatesToCellName(1, len(dims)+4)
	_ = f.SetCellValue(sh, cell, "推荐结果")
	cell, _ = excelize.CoordinatesToCellName(2, len(dims)+4)
	_ = f.SetCellValue(sh, cell, "-")
	for i, c := range cands {
		cell, _ := excelize.CoordinatesToCellName(i+3, len(dims)+3)
		_ = f.SetCellValue(sh, cell, c.Grade)
		cell, _ = excelize.CoordinatesToCellName(i+3, len(dims)+4)
		_ = f.SetCellValue(sh, cell, c.Recommendation)
	}
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (e *ExportService) HTMLPageFromMarkdownBody(body string) string {
	return fmt.Sprintf(`<!DOCTYPE html><html><head><meta charset="utf-8"><style>body{font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial; padding:40px; line-height:1.6}</style></head><body>%s</body></html>`, body)
}

func valOr(v interface{}, def interface{}) interface{} {
	if v == nil {
		return def
	}
	switch x := v.(type) {
	case json.Number:
		if i, err := x.Int64(); err == nil {
			return i
		}
		if f, err := x.Float64(); err == nil {
			return f
		}
		return def
	case *json.Number:
		if x == nil {
			return def
		}
		if i, err := x.Int64(); err == nil {
			return i
		}
		if f, err := x.Float64(); err == nil {
			return f
		}
		return def
	case string, float64, int, int64, bool:
		return x
	case fmt.Stringer:
		return x.String()
	case io.Reader:
		return def
	default:
		return def
	}
}
