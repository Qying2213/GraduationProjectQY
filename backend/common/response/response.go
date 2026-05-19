// Package response 定义多个服务可复用的统一 JSON 响应结构。
// 统一的 `{code, message, data}` 格式可以降低前端错误处理和 Swagger 示例说明的复杂度。
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 返回标准业务成功响应，约定 code=0 表示成功。
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 用于需要自定义成功文案但仍保持统一响应结构的场景。
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: message,
		Data:    data,
	})
}

// Error 保持 HTTP 200，但用业务 code=1 表示失败。
// 部分已有页面依赖这种约定；新接口也可以按场景使用真实 HTTP 状态码表达错误。
func Error(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    1,
		Message: message,
	})
}

// ErrorWithCode 在统一响应结构中返回自定义业务错误码。
func ErrorWithCode(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
	})
}

// Fail 用真实 HTTP 状态码表达认证、权限、服务端异常等传输层/系统层失败。
func Fail(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, Response{
		Code:    statusCode,
		Message: message,
	})
}
