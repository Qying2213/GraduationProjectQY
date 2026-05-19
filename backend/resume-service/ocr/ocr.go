// Package ocr 提供简历评估前置的文本提取能力。
//
// 在毕业设计 AI 链路中，它是第一步：把 PDF、图片和 Word 简历转换成统一文本，
// 让后续 Embedding、RAG 和 Coze 评估都能基于同一种输入工作。这里采用多策略设计，
// 是因为真实简历可能是原生 PDF、扫描版 PDF、图片或 doc/docx 文件。
package ocr

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"

	"github.com/ledongthuc/pdf"
	"golang.org/x/text/unicode/norm"
)

// OCRResult 是 AI 解析和评估使用的统一文本提取结果。
// Confidence 是启发式置信度，主要用于前端解释文本提取质量，不代表严格的模型概率。
type OCRResult struct {
	Text       string  `json:"text"`
	Confidence float64 `json:"confidence"`
	Pages      int     `json:"pages"`
}

type textCandidate struct {
	Method     string
	Text       string
	Score      float64
	Confidence float64
}

var (
	suspiciousASCIIBlockPattern = regexp.MustCompile(`[A-Za-z0-9_./+=:-]{16,}`)
	pageCountPattern            = regexp.MustCompile(`Pages:\s+(\d+)`)

	tesseractLangOnce sync.Once
	tesseractLang     string
)

// ExtractTextFromPDF 从 PDF 简历中提取可用文本。
//
// 这里会依次尝试多种候选方案：
//  1. pdftotext 处理原生 PDF；
//  2. Go PDF 解析作为跨平台降级方案；
//  3. Tesseract OCR 处理扫描版或噪声较高的 PDF。
//
// 最终通过文本质量评分选择最佳结果。真实简历经常带水印、表格或扫描页，因此这种
// 多策略处理对答辩演示很关键。
func ExtractTextFromPDF(filePath string) (*OCRResult, error) {
	numPages := countPDFPages(filePath)
	candidates := make([]textCandidate, 0, 3)

	if text, err := extractTextWithPDFToText(filePath); err == nil {
		if candidate := buildTextCandidate("pdftotext", text, 0.94); candidate.Text != "" {
			candidates = append(candidates, candidate)
		}
	}

	if text, pages, err := extractTextWithGoPDF(filePath); err == nil {
		if numPages == 0 {
			numPages = pages
		}
		if candidate := buildTextCandidate("go-pdf", text, 0.90); candidate.Text != "" {
			candidates = append(candidates, candidate)
		}
	}

	best := pickBestTextCandidate(candidates)

	// 文本很短或噪声比例高时，再尝试图片 OCR。这样原生 PDF 走快速路径，
	// 扫描件再付出较高的 OCR 成本。
	if shouldFallbackToOCR(best) {
		if text, err := ocrPDFWithTesseract(filePath); err == nil {
			if candidate := buildTextCandidate("tesseract", text, 0.82); candidate.Text != "" {
				if candidate.Score >= best.Score+0.03 || len([]rune(candidate.Text)) > len([]rune(best.Text))+80 {
					best = candidate
				}
			}
		}
	}

	if best.Text == "" {
		return nil, fmt.Errorf("无法从PDF提取可用文本")
	}

	return &OCRResult{
		Text:       best.Text,
		Confidence: best.Confidence,
		Pages:      numPages,
	}, nil
}

// ExtractTextFromImage 使用 Tesseract 从图片简历中识别文本。
// 扫描版 PDF 被转换为 PNG 页面后，也会复用这个函数逐页识别。
func ExtractTextFromImage(imagePath string) (*OCRResult, error) {
	// 检查tesseract是否可用
	if !isTesseractAvailable() {
		return nil, fmt.Errorf("tesseract未安装，请先安装: brew install tesseract tesseract-lang")
	}

	absPath, err := filepath.Abs(imagePath)
	if err == nil {
		imagePath = absPath
	}

	// 调用tesseract进行OCR
	args := []string{imagePath, "stdout"}
	if lang := preferredTesseractLanguage(); lang != "" {
		args = append(args, "-l", lang)
	}
	args = append(args, "--psm", "6", "-c", "preserve_interword_spaces=1")
	cmd := exec.Command("tesseract", args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("OCR识别失败: %w", err)
	}

	text := cleanExtractedText(string(output))
	score := textQualityScore(text)

	return &OCRResult{
		Text:       text,
		Confidence: clamp(0.70+score*0.22, 0.70, 0.92),
		Pages:      1,
	}, nil
}

