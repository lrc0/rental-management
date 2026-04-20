package repository

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"rental-management/internal/model"
)

type BillRepository struct {
	db *gorm.DB
}

func NewBillRepository(db *gorm.DB) *BillRepository {
	return &BillRepository{db: db}
}

// CreateMeterReading 创建抄表记录
func (r *BillRepository) CreateMeterReading(reading *model.MeterReading) error {
	return r.db.Create(reading).Error
}

// GetLatestMeterReading 获取房间最新抄表记录
func (r *BillRepository) GetLatestMeterReading(roomID uint, beforeDate time.Time) (*model.MeterReading, error) {
	var reading model.MeterReading
	err := r.db.Where("room_id = ? AND reading_date < ?", roomID, beforeDate).
		Order("reading_date DESC").First(&reading).Error
	if err != nil {
		return nil, err
	}
	return &reading, nil
}

// ListMeterReadings 获取抄表记录列表
func (r *BillRepository) ListMeterReadings(userID uint, roomID uint, startDate, endDate *time.Time, page, pageSize int) ([]model.MeterReading, int64, error) {
	var readings []model.MeterReading
	var total int64

	db := r.db.Model(&model.MeterReading{}).
		Joins("JOIN rooms ON rooms.id = meter_readings.room_id").
		Joins("JOIN properties ON properties.id = rooms.property_id").
		Where("properties.user_id = ?", userID)

	if roomID > 0 {
		db = db.Where("meter_readings.room_id = ?", roomID)
	}
	if startDate != nil {
		db = db.Where("meter_readings.reading_date >= ?", startDate)
	}
	if endDate != nil {
		db = db.Where("meter_readings.reading_date <= ?", endDate)
	}

	db.Count(&total)

	offset := (page - 1) * pageSize
	err := db.Offset(offset).Limit(pageSize).
		Preload("Room").
		Order("meter_readings.reading_date DESC").
		Find(&readings).Error
	return readings, total, err
}

// CreateBill 创建账单
func (r *BillRepository) CreateBill(bill *model.Bill) error {
	return r.db.Create(bill).Error
}

// FindBillByID 根据ID查找账单
func (r *BillRepository) FindBillByID(id uint) (*model.Bill, error) {
	var bill model.Bill
	err := r.db.Preload("Room").Preload("Tenant").First(&bill, id).Error
	if err != nil {
		return nil, err
	}
	return &bill, nil
}

// FindBillByIDAndUserID 根据ID和用户ID查找账单
func (r *BillRepository) FindBillByIDAndUserID(id, userID uint) (*model.Bill, error) {
	var bill model.Bill
	err := r.db.Where("id = ? AND user_id = ?", id, userID).
		Preload("Room").Preload("Tenant").First(&bill).Error
	if err != nil {
		return nil, err
	}
	return &bill, nil
}

// ListBills 获取账单列表
func (r *BillRepository) ListBills(userID uint, roomID uint, status *int8, billMonth string, page, pageSize int) ([]model.Bill, int64, error) {
	var bills []model.Bill
	var total int64

	db := r.db.Model(&model.Bill{}).Where("user_id = ?", userID)

	if roomID > 0 {
		db = db.Where("room_id = ?", roomID)
	}
	if status != nil {
		db = db.Where("status = ?", *status)
	}
	if billMonth != "" {
		db = db.Where("bill_month = ?", billMonth)
	}

	db.Count(&total)

	offset := (page - 1) * pageSize
	err := db.Offset(offset).Limit(pageSize).
		Preload("Room").Preload("Tenant").
		Order("created_at DESC").
		Find(&bills).Error
	return bills, total, err
}

// UpdateBill 更新账单
func (r *BillRepository) UpdateBill(bill *model.Bill) error {
	return r.db.Save(bill).Error
}

// CreatePayment 创建收款记录
func (r *BillRepository) CreatePayment(payment *model.Payment) error {
	return r.db.Create(payment).Error
}

// GetPaymentsByBillID 获取账单的收款记录
func (r *BillRepository) GetPaymentsByBillID(billID uint) ([]model.Payment, error) {
	var payments []model.Payment
	err := r.db.Where("bill_id = ?", billID).Order("paid_at DESC").Find(&payments).Error
	return payments, err
}

