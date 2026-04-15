package service

import (
	"errors"
	"time"

	"rental-management/internal/model"
	"rental-management/internal/repository"
)

type TenantService struct {
	tenantRepo   *repository.TenantRepository
	contractRepo *repository.ContractRepository
	roomRepo     *repository.RoomRepository
	userRepo     *repository.UserRepository
}

func NewTenantService(
	tenantRepo *repository.TenantRepository,
	contractRepo *repository.ContractRepository,
	roomRepo *repository.RoomRepository,
	userRepo *repository.UserRepository,
) *TenantService {
	return &TenantService{
		tenantRepo:   tenantRepo,
		contractRepo: contractRepo,
		roomRepo:     roomRepo,
		userRepo:     userRepo,
	}
}

// CreateTenantRequest 创建租客请求
type CreateTenantRequest struct {
	Name             string `json:"name" binding:"required"`
	Phone            string `json:"phone"`
	IDCard           string `json:"id_card"`
	Gender           int8   `json:"gender"`
	EmergencyContact string `json:"emergency_contact"`
	EmergencyPhone   string `json:"emergency_phone"`
	Remark           string `json:"remark"`
}

// UpdateTenantRequest 更新租客请求
type UpdateTenantRequest struct {
	Name             string `json:"name"`
	Phone            string `json:"phone"`
	IDCard           string `json:"id_card"`
	Gender           *int8  `json:"gender"`
	EmergencyContact string `json:"emergency_contact"`
	EmergencyPhone   string `json:"emergency_phone"`
	Remark           string `json:"remark"`
}

// CreateTenant 创建租客
func (s *TenantService) CreateTenant(userID uint, req *CreateTenantRequest) (*model.Tenant, error) {
	// 检查身份证是否已存在
	if req.IDCard != "" {
		if _, err := s.tenantRepo.FindByIDCard(userID, req.IDCard); err == nil {
			return nil, errors.New("该身份证号已存在")
		}
	}

	tenant := &model.Tenant{
		UserID:           userID,
		Name:             req.Name,
		Phone:            req.Phone,
		IDCard:           req.IDCard,
		Gender:           req.Gender,
		EmergencyContact: req.EmergencyContact,
		EmergencyPhone:   req.EmergencyPhone,
		Remark:           req.Remark,
		Status:           model.TenantStatusActive,
	}

	if err := s.tenantRepo.Create(tenant); err != nil {
		return nil, err
	}

	return tenant, nil
}

// GetTenant 获取租客详情
func (s *TenantService) GetTenant(id, userID uint) (*model.Tenant, error) {
	return s.tenantRepo.FindByIDAndUserID(id, userID)
}

// ListTenants 获取租客列表
func (s *TenantService) ListTenants(userID uint, status *int8, page, pageSize int) ([]repository.TenantWithRoom, int64, error) {
	return s.tenantRepo.ListWithRoom(userID, status, page, pageSize)
}

// UpdateTenant 更新租客
func (s *TenantService) UpdateTenant(id, userID uint, req *UpdateTenantRequest) (*model.Tenant, error) {
	tenant, err := s.tenantRepo.FindByIDAndUserID(id, userID)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		tenant.Name = req.Name
	}
	if req.Phone != "" {
		tenant.Phone = req.Phone
	}
	if req.IDCard != "" {
		tenant.IDCard = req.IDCard
	}
	if req.Gender != nil {
		tenant.Gender = *req.Gender
	}
	if req.EmergencyContact != "" {
		tenant.EmergencyContact = req.EmergencyContact
	}
	if req.EmergencyPhone != "" {
		tenant.EmergencyPhone = req.EmergencyPhone
	}
	if req.Remark != "" {
		tenant.Remark = req.Remark
	}

	if err := s.tenantRepo.Update(tenant); err != nil {
		return nil, err
	}

	return tenant, nil
}

