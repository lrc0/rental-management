package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"rental-management/internal/middleware"
	"rental-management/internal/pkg/errors"
	"rental-management/internal/pkg/response"
	"rental-management/internal/service"
)

type BillHandler struct {
	billService *service.BillService
}

func NewBillHandler(billService *service.BillService) *BillHandler {
	return &BillHandler{billService: billService}
}

// CreateMeterReading 创建抄表记录
// @Summary 创建抄表记录
// @Description 录入水电气表读数
// @Tags 抄表管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.CreateMeterReadingRequest true "抄表信息"
// @Success 200 {object} response.Response
// @Router /meter-readings [post]
func (h *BillHandler) CreateMeterReading(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req service.CreateMeterReadingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, "参数错误: "+err.Error())
		return
	}

	reading, err := h.billService.CreateMeterReading(userID, &req)
	if err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, err.Error())
		return
	}

	response.Success(c, reading)
}

// ListMeterReadings 获取抄表记录列表
// @Summary 获取抄表记录列表
// @Description 获取抄表记录列表
// @Tags 抄表管理
// @Produce json
// @Security BearerAuth
// @Param room_id query int false "房间ID"
// @Param start_date query string false "开始日期"
// @Param end_date query string false "结束日期"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} response.Response
// @Router /meter-readings [get]
func (h *BillHandler) ListMeterReadings(c *gin.Context) {
	userID := middleware.GetUserID(c)
	roomID, _ := strconv.ParseUint(c.Query("room_id"), 10, 64)

	var startDate, endDate *time.Time
	if s := c.Query("start_date"); s != "" {
		t, err := time.Parse("2006-01-02", s)
		if err == nil {
			startDate = &t
		}
	}
	if s := c.Query("end_date"); s != "" {
		t, err := time.Parse("2006-01-02", s)
		if err == nil {
			endDate = &t
		}
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	readings, total, err := h.billService.ListMeterReadings(userID, uint(roomID), startDate, endDate, page, pageSize)
	if err != nil {
		response.Fail(c, errors.CodeInternalError)
		return
	}

	response.Page(c, readings, total, page, pageSize)
}

// CreateBill 生成账单
// @Summary 生成账单
// @Description 生成账单
// @Tags 账单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.CreateBillRequest true "账单信息"
// @Success 200 {object} response.Response
// @Router /bills [post]
func (h *BillHandler) CreateBill(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req service.CreateBillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, "参数错误: "+err.Error())
		return
	}

	bill, err := h.billService.CreateBill(userID, &req)
	if err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, err.Error())
		return
	}

	response.Success(c, bill)
}

// GetBill 获取账单详情
// @Summary 获取账单详情
// @Description 根据ID获取账单详情
// @Tags 账单管理
// @Produce json
// @Security BearerAuth
// @Param id path int true "账单ID"
// @Success 200 {object} response.Response
// @Router /bills/{id} [get]
func (h *BillHandler) GetBill(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, "无效的ID")
		return
	}

	bill, err := h.billService.GetBill(uint(id), userID)
	if err != nil {
		response.Fail(c, errors.CodeNotFound)
		return
	}

	response.Success(c, bill)
}

// ListBills 获取账单列表
// @Summary 获取账单列表
// @Description 获取账单列表
// @Tags 账单管理
// @Produce json
// @Security BearerAuth
// @Param room_id query int false "房间ID"
// @Param status query int false "状态"
// @Param bill_month query string false "账单月份"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} response.Response
// @Router /bills [get]
func (h *BillHandler) ListBills(c *gin.Context) {
	userID := middleware.GetUserID(c)
	roomID, _ := strconv.ParseUint(c.Query("room_id"), 10, 64)
	billMonth := c.Query("bill_month")

	var status *int8
	if s := c.Query("status"); s != "" {
		st, _ := strconv.Atoi(s)
		st8 := int8(st)
		status = &st8
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	bills, total, err := h.billService.ListBills(userID, uint(roomID), status, billMonth, page, pageSize)
	if err != nil {
		response.Fail(c, errors.CodeInternalError)
		return
	}

	response.Page(c, bills, total, page, pageSize)
}

// PayBill 支付账单
// @Summary 支付账单
// @Description 标记账单已支付
// @Tags 账单管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "账单ID"
// @Param request body service.PayBillRequest true "支付信息"
// @Success 200 {object} response.Response
// @Router /bills/{id}/pay [put]
func (h *BillHandler) PayBill(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, "无效的ID")
		return
	}

	var req service.PayBillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := h.billService.PayBill(uint(id), userID, &req); err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetBillStatistics 获取账单统计
