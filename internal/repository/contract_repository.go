package repository

import (
	"time"

	"gorm.io/gorm"
	"rental-management/internal/model"
)

type ContractRepository struct {
	db *gorm.DB
}

func NewContractRepository(db *gorm.DB) *ContractRepository {
	return &ContractRepository{db: db}
}

// Create 创建合同
func (r *ContractRepository) Create(contract *model.Contract) error {
	return r.db.Create(contract).Error
}

// FindByID 根据ID查找合同
func (r *ContractRepository) FindByID(id uint) (*model.Contract, error) {
	var contract model.Contract
	err := r.db.Preload("Room").Preload("Tenant").First(&contract, id).Error
	if err != nil {
		return nil, err
	}
	return &contract, nil
}

// FindByIDAndUserID 根据ID和用户ID查找合同
func (r *ContractRepository) FindByIDAndUserID(id, userID uint) (*model.Contract, error) {
	var contract model.Contract
	err := r.db.Joins("JOIN rooms ON rooms.id = contracts.room_id").
		Joins("JOIN properties ON properties.id = rooms.property_id").
		Where("contracts.id = ? AND properties.user_id = ?", id, userID).
		Preload("Room").Preload("Tenant").
		First(&contract).Error
	if err != nil {
		return nil, err
	}
	return &contract, nil
}

// List 获取合同列表
func (r *ContractRepository) List(userID uint, status *int8, page, pageSize int) ([]model.Contract, int64, error) {
	var contracts []model.Contract
	var total int64

	db := r.db.Model(&model.Contract{}).
		Joins("JOIN rooms ON rooms.id = contracts.room_id").
		Joins("JOIN properties ON properties.id = rooms.property_id").
		Where("properties.user_id = ?", userID)

	if status != nil {
		db = db.Where("contracts.status = ?", *status)
	}

	db.Count(&total)

	offset := (page - 1) * pageSize
	err := db.Offset(offset).Limit(pageSize).
		Preload("Room").Preload("Tenant").
		Order("contracts.created_at DESC").
		Find(&contracts).Error
	return contracts, total, err
}

// Update 更新合同
func (r *ContractRepository) Update(contract *model.Contract) error {
	return r.db.Save(contract).Error
}

// FindActiveByRoomID 查找房间当前有效合同
func (r *ContractRepository) FindActiveByRoomID(roomID uint) (*model.Contract, error) {
	var contract model.Contract
	err := r.db.Where("room_id = ? AND status = ?", roomID, model.ContractStatusActive).
		First(&contract).Error
	if err != nil {
		return nil, err
	}
	return &contract, nil
}

// FindActiveByTenantID 查找租客当前有效合同
func (r *ContractRepository) FindActiveByTenantID(tenantID uint) (*model.Contract, error) {
	var contract model.Contract
	err := r.db.Where("tenant_id = ? AND status = ?", tenantID, model.ContractStatusActive).
		First(&contract).Error
	if err != nil {
		return nil, err
	}
	return &contract, nil
}

// FindByRoomID 根据房间ID查找合同
func (r *ContractRepository) FindByRoomID(roomID uint) ([]model.Contract, error) {
	var contracts []model.Contract
	err := r.db.Where("room_id = ?", roomID).Preload("Tenant").Find(&contracts).Error
	return contracts, err
}

// CheckRoomAvailable 检查房间是否可租
func (r *ContractRepository) CheckRoomAvailable(roomID uint, startDate time.Time) (bool, error) {
	var count int64
	err := r.db.Model(&model.Contract{}).
		Where("room_id = ? AND status = ? AND end_date >= ?", roomID, model.ContractStatusActive, startDate).
		Count(&count).Error
	return count == 0, err
}
