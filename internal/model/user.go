package model

import (
	"time"

	"gorm.io/gorm"
)

// User 房东用户表
type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Phone        string         `gorm:"uniqueIndex;size:20;not null" json:"phone"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	Name         string         `gorm:"size:50" json:"name"`
	Avatar       string         `gorm:"size:255" json:"avatar"`
	Status       int8           `gorm:"default:1;comment:1正常 2禁用" json:"status"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "users"
}

// FeeRate 费率配置表
type FeeRate struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	UserID           uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	WaterRate        float64   `gorm:"type:decimal(10,4);default:0" json:"water_rate"`        // 水费单价(元/吨)
	ElectricityRate  float64   `gorm:"type:decimal(10,4);default:0" json:"electricity_rate"`  // 电费单价(元/度)
	GasRate          float64   `gorm:"type:decimal(10,4);default:0" json:"gas_rate"`          // 气费单价(元/立方)
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (FeeRate) TableName() string {
	return "fee_rates"
}
