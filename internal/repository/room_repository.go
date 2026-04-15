package repository

import (
	"gorm.io/gorm"
	"rental-management/internal/model"
)

type RoomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) *RoomRepository {
	return &RoomRepository{db: db}
}

// Create 创建房间
func (r *RoomRepository) Create(room *model.Room) error {
	return r.db.Create(room).Error
}

// FindByID 根据ID查找房间
func (r *RoomRepository) FindByID(id uint) (*model.Room, error) {
	var room model.Room
	err := r.db.Preload("Property").First(&room, id).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

// FindByIDAndUserID 根据ID和用户ID查找房间(通过房源关联)
func (r *RoomRepository) FindByIDAndUserID(id, userID uint) (*model.Room, error) {
	var room model.Room
	err := r.db.Joins("JOIN properties ON properties.id = rooms.property_id").
		Where("rooms.id = ? AND properties.user_id = ?", id, userID).
		Preload("Property").
		First(&room).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

// List 获取房间列表
func (r *RoomRepository) List(userID uint, propertyID uint, status *int8, page, pageSize int) ([]model.Room, int64, error) {
	var rooms []model.Room
	var total int64

	db := r.db.Model(&model.Room{}).Joins("JOIN properties ON properties.id = rooms.property_id").
		Where("properties.user_id = ?", userID)

	if propertyID > 0 {
		db = db.Where("rooms.property_id = ?", propertyID)
	}
	if status != nil {
		db = db.Where("rooms.status = ?", *status)
	}

	db.Count(&total)

	offset := (page - 1) * pageSize
	err := db.Offset(offset).Limit(pageSize).
		Preload("Property").
		Order("rooms.created_at DESC").
		Find(&rooms).Error
	return rooms, total, err
}

// Update 更新房间
func (r *RoomRepository) Update(room *model.Room) error {
	return r.db.Save(room).Error
}

// Delete 删除房间
func (r *RoomRepository) Delete(id uint) error {
	return r.db.Delete(&model.Room{}, id).Error
}

// CountByPropertyID 统计房源下的房间数
func (r *RoomRepository) CountByPropertyID(propertyID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Room{}).Where("property_id = ?", propertyID).Count(&count).Error
	return count, err
}

// FindByPropertyID 根据房源ID查找房间
func (r *RoomRepository) FindByPropertyID(propertyID uint) ([]model.Room, error) {
	var rooms []model.Room
	err := r.db.Where("property_id = ?", propertyID).Find(&rooms).Error
	return rooms, err
}

// UpdateStatus 更新房间状态
func (r *RoomRepository) UpdateStatus(id uint, status int8) error {
	return r.db.Model(&model.Room{}).Where("id = ?", id).Update("status", status).Error
}

// RoomWithTenant 房间带租客信息
type RoomWithTenant struct {
	model.Room
	TenantID   *uint  `json:"tenant_id"`
	TenantName string `json:"tenant_name"`
	TenantPhone string `json:"tenant_phone"`
}

// ListWithTenant 获取房间列表（包含当前租客信息）
func (r *RoomRepository) ListWithTenant(userID uint, propertyID uint, status *int8, page, pageSize int) ([]RoomWithTenant, int64, error) {
	var results []RoomWithTenant
	var total int64

	// 统计总数
	countDB := r.db.Model(&model.Room{}).Joins("JOIN properties ON properties.id = rooms.property_id").
		Where("properties.user_id = ?", userID)
	if propertyID > 0 {
		countDB = countDB.Where("rooms.property_id = ?", propertyID)
	}
	if status != nil {
		countDB = countDB.Where("rooms.status = ?", *status)
	}
	countDB.Count(&total)

	// 查询房间及其当前租客
	offset := (page - 1) * pageSize
	query := `
		SELECT r.*, p.id as property_id, p.name as property_name, p.user_id,
			c.tenant_id, t.name as tenant_name, t.phone as tenant_phone
		FROM rooms r
		JOIN properties p ON r.property_id = p.id AND p.deleted_at IS NULL
		LEFT JOIN contracts c ON r.id = c.room_id AND c.status = 1 AND c.deleted_at IS NULL
		LEFT JOIN tenants t ON c.tenant_id = t.id AND t.deleted_at IS NULL
		WHERE p.user_id = ? AND r.deleted_at IS NULL
	`
	args := []interface{}{userID}
	if propertyID > 0 {
		query += " AND r.property_id = ?"
		args = append(args, propertyID)
	}
	if status != nil {
		query += " AND r.status = ?"
		args = append(args, *status)
	}
	query += " ORDER BY r.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, pageSize, offset)

	err := r.db.Raw(query, args...).Scan(&results).Error
	return results, total, err
}
