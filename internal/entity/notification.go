package entity

import "time"

type NotificationType string

const (
	EMAIL   NotificationType = "email"
	WEBHOOK NotificationType = "webhook"
)

type NotificationAction string

const (
	JOB_FAILING              NotificationAction = "job_failing"
	AFTER_EACH_JOB_EXECUTION NotificationAction = "after_each_job_execution"
)

type Notification struct {
	ID          uint               `gorm:"primaryKey"`
	JobId       uint               `gorm:"unique"`
	Type        NotificationType   `gorm:"type:notificationtype"`
	Action      NotificationAction `gorm:"type:notificationaction"`
	Url         *string
	Method      *Method `gorm:"type:method"`
	Sensitivity *uint
	Data        *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Job         Job
}
