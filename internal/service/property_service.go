package service

import (
	"errors"

	"rental-management/internal/model"
	"rental-management/internal/repository"
)

type PropertyService struct {
	propertyRepo *repository.PropertyRepository
	roomRepo     *repository.RoomRepository
}

func NewPropertyService(propertyRepo *repository.PropertyRepository, roomRepo *repository.RoomRepository) *PropertyService {
	return &PropertyService{
		propertyRepo: propertyRepo,
		roomRepo:     roomRepo,
	}
}

// CreatePropertyRequest 创建房源请求
type CreatePropertyRequest struct {
	Name         string `json:"name" binding:"required"`
	Address      string `json:"address"`
	PropertyType int8   `json:"property_type" binding:"required"`
}

// UpdatePropertyRequest 更新房源请求
type UpdatePropertyRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Status  *int8  `json:"status"`
}

// CreateProperty 创建房源
func (s *PropertyService) CreateProperty(userID uint, req *CreatePropertyRequest) (*model.Property, error) {
	property := &model.Property{
		UserID:       userID,
		Name:         req.Name,
		Address:      req.Address,
		PropertyType: req.PropertyType,
		TotalRooms:   0,
		Status:       1,
	}

	if err := s.propertyRepo.Create(property); err != nil {
		return nil, err
	}

	return property, nil
}

// GetProperty 获取房源详情
func (s *PropertyService) GetProperty(id, userID uint) (*model.Property, error) {
	return s.propertyRepo.FindByIDAndUserID(id, userID)
}

// ListProperties 获取房源列表
func (s *PropertyService) ListProperties(userID uint, page, pageSize int) ([]model.Property, int64, error) {
	return s.propertyRepo.List(userID, page, pageSize)
}

// UpdateProperty 更新房源
func (s *PropertyService) UpdateProperty(id, userID uint, req *UpdatePropertyRequest) (*model.Property, error) {
	property, err := s.propertyRepo.FindByIDAndUserID(id, userID)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		property.Name = req.Name
	}
	if req.Address != "" {
		property.Address = req.Address
	}
	if req.Status != nil {
		property.Status = *req.Status
	}

	if err := s.propertyRepo.Update(property); err != nil {
		return nil, err
	}

	return property, nil
}

// DeleteProperty 删除房源
func (s *PropertyService) DeleteProperty(id, userID uint) error {
	property, err := s.propertyRepo.FindByIDAndUserID(id, userID)
	if err != nil {
		return err
	}

	// 检查是否有房间
	rooms, err := s.roomRepo.FindByPropertyID(property.ID)
	if err != nil {
		return err
	}
	if len(rooms) > 0 {
		return errors.New("该房源下存在房间，无法删除")
	}

	return s.propertyRepo.DeleteByIDAndUserID(id, userID)
}

// CreateRoomRequest 创建房间请求
type CreateRoomRequest struct {
	PropertyID  uint    `json:"property_id" binding:"required"`
	RoomNumber  string  `json:"room_number" binding:"required"`
	Floor       int     `json:"floor"`
	Area        float64 `json:"area"`
	RentType    int8    `json:"rent_type"`
	RentAmount  float64 `json:"rent_amount"`
	MonthlyRent float64 `json:"monthly_rent"`
	Facilities  string  `json:"facilities"`
	Remark      string  `json:"remark"`
}

// UpdateRoomRequest 更新房间请求
type UpdateRoomRequest struct {
	RoomNumber  string  `json:"room_number"`
	Floor       *int    `json:"floor"`
	Area        float64 `json:"area"`
	RentType    *int8   `json:"rent_type"`
	RentAmount  float64 `json:"rent_amount"`
	MonthlyRent float64 `json:"monthly_rent"`
	Facilities  string  `json:"facilities"`
	Remark      string  `json:"remark"`
}

// CreateRoom 创建房间
func (s *PropertyService) CreateRoom(userID uint, req *CreateRoomRequest) (*model.Room, error) {
	// 验证房源归属
	property, err := s.propertyRepo.FindByIDAndUserID(req.PropertyID, userID)
	if err != nil {
		return nil, errors.New("房源不存在或无权限")
	}

	// Facilities 必须是有效的 JSON，空时使用空数组
	facilities := req.Facilities
	if facilities == "" {
		facilities = "[]"
	}

	// 默认月租类型
	rentType := req.RentType
	if rentType == 0 {
		rentType = model.RentTypeMonthly
	}

	room := &model.Room{
		PropertyID:  req.PropertyID,
		RoomNumber:  req.RoomNumber,
		Floor:       req.Floor,
		Area:        req.Area,
		RentType:    rentType,
		RentAmount:  req.RentAmount,
		MonthlyRent: req.MonthlyRent,
		Status:      model.RoomStatusVacant,
		Facilities:  facilities,
		Remark:      req.Remark,
	}

	if err := s.roomRepo.Create(room); err != nil {
		return nil, err
	}

	// 更新房源房间数
	count, _ := s.roomRepo.CountByPropertyID(property.ID)
	_ = s.propertyRepo.UpdateRoomCount(property.ID, int(count))

	return room, nil
}

// GetRoom 获取房间详情
func (s *PropertyService) GetRoom(id, userID uint) (*model.Room, error) {
	return s.roomRepo.FindByIDAndUserID(id, userID)
}

// ListRooms 获取房间列表
func (s *PropertyService) ListRooms(userID uint, propertyID uint, status *int8, page, pageSize int) ([]repository.RoomWithTenant, int64, error) {
	return s.roomRepo.ListWithTenant(userID, propertyID, status, page, pageSize)
}

// UpdateRoom 更新房间
func (s *PropertyService) UpdateRoom(id, userID uint, req *UpdateRoomRequest) (*model.Room, error) {
	room, err := s.roomRepo.FindByIDAndUserID(id, userID)
	if err != nil {
		return nil, err
	}

	if req.RoomNumber != "" {
		room.RoomNumber = req.RoomNumber
	}
	if req.Floor != nil {
		room.Floor = *req.Floor
	}
	room.Area = req.Area
	if req.RentType != nil {
		room.RentType = *req.RentType
	}
	room.RentAmount = req.RentAmount
	room.MonthlyRent = req.MonthlyRent
	if req.Facilities != "" {
		room.Facilities = req.Facilities
	} else {
		room.Facilities = "[]"
	}
	if req.Remark != "" {
		room.Remark = req.Remark
	}

	if err := s.roomRepo.Update(room); err != nil {
		return nil, err
	}

	return room, nil
}

// DeleteRoom 删除房间
func (s *PropertyService) DeleteRoom(id, userID uint) error {
	room, err := s.roomRepo.FindByIDAndUserID(id, userID)
	if err != nil {
		return err
	}

	// 检查房间状态
	if room.Status == model.RoomStatusRented {
		return errors.New("房间已出租，无法删除")
	}

	if err := s.roomRepo.DeleteByIDAndUserID(id, userID); err != nil {
		return err
	}

	// 更新房源房间数
	count, _ := s.roomRepo.CountByPropertyID(room.PropertyID)
	_ = s.propertyRepo.UpdateRoomCount(room.PropertyID, int(count))

	return nil
}

// UpdateRoomStatus 更新房间状态
func (s *PropertyService) UpdateRoomStatus(id, userID uint, status int8) error {
	// 先验证权限
	_, err := s.roomRepo.FindByIDAndUserID(id, userID)
	if err != nil {
		return err
	}

	return s.roomRepo.UpdateStatusByUserID(id, userID, status)
}
