package entity

import "time"

type Log struct {
	ID            uint `gorm:"primaryKey"`
	RequestHttpId uint
	JobId         uint
	UserId        uint
	Url           string
	Method        string
	Res           string
	ResStatus     uint
	ResTime       uint
	ErrorMessage  string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	RequestHttp   RequestHttp
	Job           Job
	User          User
}
