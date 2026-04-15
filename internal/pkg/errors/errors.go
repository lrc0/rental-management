package errors

import "errors"

// 标准错误
var (
	ErrBadRequest       = errors.New("请求参数错误")
	ErrUnauthorized     = errors.New("未授权访问")
	ErrForbidden        = errors.New("禁止访问")
	ErrNotFound         = errors.New("资源不存在")
	ErrInternalServer   = errors.New("服务器内部错误")
	ErrPhoneExists      = errors.New("手机号已注册")
	ErrPhoneNotFound    = errors.New("手机号未注册")
	ErrPasswordWrong    = errors.New("密码错误")
	ErrTokenInvalid     = errors.New("Token无效")
	ErrTokenExpired     = errors.New("Token已过期")
	ErrPermissionDenied = errors.New("无权限操作")
)

// 业务错误码
const (
	CodeSuccess         = 0
	CodeBadRequest      = 400
	CodeUnauthorized    = 401
	CodeForbidden       = 403
	CodeNotFound        = 404
	CodeInternalError   = 500
	CodePhoneExists     = 10001
	CodePhoneNotFound   = 10002
	CodePasswordWrong   = 10003
	CodeTokenInvalid    = 10004
	CodeTokenExpired    = 10005
	CodePermissionError = 10006
)

// 错误码消息映射
var codeMessages = map[int]string{
	CodeSuccess:         "成功",
	CodeBadRequest:      "请求参数错误",
	CodeUnauthorized:    "未授权访问",
	CodeForbidden:       "禁止访问",
	CodeNotFound:        "资源不存在",
	CodeInternalError:   "服务器内部错误",
	CodePhoneExists:     "手机号已注册",
	CodePhoneNotFound:   "手机号未注册",
	CodePasswordWrong:   "密码错误",
	CodeTokenInvalid:    "Token无效",
	CodeTokenExpired:    "Token已过期",
	CodePermissionError: "无权限操作",
}

func GetCodeMessage(code int) string {
	if msg, ok := codeMessages[code]; ok {
		return msg
	}
	return "未知错误"
}
