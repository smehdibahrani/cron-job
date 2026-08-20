package usecase

import (
	"cron_job/internal/entity"
	"cron_job/internal/repository"
	"cron_job/pkg"
	"cron_job/pkg/exception"
	"fmt"
)

type JobUseCase struct {
	repo                *repository.JobRepository
	reqHttpUseCase      *ReqHttpUseCase
	notificationUseCase *NotificationUseCase
}

func NewJobUseCase() *JobUseCase {
	return &JobUseCase{repo: repository.NewJobRepository(), reqHttpUseCase: NewReqHttpUseCase(), notificationUseCase: NewNotificationUseCase()}
}

func (u *JobUseCase) CreateJob(job *entity.Job, reqHttp *entity.RequestHttp, notification *entity.Notification) (entity.Job, *exception.AppError) {
	if !pkg.IsValidCron(job.Schedule) {
		return entity.Job{}, exception.NewBadRequest("cron pattern invalid", "")
	}

	job.ExecutionPerDay = pkg.CalculateEpd(job.Schedule)

	job.ExecutionPerDay = pkg.CalculateEpd(job.Schedule)
	jobCreated, err := u.repo.Create(job)
	if err != nil {
		return entity.Job{}, exception.NewDatabaseError(fmt.Sprintf("failed to create job: %v", err))
	}

	reqHttp.JobId = jobCreated.ID
	createdReqHttp, err := u.reqHttpUseCase.create(reqHttp)
	if err != nil {
		// Rollback: delete the created job if request creation fails
		u.repo.Delete(jobCreated.ID)
		return entity.Job{}, exception.NewDatabaseError(fmt.Sprintf("failed to create request HTTP: %v", err))
	}
	jobCreated.RequestHttp = createdReqHttp

	if notification != nil {
		notification.JobId = jobCreated.ID
		createdNotification, err := u.notificationUseCase.create(notification)
		if err != nil {
			// Rollback: delete the created job and request if notification creation fails
			u.repo.Delete(jobCreated.ID)
			return entity.Job{}, exception.NewDatabaseError(fmt.Sprintf("failed to create notification: %v", err))
		}
		jobCreated.Notifications = []entity.Notification{createdNotification}
	}

	return jobCreated, nil
}

func (u *JobUseCase) Update(id uint, jobData *entity.Job, reqHttp *entity.RequestHttp, notification *entity.Notification) (*entity.Job, bool, *exception.AppError) {
	job, err := u.repo.GetByIdAndUserId(id, jobData.UserId)
	if err != nil {
		return nil, false, exception.NewNotFound(err.Error(), "job")
	}
	canActiveJob := false
	if *jobData.IsActive == true && *job.IsActive == false {
		canActiveJob = true
	}

	if jobData.Schedule != "" {
		if !pkg.IsValidCron(jobData.Schedule) {
			return nil, false, exception.NewBadRequest("cron pattern incorrect", "")
		}

		job.ExecutionPerDay = pkg.CalculateEpd(jobData.Schedule)

		job.Schedule = jobData.Schedule

	}
	if jobData.TimeZone != "" {
		job.TimeZone = jobData.TimeZone
	}
	if jobData.Name != "" {
		job.Name = jobData.Name
	}
	if jobData.Description != "" {
		job.Description = jobData.Description
	}
	if jobData.IsActive != nil {
		job.IsActive = jobData.IsActive
	}
	job.GroupId = jobData.GroupId

	updatedJob, err := u.repo.Update(id, &job)
	if err != nil {
		return nil, false, exception.NewInternal(err.Error())
	}

	if reqHttp != nil {
		updateReqHttp, err := u.reqHttpUseCase.Update(id, reqHttp)
		if err != nil {
			return nil, false, exception.NewInternal(err.Error())
		}
		updatedJob.RequestHttp = *updateReqHttp
	}

	if notification != nil {
		jobNotification, _ := u.notificationUseCase.repo.GetByJobId(id)
		if jobNotification != nil {
			updateNotification, _ := u.notificationUseCase.Update(id, notification)
			updatedJob.Notifications = []entity.Notification{updateNotification}
		} else {
			notification.JobId = id
			updateNotification, _ := u.notificationUseCase.create(notification)
			updatedJob.Notifications = []entity.Notification{updateNotification}
		}
	}

	return updatedJob, canActiveJob, nil
}

func (u *JobUseCase) GetAllJobsByUserId(userId uint) ([]entity.Job, error) {
	return u.repo.GetAllJobsByUserId(userId)
}

func (u *JobUseCase) GetCountOfEPDByUserId(userId uint) uint {
	count, _ := u.repo.CountOFEpdJobsByUserId(userId)
	return count
}

func (u *JobUseCase) GetByIdAndUserId(id uint, userId uint) (entity.Job, *exception.AppError) {
	job, err := u.repo.GetByIdAndUserId(id, userId)
	fmt.Println(job, err)
	if err != nil {
		return entity.Job{}, exception.NewNotFound("no job with this id", "job")
	}
	return job, nil
}

func (u *JobUseCase) GetById(id uint) (entity.Job, error) {
	return u.repo.GetById(id)
}

func (u *JobUseCase) UpdateGroupId(id uint, grpId uint) {
	u.repo.UpdateGroupId(id, grpId)
}

func (u *JobUseCase) UpdateTotal(id uint, totalSuccess uint, totalFail uint) {
	u.repo.UpdateTotal(id, totalSuccess, totalFail)
}

func (u *JobUseCase) DeleteJob(id uint, userId uint) *exception.AppError {
	_, err := u.GetByIdAndUserId(id, userId)
	if err != nil {
		return exception.NewNotFound("no job with this id", "job")
	}
	u.repo.Delete(id)
	return nil
}