// ocrPDFWithTesseract 先把 PDF 每页渲染成图片，再逐页 OCR。
// 它比直接文本提取慢，但能让扫描版简历进入后续 AI 链路。
func ocrPDFWithTesseract(pdfPath string) (string, error) {
	if !isTesseractAvailable() {
		return "", fmt.Errorf("tesseract未安装")
	}

	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "pdf_ocr_")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tempDir)

	if !isPdftoppmAvailable() {
		return "", fmt.Errorf("pdftoppm未安装，无法处理扫描版PDF")
	}

	// 提高渲染分辨率并转为灰度图，能显著改善带水印和小字号 PDF 的识别效果。
	imgPrefix := filepath.Join(tempDir, "page")
	cmd := exec.Command("pdftoppm", "-r", "300", "-gray", "-png", pdfPath, imgPrefix)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("PDF转图片失败: %w", err)
	}

	files, _ := filepath.Glob(imgPrefix + "*.png")
	sort.Strings(files)

	var parts []string
	for _, imgFile := range files {
		result, err := ExtractTextFromImage(imgFile)
		if err != nil {
			continue
		}
		if strings.TrimSpace(result.Text) != "" {
			parts = append(parts, result.Text)
		}
	}

	if len(parts) == 0 {
		return "", fmt.Errorf("OCR未提取到可用文本")
	}

	return strings.Join(parts, "\n\n"), nil
}

// isTesseractAvailable 检查tesseract是否可用
func isTesseractAvailable() bool {
	cmd := exec.Command("tesseract", "--version")
	return cmd.Run() == nil
}

// isPdftoppmAvailable 检查pdftoppm是否可用
func isPdftoppmAvailable() bool {
	cmd := exec.Command("pdftoppm", "-v")
	return cmd.Run() == nil
}

func isPdftotextAvailable() bool {
	cmd := exec.Command("pdftotext", "-v")
	return cmd.Run() == nil
}

func isPDFInfoAvailable() bool {
	cmd := exec.Command("pdfinfo", "-v")
	return cmd.Run() == nil
}

// ExtractTextFromFile 是 handler 层统一调用的文本提取入口。
// 文件类型分发集中在这里，AI 评估代码就不需要关心 PDF、图片和 Word 的实现细节。
func ExtractTextFromFile(filePath string) (*OCRResult, error) {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".pdf":
		return ExtractTextFromPDF(filePath)
	case ".png", ".jpg", ".jpeg", ".bmp", ".tiff", ".gif":
		return ExtractTextFromImage(filePath)
	case ".doc", ".docx":
		return extractTextFromWord(filePath)
	default:
		return nil, fmt.Errorf("不支持的文件格式: %s", ext)
	}
}

// extractTextFromWord 从 doc/docx 简历中提取文本。
// 很多候选人上传的是可编辑 Word 简历而不是 PDF，因此这是 OCR 之外的重要补充。
func extractTextFromWord(filePath string) (*OCRResult, error) {
	ext := strings.ToLower(filepath.Ext(filePath))

	if ext == ".docx" {
		return extractTextFromDocx(filePath)
	}

	// .doc格式需要使用antiword或其他工具
	if isAntiwordAvailable() {
		cmd := exec.Command("antiword", filePath)
		output, err := cmd.Output()
		if err != nil {
			return nil, fmt.Errorf("读取DOC文件失败: %w", err)
		}
		return &OCRResult{
			Text:       cleanExtractedText(string(output)),
			Confidence: 0.95,
			Pages:      1,
		}, nil
	}

	return nil, fmt.Errorf("无法处理.doc文件，请安装antiword: brew install antiword")
}

// extractTextFromDocx 从DOCX文件提取文本
func extractTextFromDocx(filePath string) (*OCRResult, error) {
	// 使用unzip解压docx并读取document.xml
	tempDir, err := os.MkdirTemp("", "docx_")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tempDir)

	cmd := exec.Command("unzip", "-q", filePath, "-d", tempDir)
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("解压DOCX失败: %w", err)
	}

	// 读取document.xml
	docPath := filepath.Join(tempDir, "word", "document.xml")
	content, err := os.ReadFile(docPath)
	if err != nil {
		return nil, fmt.Errorf("读取文档内容失败: %w", err)
	}

	// 简单提取文本（去除XML标签）
	text := extractTextFromXML(string(content))

	return &OCRResult{
		Text:       cleanExtractedText(text),
		Confidence: 0.95,
		Pages:      1,
	}, nil
}

// extractTextFromXML 从XML中提取纯文本
func extractTextFromXML(xml string) string {
	var result strings.Builder
	inTag := false
	inText := false

	for i := 0; i < len(xml); i++ {
		c := xml[i]
		if c == '<' {
			inTag = true
			// 检查是否是<w:t>标签
			if i+4 < len(xml) && xml[i:i+4] == "<w:t" {
				inText = true
			} else if i+5 < len(xml) && xml[i:i+5] == "</w:t" {
				inText = false
				result.WriteString(" ")
			} else if i+5 < len(xml) && xml[i:i+5] == "</w:p" {
				result.WriteString("\n")
			}
		} else if c == '>' {
			inTag = false
		} else if !inTag && inText {
			result.WriteByte(c)
		}
	}

	return result.String()
}

