package model

import (
	"time"
)

// MeterReading 抄表记录表
type MeterReading struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	RoomID              uint      `gorm:"index:idx_room_date;not null" json:"room_id"`
	ReadingDate         time.Time `gorm:"type:date;index:idx_room_date" json:"reading_date"`
	WaterReading        float64   `gorm:"type:decimal(10,2)" json:"water_reading"`        // 水表读数
	ElectricityReading  float64   `gorm:"type:decimal(10,2)" json:"electricity_reading"`  // 电表读数
	GasReading          float64   `gorm:"type:decimal(10,2)" json:"gas_reading"`          // 气表读数
	WaterUsage          float64   `gorm:"type:decimal(10,2)" json:"water_usage"`          // 用水量(计算值)
	ElectricityUsage    float64   `gorm:"type:decimal(10,2)" json:"electricity_usage"`    // 用电量(计算值)
	GasUsage            float64   `gorm:"type:decimal(10,2)" json:"gas_usage"`            // 用气量(计算值)
	Remark              string    `gorm:"size:255" json:"remark"`
	CreatedAt           time.Time `json:"created_at"`

	// 关联
	Room *Room `gorm:"foreignKey:RoomID" json:"room,omitempty"`
}

func (MeterReading) TableName() string {
	return "meter_readings"
}

// Bill 账单表
type Bill struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	UserID           uint       `gorm:"index:idx_user_month;not null" json:"user_id"`
	RoomID           uint       `gorm:"index" json:"room_id"`
	TenantID         uint       `gorm:"index" json:"tenant_id"`
	BillType         int8       `gorm:"comment:1租金 2水费 3电费 4气费 5综合账单" json:"bill_type"`
	BillMonth        string     `gorm:"size:7;index:idx_user_month" json:"bill_month"` // 2024-01
	Amount           float64    `gorm:"type:decimal(10,2)" json:"amount"`              // 总金额
	WaterFee         float64    `gorm:"type:decimal(10,2)" json:"water_fee"`
	ElectricityFee   float64    `gorm:"type:decimal(10,2)" json:"electricity_fee"`
	GasFee           float64    `gorm:"type:decimal(10,2)" json:"gas_fee"`
	RentFee          float64    `gorm:"type:decimal(10,2)" json:"rent_fee"`
	Status           int8       `gorm:"default:1;comment:1待支付 2已支付 3已逾期" json:"status"`
	DueDate          *time.Time `gorm:"type:date" json:"due_date"`
	PaidAt           *time.Time `json:"paid_at"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`

	// 关联
	Room   *Room   `gorm:"foreignKey:RoomID" json:"room,omitempty"`
	Tenant *Tenant `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
}

func (Bill) TableName() string {
	return "bills"
}

// 账单类型
const (
	BillTypeRent     int8 = 1 // 租金
	BillTypeWater    int8 = 2 // 水费
	BillTypeElectric int8 = 3 // 电费
	BillTypeGas      int8 = 4 // 气费
	BillTypeCombined int8 = 5 // 综合账单
)

// 账单状态
const (
	BillStatusPending int8 = 1 // 待支付
	BillStatusPaid    int8 = 2 // 已支付
	BillStatusOverdue int8 = 3 // 已逾期
)

// Payment 收款记录表
type Payment struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	BillID        uint      `gorm:"index;not null" json:"bill_id"`
	Amount        float64   `gorm:"type:decimal(10,2)" json:"amount"`
	PaymentMethod int8      `gorm:"comment:1现金 2微信 3支付宝 4银行转账" json:"payment_method"`
	PaidAt        time.Time `json:"paid_at"`
	Note          string    `gorm:"size:255" json:"note"`
	CreatedAt     time.Time `json:"created_at"`

	// 关联
	Bill *Bill `gorm:"foreignKey:BillID" json:"bill,omitempty"`
}

func (Payment) TableName() string {
	return "payments"
}

// 支付方式
const (
	PaymentMethodCash     int8 = 1 // 现金
	PaymentMethodWechat   int8 = 2 // 微信
	PaymentMethodAlipay   int8 = 3 // 支付宝
	PaymentMethodTransfer int8 = 4 // 银行转账
)
