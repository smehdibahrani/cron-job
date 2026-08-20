package usecase

import (
	"cron_job/internal/entity"
	"cron_job/internal/repository"
	"github.com/gin-gonic/gin"
)

type ReqHttpUseCase struct {
	repo *repository.ReqHttpRepository
	c    *gin.Context
}

func NewReqHttpUseCase() *ReqHttpUseCase {
	repo := repository.NewReqHttpRepository()
	return &ReqHttpUseCase{repo: repo}
}

func (u *ReqHttpUseCase) create(reqHttp *entity.RequestHttp) (entity.RequestHttp, error) {
	return u.repo.Create(reqHttp)
}

func (u *ReqHttpUseCase) Update(jobId uint, reqHttp *entity.RequestHttp) (*entity.RequestHttp, error) {
	_, err := u.repo.GetByJobId(jobId)
	if err != nil {
		return nil, err
	}
	return u.repo.Update(jobId, reqHttp)
}