// isAntiwordAvailable 检查antiword是否可用
func isAntiwordAvailable() bool {
	cmd := exec.Command("antiword", "--version")
	return cmd.Run() == nil
}

// ReadFileContent 读取文件内容
func ReadFileContent(filePath string) ([]byte, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return io.ReadAll(f)
}

func extractTextWithGoPDF(filePath string) (string, int, error) {
	f, r, err := pdf.Open(filePath)
	if err != nil {
		return "", 0, fmt.Errorf("打开PDF文件失败: %w", err)
	}
	defer f.Close()

	var buf bytes.Buffer
	numPages := r.NumPage()

	for i := 1; i <= numPages; i++ {
		page := r.Page(i)
		if page.V.IsNull() {
			continue
		}

		text, err := page.GetPlainText(nil)
		if err != nil {
			continue
		}
		buf.WriteString(text)
		buf.WriteString("\n")
	}

	return buf.String(), numPages, nil
}

func extractTextWithPDFToText(filePath string) (string, error) {
	if !isPdftotextAvailable() {
		return "", fmt.Errorf("pdftotext未安装")
	}

	cmd := exec.Command("pdftotext", "-enc", "UTF-8", "-raw", "-nopgbrk", filePath, "-")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("pdftotext提取失败: %w", err)
	}

	return string(output), nil
}

func buildTextCandidate(method string, text string, baseConfidence float64) textCandidate {
	cleaned := cleanExtractedText(text)
	if cleaned == "" {
		return textCandidate{}
	}

	score := textQualityScore(cleaned)
	return textCandidate{
		Method:     method,
		Text:       cleaned,
		Score:      score,
		Confidence: clamp(baseConfidence*0.75+score*0.20, 0.70, 0.98),
	}
}

func pickBestTextCandidate(candidates []textCandidate) textCandidate {
	var best textCandidate
	for i, candidate := range candidates {
		if i == 0 {
			best = candidate
			continue
		}

		if candidate.Score > best.Score+0.03 {
			best = candidate
			continue
		}

		if math.Abs(candidate.Score-best.Score) <= 0.03 && len([]rune(candidate.Text)) > len([]rune(best.Text))+60 {
			best = candidate
		}
	}

	return best
}

func shouldFallbackToOCR(candidate textCandidate) bool {
	runeCount := len([]rune(candidate.Text))
	if runeCount == 0 {
		return true
	}
	if runeCount < 120 {
		return true
	}
	return candidate.Score < 0.58
}

func cleanExtractedText(text string) string {
	if strings.TrimSpace(text) == "" {
		return ""
	}

	text = strings.ReplaceAll(text, "\x00", "")
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	text = strings.ReplaceAll(text, "\f", "\n")
	text = strings.ReplaceAll(text, "\u00a0", " ")
	text = norm.NFKC.String(text)

	lines := strings.Split(text, "\n")
	normalizedLines := make([]string, 0, len(lines))
	lineCounts := make(map[string]int, len(lines))
	for _, line := range lines {
		normalized := normalizeWhitespace(line)
		normalizedLines = append(normalizedLines, normalized)
		if normalized != "" {
			lineCounts[normalized]++
		}
	}

	cleaned := make([]string, 0, len(normalizedLines))
	lastWasBlank := true
	for _, line := range normalizedLines {
		if line == "" {
			if !lastWasBlank {
				cleaned = append(cleaned, "")
				lastWasBlank = true
			}
			continue
		}

		if isNoiseLine(line, lineCounts[line]) {
			continue
		}

		cleaned = append(cleaned, line)
		lastWasBlank = false
	}

	return strings.TrimSpace(strings.Join(cleaned, "\n"))
}

func normalizeWhitespace(line string) string {
	return strings.Join(strings.Fields(line), " ")
}

func isNoiseLine(line string, count int) bool {
	if line == "" {
		return true
	}

	runeCount := utf8.RuneCountInString(line)
	if count >= 2 {
		if runeCount <= 3 {
			return true
		}
		if suspiciousASCIIBlockPattern.MatchString(line) && !looksLikeContactLine(line) {
			return true
		}
	}

	if runeCount >= 24 && suspiciousASCIIBlockPattern.MatchString(line) && !looksLikeContactLine(line) {
		return true
	}

	return false
}

