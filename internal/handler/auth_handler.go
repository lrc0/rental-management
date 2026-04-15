package handler

import (
	"github.com/gin-gonic/gin"
	"rental-management/internal/middleware"
	"rental-management/internal/pkg/errors"
	"rental-management/internal/pkg/response"
	"rental-management/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register 用户注册
// @Summary 用户注册
// @Description 房东用户注册
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body service.RegisterRequest true "注册信息"
// @Success 200 {object} response.Response
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, "参数错误: "+err.Error())
		return
	}

	user, err := h.authService.Register(&req)
	if err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, err.Error())
		return
	}

	response.Success(c, user)
}

// Login 用户登录
// @Summary 用户登录
// @Description 房东用户登录
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body service.LoginRequest true "登录信息"
// @Success 200 {object} response.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, "参数错误: "+err.Error())
		return
	}

	resp, err := h.authService.Login(&req)
	if err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, err.Error())
		return
	}

	// 生成JWT token
	token, err := middleware.GenerateToken(resp.User.ID)
	if err != nil {
		response.Fail(c, errors.CodeInternalError)
		return
	}

	resp.Token = token
	response.Success(c, resp)
}

// GetProfile 获取用户信息
// @Summary 获取用户信息
// @Description 获取当前登录用户信息
// @Tags 认证
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Router /auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	user, err := h.authService.GetProfile(userID)
	if err != nil {
		response.Fail(c, errors.CodeNotFound)
		return
	}

	response.Success(c, user)
}

// UpdateProfile 更新用户信息
// @Summary 更新用户信息
// @Description 更新当前登录用户信息
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.UpdateProfileRequest true "用户信息"
// @Success 200 {object} response.Response
// @Router /auth/profile [put]
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req service.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, "参数错误: "+err.Error())
		return
	}

	user, err := h.authService.UpdateProfile(userID, &req)
	if err != nil {
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}

	response.Success(c, user)
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 修改当前登录用户密码
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.ChangePasswordRequest true "密码信息"
// @Success 200 {object} response.Response
// @Router /auth/password [put]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req service.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := h.authService.ChangePassword(userID, &req); err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, err.Error())
		return
	}

	response.Success(c, nil)
}
