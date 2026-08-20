package repository

import (
	"cron_job/internal/config"
	"cron_job/internal/entity"
	"gorm.io/gorm"
)

type ReqHttpRepository struct {
	db *gorm.DB
}

func NewReqHttpRepository() *ReqHttpRepository {
	return &ReqHttpRepository{db: config.DB}
}

func (r *ReqHttpRepository) Create(reqHttp *entity.RequestHttp) (entity.RequestHttp, error) {
	err := r.db.Create(&reqHttp).Error
	return *reqHttp, err
}

func (r *ReqHttpRepository) GetByJobId(jobId uint) (entity.RequestHttp, error) {
	var reqHttp entity.RequestHttp
	err := r.db.Where("job_id = ?", jobId).First(&reqHttp).Error
	return reqHttp, err
}

func (r *ReqHttpRepository) Update(jobId uint, reqHttpData *entity.RequestHttp) (*entity.RequestHttp, error) {
	reqHttp, _ := r.GetByJobId(jobId)
	if reqHttpData.Url != "" {
		reqHttp.Url = reqHttpData.Url
	}
	if reqHttpData.Method != "" {
		reqHttp.Method = reqHttpData.Method
	}
	if len(reqHttpData.Headers) > 0 {
		reqHttp.Headers = reqHttpData.Headers
	}
	if reqHttpData.Body.String() != "" {
		reqHttp.Body = reqHttpData.Body
	}
	if reqHttpData.TimeOut != 0 {
		reqHttp.TimeOut = reqHttpData.TimeOut
	}

	err := r.db.Save(&reqHttp).Error
	if err != nil {
		return nil, err
	}

	updatedReqHttp, _ := r.GetByJobId(jobId)
	return &updatedReqHttp, err
}
