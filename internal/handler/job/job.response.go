package job

import (
	"cron_job/internal/entity"
	"time"
)

type Response struct {
	ID              uint                  `json:"id"`
	Name            string                `json:"name"`
	GroupId         uint                  `json:"groupId"`
	GroupName       string                `json:"groupName"`
	Schedule        string                `json:"schedule"`
	ExecutionPerDay uint                  `json:"executionPerDay"`
	TotalSuccess    uint                  `json:"totalSuccess"`
	TotalFail       uint                  `json:"totalFail"`
	IsActive        *bool                 `json:"isActive"`
	RequestHttp     entity.RequestHttp    `json:"requestHttp"`
	Notification    *NotificationResponse `json:"notification"`
	CreatedAt       time.Time             `json:"createdAt"`
	UpdatedAt       time.Time             `json:"updatedAt"`
}

type NotificationResponse struct {
	Type        entity.NotificationType   `json:"type"`
	Action      entity.NotificationAction `json:"action"`
	Sensitivity *uint                     `json:"sensitivity"`
	CreatedAt   time.Time                 `json:"created_at"`
	UpdatedAt   time.Time                 `json:"updated_at"`
}

func generateJobResponse(job entity.Job) Response {
	// Initialize empty response
	response := Response{
		ID:              job.ID,
		Name:            job.Name,
		GroupId:         job.Group.ID,
		Schedule:        job.Schedule,
		ExecutionPerDay: job.ExecutionPerDay,
		TotalSuccess:    job.TotalSuccess,
		TotalFail:       job.TotalFail,
		IsActive:        job.IsActive,
		RequestHttp:     job.RequestHttp,
		CreatedAt:       job.CreatedAt,
		UpdatedAt:       job.UpdatedAt,
	}

	// Check if there are notifications and access the first one safely
	if len(job.Notifications) > 0 {
		notification := generateNotificationResponse(job.Notifications[0])
		response.Notification = &notification
	}
	// If no notifications, Notification will be nil and will show as null in JSON

	return response
}
func generateNotificationResponse(notification entity.Notification) NotificationResponse {
	return NotificationResponse{
		Type:        notification.Type,
		Action:      notification.Action,
		Sensitivity: notification.Sensitivity,
		CreatedAt:   notification.CreatedAt,
		UpdatedAt:   notification.UpdatedAt,
	}
}

func GenerateJobResponses(jobs []entity.Job) []Response {
	res := make([]Response, 0, len(jobs))
	for _, job := range jobs {
		res = append(res, generateJobResponse(job))
	}
	return res
}

type LogResponse struct {
	ID           uint               `json:"id"`
	JobId        uint               `json:"jobId"`
	UserId       uint               `json:"userId"`
	Url          string             `json:"url"`
	Method       string             `json:"method"`
	Res          string             `json:"res"`
	ResStatus    uint               `json:"resStatus"`
	ResTime      uint               `json:"resTime"`
	ErrorMessage string             `json:"errorMessage"`
	CreatedAt    time.Time          `json:"createdAt"`
	RequestHttp  entity.RequestHttp `json:"requestHttp"`
}

func generateLogResponse(log entity.Log) LogResponse {
	return LogResponse{
		ID:           log.ID,
		JobId:        log.JobId,
		UserId:       log.UserId,
		Url:          log.Url,
		Method:       log.Method,
		Res:          log.Res,
		ResStatus:    log.ResStatus,
		ResTime:      log.ResTime,
		ErrorMessage: log.ErrorMessage,
		CreatedAt:    log.CreatedAt,
		RequestHttp:  log.RequestHttp,
	}
}

func generateLogResponses(logs []entity.Log) []LogResponse {
	res := make([]LogResponse, 0, len(logs))
	for _, log := range logs {
		res = append(res, generateLogResponse(log))
	}
	return res
}
