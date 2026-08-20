package usecase

import (
	"cron_job/internal/entity"
	"cron_job/internal/repository"
	"cron_job/pkg/exception"
	"fmt"
	"github.com/gin-gonic/gin"
)

type LogUseCase struct {
	repo *repository.LogRepository
	ctx  *gin.Context
}

func NewLogUseCase() *LogUseCase {
	return &LogUseCase{repo: repository.NewLogRepository()}
}

func (u *LogUseCase) Create(log *entity.Log) (entity.Log, *exception.AppError) {
	createdLog, err := u.repo.Create(log)
	if err != nil {
		return entity.Log{}, exception.NewDatabaseError(fmt.Sprintf("failed to create log: %v", err))
	}

	return createdLog, nil
}

func (u *LogUseCase) GetByIdAndUserId(userId uint) ([]entity.Log, *exception.AppError) {
	if userId == 0 {
		return nil, exception.NewValidationError("userId", "userId is required")
	}

	logs, err := u.repo.GetAllByUserId(userId)
	if err != nil {
		return nil, exception.NewDatabaseError(fmt.Sprintf("failed to get logs for user %d: %v", userId, err))
	}

	return logs, nil
}

func (u *LogUseCase) GetByJobIdAndUserId(jobId uint, userId uint) ([]entity.Log, *exception.AppError) {
	if jobId == 0 {
		return nil, exception.NewValidationError("jobId", "job id required")
	}

	if userId == 0 {
		return nil, exception.NewValidationError("userId", "user id required")
	}

	logs, err := u.repo.GetAllByJobIdAndUserId(jobId, userId)
	if err != nil {
		return nil, exception.NewDatabaseError(fmt.Sprintf("failed to get logs for job %d and user %d: %v", jobId, userId, err))
	}

	return logs, nil
}

func (u *LogUseCase) Clean() *exception.AppError {
	u.repo.Clean()
	// Note: Clean() method in repository doesn't return error, but we could add logging here
	return nil
}
