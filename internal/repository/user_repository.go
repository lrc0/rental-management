package repository

import (
	"gorm.io/gorm"
	"rental-management/internal/model"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create 创建用户
func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// FindByPhone 根据手机号查找用户
func (r *UserRepository) FindByPhone(phone string) (*model.User, error) {
	var user model.User
	err := r.db.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername 根据用户名查找用户
func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID 根据ID查找用户
func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户
func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// ExistsByPhone 检查手机号是否已存在
func (r *UserRepository) ExistsByPhone(phone string) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("phone = ?", phone).Count(&count).Error
	return count > 0, err
}

// ExistsByUsername 检查用户名是否已存在
func (r *UserRepository) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

// CreateFeeRate 创建费率配置
func (r *UserRepository) CreateFeeRate(feeRate *model.FeeRate) error {
	return r.db.Create(feeRate).Error
}

// GetFeeRateByUserID 获取用户费率配置
func (r *UserRepository) GetFeeRateByUserID(userID uint) (*model.FeeRate, error) {
	var feeRate model.FeeRate
	err := r.db.Where("user_id = ?", userID).First(&feeRate).Error
	if err != nil {
		return nil, err
	}
	return &feeRate, nil
}

// UpdateFeeRate 更新费率配置
func (r *UserRepository) UpdateFeeRate(feeRate *model.FeeRate) error {
	return r.db.Save(feeRate).Error
}
