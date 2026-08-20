package entity

import (
	"gorm.io/gorm"
	"time"
)

type Job struct {
	ID              uint `gorm:"primaryKey"`
	GroupId         uint
	UserId          uint
	TimeZone        string
	Name            string
	Description     string
	Schedule        string
	ExecutionPerDay uint           `gorm:"default:0"`
	TotalSuccess    uint           `gorm:"default:0"`
	TotalFail       uint           `gorm:"default:0"`
	IsActive        *bool          `gorm:"default:false"`
	DeletedAt       gorm.DeletedAt `gorm:"index"` // Soft delete field
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Group           Group
	User            User
	RequestHttp     RequestHttp
	Notifications   []Notification ` json:"notifications" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
