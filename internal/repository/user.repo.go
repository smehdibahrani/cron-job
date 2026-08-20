package repository

import (
	"cron_job/internal/config"
	"cron_job/internal/entity"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{db: config.DB}
}

func (r *UserRepository) Register(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindOneByEmailAndIsActiveAndVerified(email string) (entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ? and is_email_verified = ? and is_active = ?", email, true, true).First(&user).Error
	return user, err
}

func (r *UserRepository) GetById(userId uint) (entity.User, error) {
	var user entity.User
	err := r.db.Where("id = ?", userId).First(&user).Error
	return user, err
}

func (r *UserRepository) FindOneByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (r *UserRepository) FindOneByEmailAndIsVerified(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ? and is_email_verified = ? ", email, true).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (r *UserRepository) FindOneByPhoneAndIsVerified(phone string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("phone = ? and is_phone_verified = true", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (r *UserRepository) Update(id uint, user *entity.User) (entity.User, error) {
	err := r.db.Where("id = ?", id).Updates(user).Error
	updatedUser, err := r.GetById(id)
	return updatedUser, err
}
