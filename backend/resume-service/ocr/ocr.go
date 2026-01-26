package ocr

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ledongthuc/pdf"
)

// OCRResult OCR识别结果
type OCRResult struct {
	Text       string  `json:"text"`
	Confidence float64 `json:"confidence"`
	Pages      int     `json:"pages"`
}

// ExtractTextFromPDF 从PDF文件提取文本
func ExtractTextFromPDF(filePath string) (*OCRResult, error) {
	f, r, err := pdf.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开PDF文件失败: %w", err)
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

	result := &OCRResult{
		Text:       buf.String(),
		Confidence: 0.95, // PDF直接提取文本置信度高
		Pages:      numPages,
	}

	// 如果提取的文本太少，可能是扫描版PDF，需要OCR
	if len(strings.TrimSpace(result.Text)) < 100 && numPages > 0 {
		// 尝试使用tesseract进行OCR
		ocrText, err := ocrPDFWithTesseract(filePath)
		if err == nil && len(ocrText) > len(result.Text) {
			result.Text = ocrText
			result.Confidence = 0.85 // OCR识别置信度较低
		}
	}

	return result, nil
}

// ExtractTextFromImage 从图片文件提取文本（使用tesseract）
func ExtractTextFromImage(imagePath string) (*OCRResult, error) {
	// 检查tesseract是否可用
	if !isTesseractAvailable() {
		return nil, fmt.Errorf("tesseract未安装，请先安装: brew install tesseract tesseract-lang")
	}

	// 调用tesseract进行OCR
	cmd := exec.Command("tesseract", imagePath, "stdout", "-l", "chi_sim+eng")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("OCR识别失败: %w", err)
	}

	return &OCRResult{
		Text:       string(output),
		Confidence: 0.85,
		Pages:      1,
	}, nil
}

// ocrPDFWithTesseract 使用tesseract对PDF进行OCR
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

	// 使用pdftoppm将PDF转换为图片（如果可用）
	if isPdftoppmAvailable() {
		imgPrefix := filepath.Join(tempDir, "page")
		cmd := exec.Command("pdftoppm", "-png", pdfPath, imgPrefix)
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("PDF转图片失败: %w", err)
		}

		// 对每个图片进行OCR
		var allText strings.Builder
		files, _ := filepath.Glob(imgPrefix + "*.png")
		for _, imgFile := range files {
			result, err := ExtractTextFromImage(imgFile)
			if err == nil {
				allText.WriteString(result.Text)
				allText.WriteString("\n")
			}
		}
		return allText.String(), nil
	}

	return "", fmt.Errorf("pdftoppm未安装，无法处理扫描版PDF")
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

// ExtractTextFromFile 根据文件类型自动选择提取方法
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

// extractTextFromWord 从Word文档提取文本
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
			Text:       string(output),
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
		Text:       text,
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
