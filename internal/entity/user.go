package entity

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	FirstName string `gorm:""`
	LastName  string `gorm:""`
	Email     string `gorm:""`
	Password  string
	TimeZone  string `gorm:"default:UTC+3:30"`
	ApiKey    string `gorm:"default:''"`
	IsActive  bool   `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Jobs      []Job   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	JobGroups []Group `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (u *User) GenerateNewPassword(password string) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	u.Password = string(hashedPassword)
}