func looksLikeContactLine(line string) bool {
	return strings.Contains(line, "@") || strings.Contains(line, "http://") || strings.Contains(line, "https://")
}

func textQualityScore(text string) float64 {
	if strings.TrimSpace(text) == "" {
		return 0
	}

	lines := strings.Split(text, "\n")
	nonEmptyLines := 0
	informativeLines := 0
	shortASCIIOnlyLines := 0
	totalRunes := 0
	contentRunes := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		nonEmptyLines++
		runeCount := utf8.RuneCountInString(trimmed)
		if runeCount >= 8 || containsCJK(trimmed) || strings.Contains(trimmed, "@") || strings.Contains(trimmed, "：") || strings.Contains(trimmed, ":") {
			informativeLines++
		}
		if runeCount <= 2 && isASCIIWord(trimmed) {
			shortASCIIOnlyLines++
		}

		for _, r := range trimmed {
			totalRunes++
			if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) || unicode.In(r, unicode.Han) || isCommonPunctuation(r) {
				contentRunes++
			}
		}
	}

	if nonEmptyLines == 0 || totalRunes == 0 {
		return 0
	}

	signalRatio := float64(informativeLines) / float64(nonEmptyLines)
	contentRatio := float64(contentRunes) / float64(totalRunes)
	lengthScore := math.Min(float64(utf8.RuneCountInString(text))/1200.0, 1.0)
	noisePenalty := math.Min(float64(shortASCIIOnlyLines)/float64(nonEmptyLines), 1.0)
	suspiciousPenalty := math.Min(float64(len(suspiciousASCIIBlockPattern.FindAllString(text, -1)))*0.08, 0.40)

	score := 0.38*signalRatio + 0.34*contentRatio + 0.20*lengthScore - 0.20*noisePenalty - suspiciousPenalty
	score += resumeSignalBonus(text)

	return clamp(score, 0, 1)
}

func containsCJK(text string) bool {
	for _, r := range text {
		if unicode.In(r, unicode.Han) {
			return true
		}
	}
	return false
}

func isASCIIWord(text string) bool {
	for _, r := range text {
		if r > unicode.MaxASCII || !(unicode.IsLetter(r) || unicode.IsDigit(r)) {
			return false
		}
	}
	return text != ""
}

func isCommonPunctuation(r rune) bool {
	switch r {
	case '，', '。', '；', '：', '（', '）', '、', '【', '】', '《', '》', ',', '.', ';', ':', '(', ')', '/', '-', '+', '&', '%', '_':
		return true
	default:
		return false
	}
}

func resumeSignalBonus(text string) float64 {
	signals := []string{
		"工作经历",
		"项目经历",
		"教育经历",
		"手机号",
		"毕业学校",
		"求职意愿",
		"技能",
		"GitHub",
		"E-Mail",
		"邮箱",
	}

	hits := 0
	for _, signal := range signals {
		if strings.Contains(text, signal) {
			hits++
		}
	}

	return math.Min(float64(hits)*0.03, 0.15)
}

func countPDFPages(filePath string) int {
	if isPDFInfoAvailable() {
		cmd := exec.Command("pdfinfo", filePath)
		output, err := cmd.Output()
		if err == nil {
			matches := pageCountPattern.FindStringSubmatch(string(output))
			if len(matches) == 2 {
				if pages, err := strconv.Atoi(matches[1]); err == nil {
					return pages
				}
			}
		}
	}

	f, r, err := pdf.Open(filePath)
	if err != nil {
		return 0
	}
	defer f.Close()

	return r.NumPage()
}

func preferredTesseractLanguage() string {
	tesseractLangOnce.Do(func() {
		langs := availableTesseractLanguages()

		switch {
		case langs["chi_sim"] && langs["eng"]:
			tesseractLang = "chi_sim+eng"
		case langs["chi_sim"]:
			tesseractLang = "chi_sim"
		case langs["chi_tra"] && langs["eng"]:
			tesseractLang = "chi_tra+eng"
		case langs["chi_tra"]:
			tesseractLang = "chi_tra"
		case langs["eng"]:
			tesseractLang = "eng"
		default:
			tesseractLang = ""
		}
	})

	return tesseractLang
}

func availableTesseractLanguages() map[string]bool {
	if !isTesseractAvailable() {
		return map[string]bool{}
	}

	cmd := exec.Command("tesseract", "--list-langs")
	output, err := cmd.Output()
	if err != nil {
		return map[string]bool{}
	}

	langs := make(map[string]bool)
	for _, line := range strings.Split(string(output), "\n") {
		lang := strings.TrimSpace(line)
		if lang == "" || strings.HasPrefix(lang, "List of available languages") {
			continue
		}
		langs[lang] = true
	}

	return langs
}

func clamp(value float64, min float64, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
