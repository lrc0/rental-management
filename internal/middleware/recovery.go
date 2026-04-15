package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"rental-management/internal/pkg/errors"
	"rental-management/internal/pkg/response"
)

// Recovery 异常恢复中间件
func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("stack", string(debug.Stack())),
				)

				c.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{
					Code:    errors.CodeInternalError,
					Message: "服务器内部错误",
				})
			}
		}()
		c.Next()
	}
}
