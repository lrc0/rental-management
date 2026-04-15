package model

import (
	"time"

	"gorm.io/gorm"
)

// Property 房源/房产表
type Property struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	UserID       uint           `gorm:"index;not null" json:"user_id"`
	Name         string         `gorm:"size:100;not null" json:"name"`          // 房源名称
	Address      string         `gorm:"size:255" json:"address"`                // 地址
	PropertyType int8           `gorm:"not null;comment:1整栋 2单套 3商铺" json:"property_type"`
	TotalRooms   int            `gorm:"default:0" json:"total_rooms"`           // 总房间数
	Status       int8           `gorm:"default:1;comment:1正常 2已下架" json:"status"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Rooms []Room `gorm:"foreignKey:PropertyID" json:"rooms,omitempty"`
}

func (Property) TableName() string {
	return "properties"
}

// Room 房间表
type Room struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	PropertyID   uint           `gorm:"index;not null" json:"property_id"`
	RoomNumber   string         `gorm:"size:20;not null" json:"room_number"`    // 房间号
	Floor        int            `gorm:"default:0" json:"floor"`                 // 楼层
	Area         float64        `gorm:"type:decimal(10,2)" json:"area"`         // 面积(平米)
	RentType     int8           `gorm:"default:1;comment:1月租 2季租 3年租" json:"rent_type"` // 租金类型
	RentAmount   float64        `gorm:"type:decimal(10,2)" json:"rent_amount"`  // 租金金额
	MonthlyRent  float64        `gorm:"type:decimal(10,2)" json:"monthly_rent"` // 月租金(兼容旧数据)
	Status       int8           `gorm:"default:1;comment:1空置 2已租 3维修中" json:"status"`
	Facilities   string         `gorm:"type:json" json:"facilities"`            // 设施配置
	Remark       string         `gorm:"size:255" json:"remark"`                 // 备注
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Property *Property `gorm:"foreignKey:PropertyID" json:"property,omitempty"`
}

func (Room) TableName() string {
	return "rooms"
}

// 房间状态常量
const (
	RoomStatusVacant   int8 = 1 // 空置
	RoomStatusRented   int8 = 2 // 已租
	RoomStatusRepair   int8 = 3 // 维修中
)

// 租金类型常量
const (
	RentTypeMonthly int8 = 1 // 月租
	RentTypeQuarterly int8 = 2 // 季租
	RentTypeYearly int8 = 3 // 年租
)

// 房源类型常量
const (
	PropertyTypeBuilding int8 = 1 // 整栋
	PropertyTypeApartment int8 = 2 // 单套
	PropertyTypeShop      int8 = 3 // 商铺
)
