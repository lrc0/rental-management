package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"rental-management/internal/middleware"
	"rental-management/internal/pkg/errors"
	"rental-management/internal/pkg/response"
	"rental-management/internal/service"
)

type TenantHandler struct {
	tenantService *service.TenantService
}

func NewTenantHandler(tenantService *service.TenantService) *TenantHandler {
	return &TenantHandler{tenantService: tenantService}
}

// CreateTenant 创建租客
// @Summary 创建租客
// @Description 创建新租客
// @Tags 租客管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.CreateTenantRequest true "租客信息"
// @Success 200 {object} response.Response
// @Router /tenants [post]
func (h *TenantHandler) CreateTenant(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req service.CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, "参数错误: "+err.Error())
		return
	}

	tenant, err := h.tenantService.CreateTenant(userID, &req)
	if err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, err.Error())
		return
	}

	response.Success(c, tenant)
}

// GetTenant 获取租客详情
// @Summary 获取租客详情
// @Description 根据ID获取租客详情
// @Tags 租客管理
// @Produce json
// @Security BearerAuth
// @Param id path int true "租客ID"
// @Success 200 {object} response.Response
// @Router /tenants/{id} [get]
func (h *TenantHandler) GetTenant(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, "无效的ID")
		return
	}

	tenant, err := h.tenantService.GetTenant(uint(id), userID)
	if err != nil {
		response.Fail(c, errors.CodeNotFound)
		return
	}

	response.Success(c, tenant)
}

// ListTenants 获取租客列表
// @Summary 获取租客列表
// @Description 获取房东的租客列表
// @Tags 租客管理
// @Produce json
// @Security BearerAuth
// @Param status query int false "状态"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} response.Response
// @Router /tenants [get]
func (h *TenantHandler) ListTenants(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var status *int8
	if s := c.Query("status"); s != "" {
		st, _ := strconv.Atoi(s)
		st8 := int8(st)
		status = &st8
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	tenants, total, err := h.tenantService.ListTenants(userID, status, page, pageSize)
	if err != nil {
		response.Fail(c, errors.CodeInternalError)
		return
	}

	response.Page(c, tenants, total, page, pageSize)
}

// UpdateTenant 更新租客
// @Summary 更新租客
// @Description 更新租客信息
// @Tags 租客管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "租客ID"
// @Param request body service.UpdateTenantRequest true "租客信息"
// @Success 200 {object} response.Response
// @Router /tenants/{id} [put]
func (h *TenantHandler) UpdateTenant(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, "无效的ID")
		return
	}

	var req service.UpdateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, "参数错误: "+err.Error())
		return
	}

	tenant, err := h.tenantService.UpdateTenant(uint(id), userID, &req)
	if err != nil {
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}

	response.Success(c, tenant)
}

// DeleteTenant 删除租客
// @Summary 删除租客
// @Description 删除租客
// @Tags 租客管理
// @Produce json
// @Security BearerAuth
// @Param id path int true "租客ID"
// @Success 200 {object} response.Response
// @Router /tenants/{id} [delete]
func (h *TenantHandler) DeleteTenant(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, "无效的ID")
		return
	}

	if err := h.tenantService.DeleteTenant(uint(id), userID); err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, err.Error())
		return
	}

	response.Success(c, nil)
}

// CreateContract 签订合同
// @Summary 签订合同
// @Description 为租客签订租房合同
// @Tags 合同管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.CreateContractRequest true "合同信息"
// @Success 200 {object} response.Response
// @Router /contracts [post]
func (h *TenantHandler) CreateContract(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req service.CreateContractRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, "参数错误: "+err.Error())
		return
	}

	contract, err := h.tenantService.CreateContract(userID, &req)
	if err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, err.Error())
		return
	}

	response.Success(c, contract)
}

// GetContract 获取合同详情
// @Summary 获取合同详情
// @Description 根据ID获取合同详情
// @Tags 合同管理
// @Produce json
// @Security BearerAuth
// @Param id path int true "合同ID"
// @Success 200 {object} response.Response
// @Router /contracts/{id} [get]
func (h *TenantHandler) GetContract(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, "无效的ID")
		return
	}

	contract, err := h.tenantService.GetContract(uint(id), userID)
	if err != nil {
		response.Fail(c, errors.CodeNotFound)
		return
	}

	response.Success(c, contract)
}

// ListContracts 获取合同列表
// @Summary 获取合同列表
// @Description 获取合同列表
// @Tags 合同管理
// @Produce json
// @Security BearerAuth
// @Param status query int false "状态"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} response.Response
// @Router /contracts [get]
func (h *TenantHandler) ListContracts(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var status *int8
	if s := c.Query("status"); s != "" {
		st, _ := strconv.Atoi(s)
		st8 := int8(st)
		status = &st8
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	contracts, total, err := h.tenantService.ListContracts(userID, status, page, pageSize)
	if err != nil {
		response.Fail(c, errors.CodeInternalError)
		return
	}

	response.Page(c, contracts, total, page, pageSize)
}

// TerminateContract 解约
// @Summary 解约
// @Description 终止合同
// @Tags 合同管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "合同ID"
// @Param request body service.TerminateContractRequest true "解约原因"
// @Success 200 {object} response.Response
// @Router /contracts/{id}/terminate [put]
func (h *TenantHandler) TerminateContract(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, "无效的ID")
		return
	}

	var req service.TerminateContractRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := h.tenantService.TerminateContract(uint(id), userID, &req); err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, err.Error())
		return
	}

	response.Success(c, nil)
}

// UpdateContract 更新合同
// @Summary 更新合同
// @Description 更新合同信息
// @Tags 合同管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "合同ID"
// @Param request body service.UpdateContractRequest true "合同信息"
// @Success 200 {object} response.Response
// @Router /contracts/{id} [put]
func (h *TenantHandler) UpdateContract(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, "无效的ID")
		return
	}

	var req service.UpdateContractRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, "参数错误: "+err.Error())
		return
	}

	contract, err := h.tenantService.UpdateContract(uint(id), userID, &req)
	if err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, err.Error())
		return
	}

	response.Success(c, contract)
}

// DeleteContract 删除合同
// @Summary 删除合同
// @Description 删除合同
// @Tags 合同管理
// @Produce json
// @Security BearerAuth
// @Param id path int true "合同ID"
// @Success 200 {object} response.Response
// @Router /contracts/{id} [delete]
func (h *TenantHandler) DeleteContract(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, "无效的ID")
		return
	}

	if err := h.tenantService.DeleteContract(uint(id), userID); err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, err.Error())
		return
	}

	response.Success(c, nil)
}
