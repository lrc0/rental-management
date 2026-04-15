package service

import (
	"errors"
	"time"

	"rental-management/internal/model"
	"rental-management/internal/repository"
)

type BillService struct {
	billRepo     *repository.BillRepository
	roomRepo     *repository.RoomRepository
	userRepo     *repository.UserRepository
	contractRepo *repository.ContractRepository
}

func NewBillService(
	billRepo *repository.BillRepository,
	roomRepo *repository.RoomRepository,
	userRepo *repository.UserRepository,
	contractRepo *repository.ContractRepository,
) *BillService {
	return &BillService{
		billRepo:     billRepo,
		roomRepo:     roomRepo,
		userRepo:     userRepo,
		contractRepo: contractRepo,
	}
}

// CreateMeterReadingRequest 抄表请求
type CreateMeterReadingRequest struct {
	RoomID             uint    `json:"room_id" binding:"required"`
	ReadingDate        string  `json:"reading_date" binding:"required"` // YYYY-MM-DD
	WaterReading       float64 `json:"water_reading"`
	ElectricityReading float64 `json:"electricity_reading"`
	GasReading         float64 `json:"gas_reading"`
	Remark             string  `json:"remark"`
}

// MeterReadingResponse 抄表响应
type MeterReadingResponse struct {
	*model.MeterReading
	WaterUsage       float64 `json:"water_usage"`
	ElectricityUsage float64 `json:"electricity_usage"`
	GasUsage         float64 `json:"gas_usage"`
}

// CreateMeterReading 创建抄表记录
func (s *BillService) CreateMeterReading(userID uint, req *CreateMeterReadingRequest) (*MeterReadingResponse, error) {
	// 验证房间归属
	room, err := s.roomRepo.FindByIDAndUserID(req.RoomID, userID)
	if err != nil {
		return nil, errors.New("房间不存在或无权限")
	}

	// 解析日期
	readingDate, err := time.Parse("2006-01-02", req.ReadingDate)
	if err != nil {
		return nil, errors.New("日期格式错误")
	}

	// 获取上次抄表记录
	var waterUsage, electricityUsage, gasUsage float64
	lastReading, err := s.billRepo.GetLatestMeterReading(room.ID, readingDate)
	if err == nil && lastReading != nil {
		waterUsage = req.WaterReading - lastReading.WaterReading
		electricityUsage = req.ElectricityReading - lastReading.ElectricityReading
		gasUsage = req.GasReading - lastReading.GasReading

		// 防止负数（可能是换表）
		if waterUsage < 0 {
			waterUsage = 0
		}
		if electricityUsage < 0 {
			electricityUsage = 0
		}
		if gasUsage < 0 {
			gasUsage = 0
		}
	}

	reading := &model.MeterReading{
		RoomID:             req.RoomID,
		ReadingDate:        readingDate,
		WaterReading:       req.WaterReading,
		ElectricityReading: req.ElectricityReading,
		GasReading:         req.GasReading,
		WaterUsage:         waterUsage,
		ElectricityUsage:   electricityUsage,
		GasUsage:           gasUsage,
		Remark:             req.Remark,
	}

	if err := s.billRepo.CreateMeterReading(reading); err != nil {
		return nil, err
	}

	return &MeterReadingResponse{
		MeterReading:     reading,
		WaterUsage:       waterUsage,
		ElectricityUsage: electricityUsage,
		GasUsage:         gasUsage,
	}, nil
}

// ListMeterReadings 获取抄表记录列表
func (s *BillService) ListMeterReadings(userID uint, roomID uint, startDate, endDate *time.Time, page, pageSize int) ([]model.MeterReading, int64, error) {
	return s.billRepo.ListMeterReadings(userID, roomID, startDate, endDate, page, pageSize)
}

// CreateBillRequest 生成账单请求
type CreateBillRequest struct {
	RoomID      uint   `json:"room_id" binding:"required"`
	TenantID    uint   `json:"tenant_id"`
	BillMonth   string `json:"bill_month" binding:"required"` // YYYY-MM
	DueDate     string `json:"due_date"`
	RentFee     float64 `json:"rent_fee"`
	WaterFee    float64 `json:"water_fee"`
	ElectricityFee float64 `json:"electricity_fee"`
	GasFee      float64 `json:"gas_fee"`
}

// CreateBill 生成账单
func (s *BillService) CreateBill(userID uint, req *CreateBillRequest) (*model.Bill, error) {
	// 验证房间归属
	room, err := s.roomRepo.FindByIDAndUserID(req.RoomID, userID)
	if err != nil {
		return nil, errors.New("房间不存在或无权限")
	}

	// 检查账单是否已存在
	exists, err := s.billRepo.CheckBillExists(room.ID, req.BillMonth)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("该月份账单已存在")
	}

	// 计算总金额
	amount := req.RentFee + req.WaterFee + req.ElectricityFee + req.GasFee

	bill := &model.Bill{
		UserID:         userID,
		RoomID:         req.RoomID,
		TenantID:       req.TenantID,
		BillType:       model.BillTypeCombined,
		BillMonth:      req.BillMonth,
		Amount:         amount,
		WaterFee:       req.WaterFee,
		ElectricityFee: req.ElectricityFee,
		GasFee:         req.GasFee,
		RentFee:        req.RentFee,
		Status:         model.BillStatusPending,
	}

	// 设置应付日期
	if req.DueDate != "" {
		dueDate, err := time.Parse("2006-01-02", req.DueDate)
		if err == nil {
			bill.DueDate = &dueDate
		}
	}

	if err := s.billRepo.CreateBill(bill); err != nil {
		return nil, err
	}

	return bill, nil
}

