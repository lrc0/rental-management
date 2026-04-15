package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"rental-management/internal/pkg/errors"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PageData 分页数据结构
type PageData struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    errors.CodeSuccess,
		Message: errors.GetCodeMessage(errors.CodeSuccess),
		Data:    data,
	})
}

// SuccessWithMsg 成功响应（自定义消息）
func SuccessWithMsg(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    errors.CodeSuccess,
		Message: msg,
		Data:    data,
	})
}

// Fail 失败响应
func Fail(c *gin.Context, code int) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: errors.GetCodeMessage(code),
	})
}

// FailWithMsg 失败响应（自定义消息）
func FailWithMsg(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: msg,
	})
}

// FailWithData 失败响应（带数据）
func FailWithData(c *gin.Context, code int, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: errors.GetCodeMessage(code),
		Data:    data,
	})
}

// Page 分页响应
func Page(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	Success(c, PageData{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}
