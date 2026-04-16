package model

import (
	"time"

	"gorm.io/gorm"
)

// Tenant 租客表
type Tenant struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	UserID           uint           `gorm:"index;not null" json:"user_id"` // 所属房东
	Name             string         `gorm:"size:50;not null" json:"name"`
	Phone            string         `gorm:"size:20" json:"phone"`
	IDCard           string         `gorm:"size:18" json:"id_card"`            // 身份证
	Gender           int8           `gorm:"default:0;comment:0未知 1男 2女" json:"gender"`
	EmergencyContact string         `gorm:"size:50" json:"emergency_contact"`   // 紧急联系人
	EmergencyPhone   string         `gorm:"size:20" json:"emergency_phone"`     // 紧急联系电话
	Status           int8           `gorm:"default:1;comment:1正常 2已退租" json:"status"`
	Remark           string         `gorm:"size:255" json:"remark"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Tenant) TableName() string {
	return "tenants"
}

// 租客状态
const (
	TenantStatusActive    int8 = 1 // 正常
	TenantStatusCheckout  int8 = 2 // 已退租
)

// Contract 合同表
type Contract struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	RoomID          uint           `gorm:"index;not null" json:"room_id"`
	TenantID        uint           `gorm:"index;not null" json:"tenant_id"`
	StartDate       time.Time      `gorm:"type:date" json:"start_date"`
	EndDate         time.Time      `gorm:"type:date" json:"end_date"`
	MonthlyRent     float64        `gorm:"type:decimal(10,2)" json:"monthly_rent"`
	Deposit         float64        `gorm:"type:decimal(10,2)" json:"deposit"` // 押金
	PaymentDay      int8           `gorm:"comment:每月几号交租" json:"payment_day"`
	Status          int8           `gorm:"default:1;comment:1生效 2已到期 3已解约" json:"status"`
	TerminateReason string         `gorm:"size:255" json:"terminate_reason"` // 解约原因
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联 - 不创建外键约束
	Room    *Room   `gorm:"foreignKey:RoomID;references:ID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION" json:"room,omitempty"`
	Tenant  *Tenant `gorm:"foreignKey:TenantID;references:ID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION" json:"tenant,omitempty"`
}

func (Contract) TableName() string {
	return "contracts"
}

// 合同状态
const (
	ContractStatusActive    int8 = 1 // 生效
	ContractStatusExpired   int8 = 2 // 已到期
	ContractStatusTerminated int8 = 3 // 已解约
)
