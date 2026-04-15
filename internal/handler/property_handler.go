package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"rental-management/internal/middleware"
	"rental-management/internal/pkg/errors"
	"rental-management/internal/pkg/response"
	"rental-management/internal/service"
)

type PropertyHandler struct {
	propertyService *service.PropertyService
}

func NewPropertyHandler(propertyService *service.PropertyService) *PropertyHandler {
	return &PropertyHandler{propertyService: propertyService}
}

// CreateProperty 创建房源
// @Summary 创建房源
// @Description 创建新房源
// @Tags 房源管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.CreatePropertyRequest true "房源信息"
// @Success 200 {object} response.Response
// @Router /properties [post]
func (h *PropertyHandler) CreateProperty(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req service.CreatePropertyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, "参数错误: "+err.Error())
		return
	}

	property, err := h.propertyService.CreateProperty(userID, &req)
	if err != nil {
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}

	response.Success(c, property)
}

// GetProperty 获取房源详情
// @Summary 获取房源详情
// @Description 根据ID获取房源详情
// @Tags 房源管理
// @Produce json
// @Security BearerAuth
// @Param id path int true "房源ID"
// @Success 200 {object} response.Response
// @Router /properties/{id} [get]
func (h *PropertyHandler) GetProperty(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, "无效的ID")
		return
	}

	property, err := h.propertyService.GetProperty(uint(id), userID)
	if err != nil {
		response.Fail(c, errors.CodeNotFound)
		return
	}

	response.Success(c, property)
}

// ListProperties 获取房源列表
// @Summary 获取房源列表
// @Description 获取房东的房源列表
// @Tags 房源管理
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} response.Response
// @Router /properties [get]
func (h *PropertyHandler) ListProperties(c *gin.Context) {
	userID := middleware.GetUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	properties, total, err := h.propertyService.ListProperties(userID, page, pageSize)
	if err != nil {
		response.Fail(c, errors.CodeInternalError)
		return
	}

	response.Page(c, properties, total, page, pageSize)
}

// UpdateProperty 更新房源
// @Summary 更新房源
// @Description 更新房源信息
// @Tags 房源管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "房源ID"
// @Param request body service.UpdatePropertyRequest true "房源信息"
// @Success 200 {object} response.Response
// @Router /properties/{id} [put]
func (h *PropertyHandler) UpdateProperty(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, "无效的ID")
		return
	}

	var req service.UpdatePropertyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, "参数错误: "+err.Error())
		return
	}

	property, err := h.propertyService.UpdateProperty(uint(id), userID, &req)
	if err != nil {
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}

	response.Success(c, property)
}

// DeleteProperty 删除房源
// @Summary 删除房源
// @Description 删除房源
// @Tags 房源管理
// @Produce json
// @Security BearerAuth
// @Param id path int true "房源ID"
// @Success 200 {object} response.Response
// @Router /properties/{id} [delete]
func (h *PropertyHandler) DeleteProperty(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, "无效的ID")
		return
	}

	if err := h.propertyService.DeleteProperty(uint(id), userID); err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, err.Error())
		return
	}

	response.Success(c, nil)
}

// CreateRoom 创建房间
// @Summary 创建房间
// @Description 在房源下创建房间
// @Tags 房间管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.CreateRoomRequest true "房间信息"
// @Success 200 {object} response.Response
// @Router /rooms [post]
func (h *PropertyHandler) CreateRoom(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req service.CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, "参数错误: "+err.Error())
		return
	}

	room, err := h.propertyService.CreateRoom(userID, &req)
	if err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, err.Error())
		return
	}

	response.Success(c, room)
}

// GetRoom 获取房间详情
// @Summary 获取房间详情
// @Description 根据ID获取房间详情
// @Tags 房间管理
// @Produce json
// @Security BearerAuth
// @Param id path int true "房间ID"
// @Success 200 {object} response.Response
// @Router /rooms/{id} [get]
func (h *PropertyHandler) GetRoom(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, "无效的ID")
		return
	}

	room, err := h.propertyService.GetRoom(uint(id), userID)
	if err != nil {
		response.Fail(c, errors.CodeNotFound)
		return
	}

	response.Success(c, room)
}

// ListRooms 获取房间列表
// @Summary 获取房间列表
// @Description 获取房间列表
// @Tags 房间管理
// @Produce json
// @Security BearerAuth
// @Param property_id query int false "房源ID"
// @Param status query int false "状态"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} response.Response
// @Router /rooms [get]
func (h *PropertyHandler) ListRooms(c *gin.Context) {
	userID := middleware.GetUserID(c)
	propertyID, _ := strconv.ParseUint(c.Query("property_id"), 10, 64)

	var status *int8
	if s := c.Query("status"); s != "" {
		st, _ := strconv.Atoi(s)
		st8 := int8(st)
		status = &st8
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	rooms, total, err := h.propertyService.ListRooms(userID, uint(propertyID), status, page, pageSize)
	if err != nil {
		response.Fail(c, errors.CodeInternalError)
		return
	}

	response.Page(c, rooms, total, page, pageSize)
}

// UpdateRoom 更新房间
// @Summary 更新房间
// @Description 更新房间信息
// @Tags 房间管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "房间ID"
// @Param request body service.UpdateRoomRequest true "房间信息"
// @Success 200 {object} response.Response
// @Router /rooms/{id} [put]
func (h *PropertyHandler) UpdateRoom(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, "无效的ID")
		return
	}

	var req service.UpdateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, "参数错误: "+err.Error())
		return
	}

	room, err := h.propertyService.UpdateRoom(uint(id), userID, &req)
	if err != nil {
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}

	response.Success(c, room)
}

// DeleteRoom 删除房间
// @Summary 删除房间
// @Description 删除房间
// @Tags 房间管理
// @Produce json
// @Security BearerAuth
// @Param id path int true "房间ID"
// @Success 200 {object} response.Response
// @Router /rooms/{id} [delete]
func (h *PropertyHandler) DeleteRoom(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, "无效的ID")
		return
	}

	if err := h.propertyService.DeleteRoom(uint(id), userID); err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, err.Error())
		return
	}

	response.Success(c, nil)
}

// UpdateRoomStatus 更新房间状态
// @Summary 更新房间状态
// @Description 更新房间状态
// @Tags 房间管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "房间ID"
// @Param request body map[string]int8 true "状态"
// @Success 200 {object} response.Response
// @Router /rooms/{id}/status [put]
func (h *PropertyHandler) UpdateRoomStatus(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, "无效的ID")
		return
	}

	var req struct {
		Status int8 `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := h.propertyService.UpdateRoomStatus(uint(id), userID, req.Status); err != nil {
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}

	response.Success(c, nil)
}
