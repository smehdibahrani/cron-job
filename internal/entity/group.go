package entity

import (
	"time"
)

type Group struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserId      uint      `json:"userId"`
	Name        string    `json:"name"`
	TagName     string    `json:"tagName"`
	Description string    `json:"description"`
	JobCount    uint      `json:"jobCount"`
	DefGrp      bool      `gorm:"default:false" json:"defGrp"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	User        User      `json:"user"`
	Jobs        []Job     `json:"jobs" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
