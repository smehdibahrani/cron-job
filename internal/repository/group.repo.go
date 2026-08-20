package repository

import (
	"cron_job/internal/config"
	"cron_job/internal/entity"
	"fmt"
	"gorm.io/gorm"
)

type GroupRepository struct {
	db *gorm.DB
}

func NewGroupRepository() *GroupRepository {
	return &GroupRepository{db: config.DB}
}

func (r *GroupRepository) Create(group *entity.Group) (entity.Group, error) {
	err := r.db.Create(&group).Error
	return *group, err
}

func (r *GroupRepository) GetAllByUserId(userId uint) []entity.Group {
	var groups []entity.Group
	query := fmt.Sprintf(`
    SELECT g.id, g.user_id, g.name, g.tag_name, g.description, COUNT(j.id) AS job_count, g.def_grp, g.created_at, g.updated_at
    FROM groups g
    LEFT JOIN jobs j ON j.group_id = g.id
    where g.user_id = %d
    GROUP BY g.id, g.user_id, g.name, g.tag_name, g.description, g.def_grp, g.created_at, g.updated_at
`, userId)
	r.db.Raw(query).Scan(&groups)
	return groups
}

func (r *GroupRepository) GetByIdAndUserId(id uint, userId uint) (entity.Group, error) {
	var group entity.Group
	err := r.db.Table("groups g").
		Select("g.*, COUNT(j.id) AS job_count").
		Joins("LEFT JOIN jobs j ON j.group_id = g.id").
		Where("g.id = ? AND g.user_id = ?", id, userId).
		Group("g.id").
		Scan(&group).Error

	if err != nil || group.ID == 0 {
		return entity.Group{}, err
	}
	return group, err
}

func (r *GroupRepository) GetByIdAndUserIdWithFullRelation(id uint, userId uint) (entity.Group, error) {
	var group entity.Group
	err := r.db.Preload("Jobs.Group").Preload("Jobs.Notifications").Preload("Jobs.RequestHttp").Where("id = ? and user_id = ?", id, userId).First(&group).Error
	return group, err
}

func (r *GroupRepository) FindOneByUserIdAndDefault(userId uint) (entity.Group, error) {
	var group entity.Group
	err := r.db.Preload("Jobs").Where("def_grp = true and user_id = ?", userId).First(&group).Error
	return group, err
}

func (r *GroupRepository) GetById(id uint) (entity.Group, error) {
	var group entity.Group
	err := r.db.Where("id = ?", id).First(&group).Error
	return group, err
}

func (r *GroupRepository) Update(id uint, group *entity.Group) (entity.Group, error) {
	err := r.db.Where("id = ?", id).Updates(group).Error
	updatedGroup, _ := r.GetById(id)
	return updatedGroup, err
}

func (r *GroupRepository) Delete(id uint) {
	r.db.Delete(&entity.Group{}, id)
}
