package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"evaluator-service/internal/utils"
)

func (h *Handlers) ExportPDF(c *gin.Context) {
	markdown := c.PostForm("markdown")
	filename := c.DefaultPostForm("filename", "report")
	htmlBody, err := utils.MarkdownToHTML(markdown)
	if err != nil { fail(c, err); return }
	full := h.exprt.HTMLPageFromMarkdownBody(htmlBody)
	pdfBytes, err := h.exprt.PDFfromHTML(full)
	if err != nil { fail(c, err); return }
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename="+filename+".pdf")
	c.Header("Content-Length", "")
	c.Status(http.StatusOK)
	_, _ = c.Writer.Write(pdfBytes)
}

func (h *Handlers) ExportExcel(c *gin.Context) {
	data := c.PostForm("data")
	b, err := h.exprt.ExcelFromBatchJSON(data)
	if err != nil { fail(c, err); return }
	ts := time.Now().Format("20060102_150405")
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=resume_evaluation_"+ts+".xlsx")
	c.Status(http.StatusOK)
	_, _ = c.Writer.Write(b)
}

