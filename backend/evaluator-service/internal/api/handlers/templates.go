package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) GetJDTemplate(c *gin.Context) {
	p := filepath.Join(h.cfg.Storage.StaticDir, "templates", "jd.txt")
	b, _ := os.ReadFile(p)
	c.JSON(http.StatusOK, gin.H{"content": string(b)})
}
