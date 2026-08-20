package entity

import "time"
import "gorm.io/datatypes"

type Method string

const (
	GET    Method = "get"
	POST   Method = "post"
	PUT    Method = "put"
	DELETE Method = "delete"
	PATCH  Method = "patch"
	HEAD   Method = "head"
	OPTION Method = "options"
)

type RequestHttp struct {
	ID        uint `gorm:"primaryKey"`
	JobId     uint `gorm:"unique"`
	Url       string
	Method    Method         `gorm:"type:method"`
	Headers   datatypes.JSON `gorm:"type:json;default:'{}'"`
	Body      datatypes.JSON `gorm:"type:json;default:'{}'"`
	TimeOut   uint
	CreatedAt time.Time
	UpdatedAt time.Time
}