// DeleteTenant 删除租客
func (s *TenantService) DeleteTenant(id, userID uint) error {
	tenant, err := s.tenantRepo.FindByIDAndUserID(id, userID)
	if err != nil {
		return err
	}

	// 检查是否有有效合同
	if _, err := s.contractRepo.FindActiveByTenantID(tenant.ID); err == nil {
		return errors.New("该租客有有效合同，无法删除")
	}

	return s.tenantRepo.Delete(id)
}

// CreateContractRequest 签订合同请求
type CreateContractRequest struct {
	RoomID      uint    `json:"room_id" binding:"required"`
	TenantID    uint    `json:"tenant_id" binding:"required"`
	StartDate   string  `json:"start_date" binding:"required"`
	EndDate     string  `json:"end_date" binding:"required"`
	MonthlyRent float64 `json:"monthly_rent" binding:"required"`
	Deposit     float64 `json:"deposit"`
	PaymentDay  int8    `json:"payment_day"`
}

// TerminateContractRequest 解约请求
type TerminateContractRequest struct {
	Reason string `json:"reason"`
}

// CreateContract 签订合同
func (s *TenantService) CreateContract(userID uint, req *CreateContractRequest) (*model.Contract, error) {
	// 解析日期
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, errors.New("开始日期格式错误")
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, errors.New("结束日期格式错误")
	}

	// 验证房间归属
	room, err := s.roomRepo.FindByIDAndUserID(req.RoomID, userID)
	if err != nil {
		return nil, errors.New("房间不存在或无权限")
	}

	// 检查房间是否已被租
	available, err := s.contractRepo.CheckRoomAvailable(req.RoomID, startDate)
	if err != nil {
		return nil, err
	}
	if !available {
		return nil, errors.New("该房间已被出租")
	}

	// 验证租客归属
	tenant, err := s.tenantRepo.FindByIDAndUserID(req.TenantID, userID)
	if err != nil {
		return nil, errors.New("租客不存在或无权限")
	}

	// 检查租客是否有有效合同
	if _, err := s.contractRepo.FindActiveByTenantID(tenant.ID); err == nil {
		return nil, errors.New("该租客已有有效合同")
	}

	// 验证日期
	if endDate.Before(startDate) {
		return nil, errors.New("结束日期不能早于开始日期")
	}

	contract := &model.Contract{
		RoomID:      req.RoomID,
		TenantID:    req.TenantID,
		StartDate:   startDate,
		EndDate:     endDate,
		MonthlyRent: req.MonthlyRent,
		Deposit:     req.Deposit,
		PaymentDay:  req.PaymentDay,
		Status:      model.ContractStatusActive,
	}

	if err := s.contractRepo.Create(contract); err != nil {
		return nil, err
	}

	// 更新房间状态为已租
	_ = s.roomRepo.UpdateStatus(room.ID, model.RoomStatusRented)

	return contract, nil
}

// GetContract 获取合同详情
func (s *TenantService) GetContract(id, userID uint) (*model.Contract, error) {
	return s.contractRepo.FindByIDAndUserID(id, userID)
}

// ListContracts 获取合同列表
func (s *TenantService) ListContracts(userID uint, status *int8, page, pageSize int) ([]model.Contract, int64, error) {
	return s.contractRepo.List(userID, status, page, pageSize)
}

// TerminateContract 解约
func (s *TenantService) TerminateContract(id, userID uint, req *TerminateContractRequest) error {
	contract, err := s.contractRepo.FindByIDAndUserID(id, userID)
	if err != nil {
		return err
	}

	if contract.Status != model.ContractStatusActive {
		return errors.New("合同已失效")
	}

	contract.Status = model.ContractStatusTerminated
	contract.TerminateReason = req.Reason

	if err := s.contractRepo.Update(contract); err != nil {
		return err
	}

	// 更新房间状态为空置
	_ = s.roomRepo.UpdateStatus(contract.RoomID, model.RoomStatusVacant)

	// 更新租客状态
	tenant, _ := s.tenantRepo.FindByID(contract.TenantID)
	if tenant != nil {
		tenant.Status = model.TenantStatusCheckout
		_ = s.tenantRepo.Update(tenant)
	}

	return nil
}
