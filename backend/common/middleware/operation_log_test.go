package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestOperationLogPreservesLargeJSONBody(t *testing.T) {
	t.Setenv("ES_URL", "http://127.0.0.1:1")
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(OperationLog(&OperationLogConfig{
		ServiceName:     "test-service",
		LogRequestBody:  true,
		LogResponseBody: false,
		MaxBodySize:     16,
	}))

	type updateProfileRequest struct {
		Avatar string `json:"avatar"`
	}

	router.PUT("/profile", func(c *gin.Context) {
		var req updateProfileRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"avatar_len": len(req.Avatar),
		})
	})

	body := map[string]string{
		"avatar": "data:image/png;base64," + strings.Repeat("a", 8192),
	}
	bodyBytes, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPut, "/profile", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	require.Equal(t, http.StatusOK, recorder.Code, recorder.Body.String())
	require.Contains(t, recorder.Body.String(), "\"avatar_len\":8214")
}
