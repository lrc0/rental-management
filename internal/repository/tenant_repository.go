package repository

import (
	"gorm.io/gorm"
	"rental-management/internal/model"
)

type TenantRepository struct {
	db *gorm.DB
}

func NewTenantRepository(db *gorm.DB) *TenantRepository {
	return &TenantRepository{db: db}
}

// Create 创建租客
func (r *TenantRepository) Create(tenant *model.Tenant) error {
	return r.db.Create(tenant).Error
}

// FindByID 根据ID查找租客
func (r *TenantRepository) FindByID(id uint) (*model.Tenant, error) {
	var tenant model.Tenant
	err := r.db.First(&tenant, id).Error
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

// FindByIDAndUserID 根据ID和用户ID查找租客
func (r *TenantRepository) FindByIDAndUserID(id, userID uint) (*model.Tenant, error) {
	var tenant model.Tenant
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&tenant).Error
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

// TenantWithRoom 租客带房间信息
type TenantWithRoom struct {
	model.Tenant
	RoomID      *uint  `json:"room_id"`
	RoomNumber  string `json:"room_number"`
	PropertyName string `json:"property_name"`
}

// ListWithRoom 获取租客列表（包含当前房间信息）
func (r *TenantRepository) ListWithRoom(userID uint, status *int8, page, pageSize int) ([]TenantWithRoom, int64, error) {
	var results []TenantWithRoom
	var total int64

	// 统计总数
	countDB := r.db.Model(&model.Tenant{}).Where("user_id = ?", userID)
	if status != nil {
		countDB = countDB.Where("status = ?", *status)
	}
	countDB.Count(&total)

	// 查询租客及其当前房间
	offset := (page - 1) * pageSize
	query := `
		SELECT t.*, c.room_id, r.room_number, p.name as property_name
		FROM tenants t
		LEFT JOIN contracts c ON t.id = c.tenant_id AND c.status = 1 AND c.deleted_at IS NULL
		LEFT JOIN rooms r ON c.room_id = r.id AND r.deleted_at IS NULL
		LEFT JOIN properties p ON r.property_id = p.id AND p.deleted_at IS NULL
		WHERE t.user_id = ? AND t.deleted_at IS NULL
	`
	args := []interface{}{userID}
	if status != nil {
		query += " AND t.status = ?"
		args = append(args, *status)
	}
	query += " ORDER BY t.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, pageSize, offset)

	err := r.db.Raw(query, args...).Scan(&results).Error
	return results, total, err
}

// List 获取租客列表
func (r *TenantRepository) List(userID uint, status *int8, page, pageSize int) ([]model.Tenant, int64, error) {
	var tenants []model.Tenant
	var total int64

	db := r.db.Model(&model.Tenant{}).Where("user_id = ?", userID)
	if status != nil {
		db = db.Where("status = ?", *status)
	}

	db.Count(&total)

	offset := (page - 1) * pageSize
	err := db.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&tenants).Error
	return tenants, total, err
}

// Update 更新租客
func (r *TenantRepository) Update(tenant *model.Tenant) error {
	return r.db.Save(tenant).Error
}

// Delete 删除租客
func (r *TenantRepository) Delete(id uint) error {
	return r.db.Delete(&model.Tenant{}, id).Error
}

// FindByIDCard 根据身份证查找租客
func (r *TenantRepository) FindByIDCard(userID uint, idCard string) (*model.Tenant, error) {
	var tenant model.Tenant
	err := r.db.Where("user_id = ? AND id_card = ?", userID, idCard).First(&tenant).Error
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

// CountByUserID 统计用户的租客数
func (r *TenantRepository) CountByUserID(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Tenant{}).Where("user_id = ? AND deleted_at IS NULL", userID).Count(&count).Error
	return count, err
}