// GenerateBillFromReadings 根据抄表记录自动生成账单
func (s *BillService) GenerateBillFromReadings(userID, roomID uint, billMonth string) (*model.Bill, error) {
	// 验证房间归属
	room, err := s.roomRepo.FindByIDAndUserID(roomID, userID)
	if err != nil {
		return nil, errors.New("房间不存在或无权限")
	}

	// 检查账单是否已存在
	exists, err := s.billRepo.CheckBillExists(room.ID, billMonth)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("该月份账单已存在")
	}

	// 获取费率配置
	_, err = s.userRepo.GetFeeRateByUserID(userID)
	if err != nil {
		// 使用默认费率
	}

	// 获取当前有效合同
	var tenantID uint
	var monthlyRent float64
	contract, err := s.contractRepo.FindActiveByRoomID(room.ID)
	if err == nil && contract != nil {
		tenantID = contract.TenantID
		monthlyRent = contract.MonthlyRent
	}

	// 解析月份，获取该月的抄表记录
	monthStart, _ := time.Parse("2006-01", billMonth)
	_ = monthStart // 暂时未使用，后续可根据抄表记录计算

	// 这里简化处理，实际应该根据抄表记录计算用量
	// 实际项目中需要更复杂的逻辑

	bill := &model.Bill{
		UserID:     userID,
		RoomID:     room.ID,
		TenantID:   tenantID,
		BillType:   model.BillTypeCombined,
		BillMonth:  billMonth,
		RentFee:    monthlyRent,
		WaterFee:   0,
		ElectricityFee: 0,
		GasFee:    0,
		Amount:    monthlyRent,
		Status:    model.BillStatusPending,
	}

	if err := s.billRepo.CreateBill(bill); err != nil {
		return nil, err
	}

	return bill, nil
}

// GetBill 获取账单详情
func (s *BillService) GetBill(id, userID uint) (*model.Bill, error) {
	return s.billRepo.FindBillByIDAndUserID(id, userID)
}

// ListBills 获取账单列表
func (s *BillService) ListBills(userID uint, roomID uint, status *int8, billMonth string, page, pageSize int) ([]model.Bill, int64, error) {
	return s.billRepo.ListBills(userID, roomID, status, billMonth, page, pageSize)
}

// PayBillRequest 支付账单请求
type PayBillRequest struct {
	Amount        float64 `json:"amount" binding:"required"`
	PaymentMethod int8    `json:"payment_method" binding:"required"`
	Note          string  `json:"note"`
}

// PayBill 支付账单
func (s *BillService) PayBill(id, userID uint, req *PayBillRequest) error {
	bill, err := s.billRepo.FindBillByIDAndUserID(id, userID)
	if err != nil {
		return errors.New("账单不存在")
	}

	if bill.Status == model.BillStatusPaid {
		return errors.New("账单已支付")
	}

	// 创建收款记录
	payment := &model.Payment{
		BillID:        bill.ID,
		Amount:        req.Amount,
		PaymentMethod: req.PaymentMethod,
		PaidAt:        time.Now(),
		Note:          req.Note,
	}

	if err := s.billRepo.CreatePayment(payment); err != nil {
		return err
	}

	// 更新账单状态
	now := time.Now()
	bill.Status = model.BillStatusPaid
	bill.PaidAt = &now

	return s.billRepo.UpdateBill(bill)
}

// GetBillStatistics 获取账单统计
func (s *BillService) GetBillStatistics(userID uint, startDate, endDate time.Time) (map[string]interface{}, error) {
	return s.billRepo.GetBillStatistics(userID, startDate, endDate)
}

// GetMonthlyStatistics 获取月度统计
func (s *BillService) GetMonthlyStatistics(userID uint, year int) ([]map[string]interface{}, error) {
	return s.billRepo.GetMonthlyStatistics(userID, year)
}

// UpdateFeeRateRequest 更新费率请求
type UpdateFeeRateRequest struct {
	WaterRate       float64 `json:"water_rate"`
	ElectricityRate float64 `json:"electricity_rate"`
	GasRate         float64 `json:"gas_rate"`
}

// GetFeeRate 获取费率配置
func (s *BillService) GetFeeRate(userID uint) (*model.FeeRate, error) {
	return s.userRepo.GetFeeRateByUserID(userID)
}

// UpdateFeeRate 更新费率配置
func (s *BillService) UpdateFeeRate(userID uint, req *UpdateFeeRateRequest) (*model.FeeRate, error) {
	feeRate, err := s.userRepo.GetFeeRateByUserID(userID)
	if err != nil {
		// 如果不存在，创建新的
		feeRate = &model.FeeRate{
			UserID: userID,
		}
	}

	feeRate.WaterRate = req.WaterRate
	feeRate.ElectricityRate = req.ElectricityRate
	feeRate.GasRate = req.GasRate

	if feeRate.ID == 0 {
		err = s.userRepo.CreateFeeRate(feeRate)
	} else {
		err = s.userRepo.UpdateFeeRate(feeRate)
	}

	if err != nil {
		return nil, err
	}

	return feeRate, nil
}
