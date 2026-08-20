package usecase

import (
	"cron_job/internal/entity"
	"cron_job/internal/repository"
	"github.com/gin-gonic/gin"
)

type NotificationUseCase struct {
	repo *repository.NotificationRepository
	c    *gin.Context
}

func NewNotificationUseCase() *NotificationUseCase {
	repo := repository.NewNotificationRepository()
	return &NotificationUseCase{repo: repo}
}

func (u *NotificationUseCase) create(notif *entity.Notification) (entity.Notification, error) {
	return u.repo.Create(notif)
}

func (u *NotificationUseCase) Update(id uint, notification *entity.Notification) (entity.Notification, error) {
	return u.repo.Update(id, notification)
}