// @Summary 获取账单统计
// @Description 获取账单统计
// @Tags 账单管理
// @Produce json
// @Security BearerAuth
// @Param start_date query string false "开始日期"
// @Param end_date query string false "结束日期"
// @Success 200 {object} response.Response
// @Router /bills/statistics [get]
func (h *BillHandler) GetBillStatistics(c *gin.Context) {
	userID := middleware.GetUserID(c)

	// 默认统计本月
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 1, -1)

	if s := c.Query("start_date"); s != "" {
		t, err := time.Parse("2006-01-02", s)
		if err == nil {
			startDate = t
		}
	}
	if s := c.Query("end_date"); s != "" {
		t, err := time.Parse("2006-01-02", s)
		if err == nil {
			endDate = t
		}
	}

	stats, err := h.billService.GetBillStatistics(userID, startDate, endDate)
	if err != nil {
		response.Fail(c, errors.CodeInternalError)
		return
	}

	response.Success(c, stats)
}

// GetMonthlyStatistics 获取月度统计
// @Summary 获取月度统计
// @Description 获取年度月度统计
// @Tags 账单管理
// @Produce json
// @Security BearerAuth
// @Param year query int false "年份" default(2024)
// @Success 200 {object} response.Response
// @Router /bills/monthly-statistics [get]
func (h *BillHandler) GetMonthlyStatistics(c *gin.Context) {
	userID := middleware.GetUserID(c)
	year, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(time.Now().Year())))

	stats, err := h.billService.GetMonthlyStatistics(userID, year)
	if err != nil {
		response.Fail(c, errors.CodeInternalError)
		return
	}

	response.Success(c, stats)
}

// GetFeeRate 获取费率配置
// @Summary 获取费率配置
// @Description 获取水电气费率配置
// @Tags 系统配置
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Router /fee-rates [get]
func (h *BillHandler) GetFeeRate(c *gin.Context) {
	userID := middleware.GetUserID(c)

	feeRate, err := h.billService.GetFeeRate(userID)
	if err != nil {
		response.Fail(c, errors.CodeNotFound)
		return
	}

	response.Success(c, feeRate)
}

// UpdateFeeRate 更新费率配置
// @Summary 更新费率配置
// @Description 更新水电气费率配置
// @Tags 系统配置
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.UpdateFeeRateRequest true "费率信息"
// @Success 200 {object} response.Response
// @Router /fee-rates [put]
func (h *BillHandler) UpdateFeeRate(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req service.UpdateFeeRateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, "参数错误: "+err.Error())
		return
	}

	feeRate, err := h.billService.UpdateFeeRate(userID, &req)
	if err != nil {
		response.FailWithMsg(c, errors.CodeInternalError, err.Error())
		return
	}

	response.Success(c, feeRate)
}

// DeleteBill 删除账单
// @Summary 删除账单
// @Description 删除账单
// @Tags 账单管理
// @Produce json
// @Security BearerAuth
// @Param id path int true "账单ID"
// @Success 200 {object} response.Response
// @Router /bills/{id} [delete]
func (h *BillHandler) DeleteBill(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, "无效的ID")
		return
	}

	if err := h.billService.DeleteBill(uint(id), userID); err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, err.Error())
		return
	}

	response.Success(c, nil)
}

// DeleteMeterReading 删除抄表记录
// @Summary 删除抄表记录
// @Description 删除抄表记录
// @Tags 抄表管理
// @Produce json
// @Security BearerAuth
// @Param id path int true "抄表记录ID"
// @Success 200 {object} response.Response
// @Router /meter-readings/{id} [delete]
func (h *BillHandler) DeleteMeterReading(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, "无效的ID")
		return
	}

	if err := h.billService.DeleteMeterReading(uint(id), userID); err != nil {
		response.FailWithMsg(c, errors.CodeBadRequest, err.Error())
		return
	}

	response.Success(c, nil)
}
