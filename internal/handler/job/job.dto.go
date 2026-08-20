package job

import (
	"cron_job/internal/entity"
	"gorm.io/datatypes"
)

type CreateJobDTO struct {
	JobData      CrtJobDTO        `json:"jobData" binding:"required"`
	HttpRequest  HttpRequestDto   `json:"httpRequest" binding:"required" `
	Notification *NotificationDto `json:"notification"`
}

type UpdateJobDTO struct {
	JobData      UpdJobDTO        `json:"jobData"`
	HttpRequest  HttpRequestDto   `json:"httpRequest"`
	Notification *NotificationDto `json:"notification"`
}

type UpdJobDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Schedule    string `json:"schedule"`
	IsActive    *bool  `json:"isActive"`
	GroupId     uint   `json:"groupId"`
}
type CrtJobDTO struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Schedule    string `json:"schedule" binding:"required"`
	IsActive    *bool  `json:"isActive"`
	GroupId     uint   `json:"groupId"`
}
type HttpRequestDto struct {
	Url       string         `json:"url" binding:"required"`
	Method    entity.Method  `json:"method" binding:"required"`
	BasicAuth string         `json:"basicAuth"`
	Headers   datatypes.JSON `json:"headers"`
	Body      datatypes.JSON `json:"body"`
	TimeOut   uint           `json:"timeOut"`
}
type NotificationDto struct {
	Type        entity.NotificationType   `json:"type"`
	Action      entity.NotificationAction `json:"action"`
	Url         *string                   `json:"url"`
	Method      *entity.Method            `json:"method" `
	Sensitivity *uint                     `json:"sensitivity"`
	Data        *string                   `json:"data"`
}

func (dto *CrtJobDTO) MapDtoToJob(job *entity.Job) {
	job.GroupId = dto.GroupId
	job.Name = dto.Name
	job.Description = dto.Description
	job.Schedule = dto.Schedule
	if dto.IsActive != nil {
		job.IsActive = dto.IsActive
	}
}

func (dto *UpdJobDTO) MapDtoToJob(job *entity.Job) {
	job.GroupId = dto.GroupId
	job.Name = dto.Name
	job.Description = dto.Description
	job.Schedule = dto.Schedule
	if dto.IsActive != nil {
		job.IsActive = dto.IsActive
	}
}

func (dto *HttpRequestDto) MapDtoToReqHttp(reqHttp *entity.RequestHttp) {
	reqHttp.Url = dto.Url
	reqHttp.Method = dto.Method
	reqHttp.Headers = dto.Headers
	reqHttp.Body = dto.Body
	reqHttp.TimeOut = dto.TimeOut
}

func (dto *NotificationDto) MapDtoToNotification(notification *entity.Notification) {
	notification.Action = dto.Action
	notification.Type = dto.Type
	notification.Url = dto.Url
	notification.Method = dto.Method
	notification.Sensitivity = dto.Sensitivity
	notification.Data = dto.Data
}
