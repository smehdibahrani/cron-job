package repository

import (
	"cron_job/internal/config"
	"cron_job/internal/entity"
	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository() *NotificationRepository {
	return &NotificationRepository{db: config.DB}
}

func (r *NotificationRepository) Create(notification *entity.Notification) (entity.Notification, error) {
	err := r.db.Create(&notification).Error
	return *notification, err
}

func (r *NotificationRepository) GetByJobId(jobId uint) (*entity.Notification, error) {
	var notification entity.Notification
	err := r.db.Where("job_id = ?", jobId).First(&notification).Error
	if err != nil {
		return nil, err
	}
	return &notification, err
}

func (r *NotificationRepository) Update(jobId uint, notificationData *entity.Notification) (entity.Notification, error) {
	notification, _ := r.GetByJobId(jobId)

	if notificationData.Type != "" {
		notification.Type = notificationData.Type
	}
	if notificationData.Action != "" {
		notification.Action = notificationData.Action
	}
	if notificationData.Url != nil {
		notification.Url = notificationData.Url
	}
	if notificationData.Method != nil {
		notification.Method = notificationData.Method
	}
	if notificationData.Sensitivity != nil {
		notification.Sensitivity = notificationData.Sensitivity
	}
	if notificationData.Data != nil {
		notification.Data = notificationData.Data
	}
	err := r.db.Save(&notification).Error
	updatedNotification, _ := r.GetByJobId(jobId)
	return *updatedNotification, err
}
