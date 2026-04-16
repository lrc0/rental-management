package service

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"rental-management/internal/model"
	"rental-management/internal/pkg/utils"
	"rental-management/internal/repository"
)

type AuthService struct {
	userRepo     *repository.UserRepository
	propertyRepo *repository.PropertyRepository
	roomRepo     *repository.RoomRepository
	tenantRepo   *repository.TenantRepository
	billRepo     *repository.BillRepository
}

func NewAuthService(
	userRepo *repository.UserRepository,
	propertyRepo *repository.PropertyRepository,
	roomRepo *repository.RoomRepository,
	tenantRepo *repository.TenantRepository,
	billRepo *repository.BillRepository,
) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		propertyRepo: propertyRepo,
		roomRepo:     roomRepo,
		tenantRepo:   tenantRepo,
		billRepo:     billRepo,
	}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}

// Register 用户注册
func (s *AuthService) Register(req *RegisterRequest) (*model.User, error) {
	// 检查手机号是否已注册
	exists, err := s.userRepo.ExistsByPhone(req.Phone)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("手机号已注册")
	}

	// 加密密码
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// 创建用户
	user := &model.User{
		Phone:        req.Phone,
		PasswordHash: passwordHash,
		Name:         req.Name,
		Status:       1,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// 创建默认费率配置
	feeRate := &model.FeeRate{
		UserID:          user.ID,
		WaterRate:       5.0,   // 默认水费5元/吨
		ElectricityRate: 0.6,   // 默认电费0.6元/度
		GasRate:         3.0,   // 默认气费3元/立方
	}
	_ = s.userRepo.CreateFeeRate(feeRate)

	return user, nil
}

// Login 用户登录
func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	// 查找用户
	user, err := s.userRepo.FindByPhone(req.Phone)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("手机号未注册")
		}
		return nil, err
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("密码错误")
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, errors.New("账号已被禁用")
	}

	return &LoginResponse{
		User: user,
	}, nil
}

// GetProfile 获取用户信息
func (s *AuthService) GetProfile(userID uint) (*model.User, error) {
	return s.userRepo.FindByID(userID)
}

// UpdateProfileRequest 更新用户信息请求
type UpdateProfileRequest struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

// UpdateProfile 更新用户信息
func (s *AuthService) UpdateProfile(userID uint, req *UpdateProfileRequest) (*model.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// ChangePassword 修改密码
func (s *AuthService) ChangePassword(userID uint, req *ChangePasswordRequest) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	// 验证旧密码
	if !utils.CheckPassword(req.OldPassword, user.PasswordHash) {
		return errors.New("原密码错误")
	}

	// 加密新密码
	passwordHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = passwordHash
	return s.userRepo.Update(user)
}

// StatisticsResponse 统计响应
type StatisticsResponse struct {
	PropertyCount int64   `json:"property_count"`
	RoomCount     int64   `json:"room_count"`
	TenantCount   int64   `json:"tenant_count"`
	MonthlyIncome float64 `json:"monthly_income"`
}

// GetStatistics 获取统计数据
func (s *AuthService) GetStatistics(userID uint, startDate, endDate time.Time) (*StatisticsResponse, error) {
	// 获取房源数
	propertyCount, _ := s.propertyRepo.CountByUserID(userID)

	// 获取房间数
	roomCount, _ := s.roomRepo.CountByUserID(userID)

	// 获取租客数
	tenantCount, _ := s.tenantRepo.CountByUserID(userID)

	// 获取本月收入
	stats, _ := s.billRepo.GetBillStatistics(userID, startDate, endDate)
	var monthlyIncome float64
	if stats != nil {
		if v, ok := stats["paid_amount"].(float64); ok {
			monthlyIncome = v
		}
	}

	return &StatisticsResponse{
		PropertyCount: propertyCount,
		RoomCount:     roomCount,
		TenantCount:   tenantCount,
		MonthlyIncome: monthlyIncome,
	}, nil
}
