package usecase

import (
	"cron_job/internal/entity"
	"cron_job/internal/repository"
	"cron_job/pkg/exception"
	"fmt"
)

type UserUseCase struct {
	repo *repository.UserRepository
}

func NewUserUseCase() *UserUseCase {
	return &UserUseCase{repo: repository.NewUserRepository()}
}

func (u *UserUseCase) Register(user *entity.User) *exception.AppError {
	// Check if user already exists
	if existingUser := u.GetByEmail(user.Email); existingUser != nil {
		return exception.NewConflict("email already exists", fmt.Sprintf("email: %s", user.Email))
	}

	err := u.repo.Register(user)
	if err != nil {
		return exception.NewDatabaseError(fmt.Sprintf("failed to register user: %v", err))
	}

	return nil
}

func (u *UserUseCase) UpsertByGoogleAuth(user *entity.User) *exception.AppError {
	existingUser := u.GetByEmail(user.Email)
	if existingUser == nil {
		err := u.repo.Register(user)
		if err != nil {
			return exception.NewDatabaseError(fmt.Sprintf("failed to register user: %v", err))
		}
	} else {
		user = existingUser
	}
	return nil
}

func (u *UserUseCase) GetById(id uint) (entity.User, *exception.AppError) {
	user, err := u.repo.GetById(id)
	if err != nil {
		return entity.User{}, exception.NewNotFound("user not found", fmt.Sprintf("id: %d", id))
	}
	return user, nil
}

func (u *UserUseCase) GetByEmail(email string) *entity.User {
	userExists, _ := u.repo.FindOneByEmail(email)
	return userExists
}

func (u *UserUseCase) GetByPhoneAndIsVerified(phone string) *entity.User {
	userExists, _ := u.repo.FindOneByPhoneAndIsVerified(phone)
	return userExists
}

func (u *UserUseCase) GetByEmailAndIsVerified(email string) *entity.User {
	userExists, _ := u.repo.FindOneByEmailAndIsVerified(email)
	return userExists
}

func (u *UserUseCase) Update(id uint, user entity.User) (*entity.User, *exception.AppError) {
	_, errGet := u.GetById(id)
	if errGet != nil {
		return nil, errGet
	}

	userUpdated, err := u.repo.Update(id, &user)
	if err != nil {
		return nil, exception.NewDatabaseError(fmt.Sprintf("failed to update user: %v", err))
	}

	return &userUpdated, nil
}
