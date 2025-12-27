package handlers

import (
	"net/http"
	"strings"

	"evaluator-service/internal/api/middleware"
	"evaluator-service/internal/logging"
	"evaluator-service/internal/utils"

	"github.com/gin-gonic/gin"
)

type credentialUpsertReq struct {
	Org      string `json:"org" binding:"required"`
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type credentialStatusResp struct {
	Exists  bool   `json:"exists"`
	Account string `json:"account,omitempty"`
}

func maskAccount(acct string) string {
	acct = strings.TrimSpace(acct)
	if acct == "" {
		return ""
	}
	if i := strings.Index(acct, "@"); i > 1 {
		prefix := acct[:i]
		if len(prefix) <= 2 {
			return "**" + acct[i:]
		}
		return prefix[:2] + strings.Repeat("*", len(prefix)-2) + acct[i:]
	}
	if len(acct) <= 2 {
		return "**"
	}
	return acct[:2] + strings.Repeat("*", len(acct)-2)
}

// CredentialStatus GET /api/credentials/status?org=motern
func (h *Handlers) CredentialStatus(c *gin.Context) {
	userID := middleware.GetUserID(c)
	org := c.Query("org")
	if strings.TrimSpace(org) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "org is required"})
		return
	}
	cred, err := h.credRepo.GetByOrgAndUser(org, userID)
	if err != nil {
		fail(c, err)
		return
	}
	if cred == nil {
		ok(c, credentialStatusResp{Exists: false})
		return
	}
	ok(c, credentialStatusResp{Exists: true, Account: maskAccount(cred.Account)})
}

// CredentialUpsert POST /api/credentials
func (h *Handlers) CredentialUpsert(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req credentialUpsertReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	encKey := strings.TrimSpace(h.cfg.Credentials.EncKey)
	passwordCipher := req.Password
	if encKey != "" {
		if enc, err := utils.EncryptAESGCM([]byte(encKey), req.Password); err != nil {
			h.log.Error("encrypt credentials failed", logging.Err(err))
		} else {
			passwordCipher = enc
		}
	} else {
		h.log.Warn("credentials.enc_key is empty; storing password in plaintext (NOT RECOMMENDED)")
	}
	if _, err := h.credRepo.UpsertByUser(req.Org, req.Account, passwordCipher, userID); err != nil {
		fail(c, err)
		return
	}
	ok(c, gin.H{"success": true})
}
