package repository

import (
	"cron_job/internal/config"
	"cron_job/internal/entity"
	"gorm.io/gorm"
)

type JobRepository struct {
	db *gorm.DB
}

func NewJobRepository() *JobRepository {
	return &JobRepository{db: config.DB}
}

func (r *JobRepository) Create(job *entity.Job) (entity.Job, error) {
	err := r.db.Create(&job).Error
	return *job, err
}

func (r *JobRepository) GetAllJobsByUserId(userId uint) ([]entity.Job, error) {
	var jobs []entity.Job
	err := r.db.Preload("Group").Preload("Notifications").Preload("RequestHttp").Where("user_id = ?", userId).Find(&jobs).Error
	return jobs, err
}

func (r *JobRepository) CountOFEpdJobsByUserId(userId uint) (uint, error) {
	var count uint
	err := r.db.Model(&entity.Job{}).
		Where("user_id = ?", userId).
		Select("COUNT(execution_per_day)").
		Scan(&count).Error
	return count, err
}

func (r *JobRepository) GetByIdAndUserId(id uint, userId uint) (entity.Job, error) {
	var job entity.Job
	err := r.db.Preload("Group").Preload("Notifications").Preload("RequestHttp").Where("user_id = ? and id = ?", userId, id).First(&job).Error
	return job, err
}

func (r *JobRepository) GetById(id uint) (entity.Job, error) {
	var job entity.Job
	err := r.db.
		Preload("Group").
		Preload("RequestHttp").
		Preload("Notifications").
		Where("id = ?", id).First(&job).Error
	return job, err
}

func (r *JobRepository) Update(id uint, job *entity.Job) (*entity.Job, error) {
	updateData := map[string]interface{}{"execution_per_day": job.ExecutionPerDay, "schedule": job.Schedule, "is_active": job.IsActive}
	if job.GroupId > 0 {
		updateData["group_id"] = job.GroupId
	}
	r.db.Model(&entity.Job{}).Where("id = ?", id).Updates(updateData)
	jobUpdated, err := r.GetById(id)
	return &jobUpdated, err
}

func (r *JobRepository) UpdateGroupId(id uint, grpId uint) {
	updateData := map[string]interface{}{"group_id": grpId}
	r.db.Model(&entity.Job{}).Where("id = ?", id).Updates(updateData)
}
func (r *JobRepository) UpdateTotal(id uint, totalSuccess uint, totalFail uint) {
	updateData := map[string]interface{}{"total_success": totalSuccess, "total_fail": totalFail}
	r.db.Model(&entity.Job{}).Where("id = ?", id).Updates(updateData)
}

func (r *JobRepository) Delete(id uint) {
	r.db.Delete(&entity.Job{}, id)
}