// GetBillStatistics 获取账单统计
func (r *BillRepository) GetBillStatistics(userID uint, startDate, endDate time.Time) (map[string]interface{}, error) {
	var totalAmount float64
	var paidAmount float64
	var pendingAmount float64
	var billCount int64
	var paidCount int64

	r.db.Model(&model.Bill{}).Where("user_id = ? AND created_at >= ? AND created_at <= ?", userID, startDate, endDate).
		Count(&billCount)

	r.db.Model(&model.Bill{}).Where("user_id = ? AND created_at >= ? AND created_at <= ? AND status = ?", userID, startDate, endDate, model.BillStatusPaid).
		Select("COALESCE(SUM(amount), 0)").Scan(&paidAmount)

	r.db.Model(&model.Bill{}).Where("user_id = ? AND created_at >= ? AND created_at <= ? AND status = ?", userID, startDate, endDate, model.BillStatusPending).
		Select("COALESCE(SUM(amount), 0)").Scan(&pendingAmount)

	r.db.Model(&model.Bill{}).Where("user_id = ? AND created_at >= ? AND created_at <= ?", userID, startDate, endDate).
		Select("COALESCE(SUM(amount), 0)").Scan(&totalAmount)

	r.db.Model(&model.Bill{}).Where("user_id = ? AND created_at >= ? AND created_at <= ? AND status = ?", userID, startDate, endDate, model.BillStatusPaid).
		Count(&paidCount)

	return map[string]interface{}{
		"total_amount":   totalAmount,
		"paid_amount":    paidAmount,
		"pending_amount": pendingAmount,
		"bill_count":     billCount,
		"paid_count":     paidCount,
	}, nil
}

// GetMonthlyStatistics 获取月度统计
func (r *BillRepository) GetMonthlyStatistics(userID uint, year int) ([]map[string]interface{}, error) {
	type MonthStat struct {
		Month     string  `json:"month"`
		TotalFee  float64 `json:"total_fee"`
		PaidFee   float64 `json:"paid_fee"`
		BillCount int64   `json:"bill_count"`
	}

	var stats []MonthStat
	err := r.db.Model(&model.Bill{}).
		Select("bill_month as month, SUM(amount) as total_fee, SUM(CASE WHEN status = ? THEN amount ELSE 0 END) as paid_fee, COUNT(*) as bill_count", model.BillStatusPaid).
		Where("user_id = ? AND bill_month LIKE ?", userID, fmt.Sprintf("%d-%%", year)).
		Group("bill_month").
		Order("bill_month").
		Scan(&stats).Error

	result := make([]map[string]interface{}, len(stats))
	for i, s := range stats {
		result[i] = map[string]interface{}{
			"month":      s.Month,
			"total_fee":  s.TotalFee,
			"paid_fee":   s.PaidFee,
			"bill_count": s.BillCount,
		}
	}

	return result, err
}

// CheckBillExists 检查账单是否已存在
func (r *BillRepository) CheckBillExists(roomID uint, billMonth string) (bool, error) {
	var count int64
	err := r.db.Model(&model.Bill{}).Where("room_id = ? AND bill_month = ?", roomID, billMonth).Count(&count).Error
	return count > 0, err
}

// DeleteBill 删除账单
func (r *BillRepository) DeleteBill(id uint) error {
	return r.db.Delete(&model.Bill{}, id).Error
}

// DeleteBillByIDAndUserID 根据ID和用户ID删除账单
func (r *BillRepository) DeleteBillByIDAndUserID(id, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Bill{}).Error
}

// GetMeterReadingByID 根据ID获取抄表记录
func (r *BillRepository) GetMeterReadingByID(id uint) (*model.MeterReading, error) {
	var reading model.MeterReading
	err := r.db.First(&reading, id).Error
	if err != nil {
		return nil, err
	}
	return &reading, nil
}

// DeleteMeterReading 删除抄表记录
func (r *BillRepository) DeleteMeterReading(id uint) error {
	return r.db.Delete(&model.MeterReading{}, id).Error
}

// DeleteMeterReadingByIDAndUserID 根据ID和用户ID删除抄表记录
func (r *BillRepository) DeleteMeterReadingByIDAndUserID(id, userID uint) error {
	return r.db.Exec(`
		DELETE m FROM meter_readings m
		JOIN rooms r ON m.room_id = r.id
		JOIN properties p ON r.property_id = p.id
		WHERE m.id = ? AND p.user_id = ?
	`, id, userID).Error
}
