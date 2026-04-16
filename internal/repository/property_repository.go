package repository

import (
	"gorm.io/gorm"
	"rental-management/internal/model"
)

type PropertyRepository struct {
	db *gorm.DB
}

func NewPropertyRepository(db *gorm.DB) *PropertyRepository {
	return &PropertyRepository{db: db}
}

// Create 创建房源
func (r *PropertyRepository) Create(property *model.Property) error {
	return r.db.Create(property).Error
}

// FindByID 根据ID查找房源
func (r *PropertyRepository) FindByID(id uint) (*model.Property, error) {
	var property model.Property
	err := r.db.First(&property, id).Error
	if err != nil {
		return nil, err
	}
	return &property, nil
}

// FindByIDAndUserID 根据ID和用户ID查找房源
func (r *PropertyRepository) FindByIDAndUserID(id, userID uint) (*model.Property, error) {
	var property model.Property
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&property).Error
	if err != nil {
		return nil, err
	}
	return &property, nil
}

// List 获取房源列表
func (r *PropertyRepository) List(userID uint, page, pageSize int) ([]model.Property, int64, error) {
	var properties []model.Property
	var total int64

	db := r.db.Model(&model.Property{}).Where("user_id = ?", userID)
	db.Count(&total)

	offset := (page - 1) * pageSize
	err := db.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&properties).Error
	return properties, total, err
}

// Update 更新房源
func (r *PropertyRepository) Update(property *model.Property) error {
	return r.db.Save(property).Error
}

// Delete 删除房源
func (r *PropertyRepository) Delete(id uint) error {
	return r.db.Delete(&model.Property{}, id).Error
}

// UpdateRoomCount 更新房间数
func (r *PropertyRepository) UpdateRoomCount(propertyID uint, count int) error {
	return r.db.Model(&model.Property{}).Where("id = ?", propertyID).
		Update("total_rooms", count).Error
}

// CountByUserID 统计用户的房源数
func (r *PropertyRepository) CountByUserID(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Property{}).Where("user_id = ? AND deleted_at IS NULL", userID).Count(&count).Error
	return count, err
}
