package usecase

import (
	"cron_job/internal/entity"
	"cron_job/internal/repository"
	"cron_job/pkg/exception"
	"github.com/gin-gonic/gin"
)

type GroupUseCase struct {
	repo       *repository.GroupRepository
	jobUseCase *JobUseCase
	c          *gin.Context
}

func NewGroupUseCase() *GroupUseCase {
	groupRepo := repository.NewGroupRepository()
	return &GroupUseCase{repo: groupRepo, jobUseCase: NewJobUseCase()}
}

func (u *GroupUseCase) CreateGroup(jobGroup *entity.Group) (entity.Group, error) {
	return u.repo.Create(jobGroup)
}

func (u *GroupUseCase) CreateDefaultGroup(userId uint) {
	groups := u.repo.GetAllByUserId(userId)
	if len(groups) == 0 {
		_, _ = u.repo.Create(&entity.Group{
			UserId:      userId,
			Name:        "default",
			TagName:     "def",
			Description: "",
			DefGrp:      true,
		})
	}
}

func (u *GroupUseCase) GetByIdAndUserId(id uint, userId uint) (entity.Group, *exception.AppError) {
	grp, err := u.repo.GetByIdAndUserId(id, userId)

	if err != nil || grp.ID == 0 {
		return entity.Group{}, exception.NewNotFound("group not found", "")
	}
	return grp, nil
}
func (u *GroupUseCase) GetByIdAndUserIdWithFullRelation(id uint, userId uint) (entity.Group, *exception.AppError) {
	grp, err := u.repo.GetByIdAndUserIdWithFullRelation(id, userId)
	if err != nil {
		return entity.Group{}, exception.NewNotFound("group not found", "")
	}
	return grp, nil
}

func (u *GroupUseCase) GetDefault(userId uint) (entity.Group, error) {
	return u.repo.FindOneByUserIdAndDefault(userId)
}

func (u *GroupUseCase) GetAllGroupsByUserId(userId uint) []entity.Group {
	return u.repo.GetAllByUserId(userId)
}

func (u *GroupUseCase) Update(id uint, group entity.Group) (entity.Group, *exception.AppError) {
	g, err := u.repo.GetByIdAndUserId(id, group.UserId)
	if err != nil || g.ID == 0 {
		return entity.Group{}, exception.NewNotFound("group not found!", "")
	}
	grp, err := u.repo.Update(id, &group)
	if err != nil {
		return entity.Group{}, exception.NewInternal(err.Error())
	}
	return grp, nil
}

func (u *GroupUseCase) Delete(id uint, userId uint) *exception.AppError {
	group, err := u.repo.GetByIdAndUserId(id, userId)
	if err != nil || group.ID == 0 {
		return exception.NewNotFound("group not found!!", "")
	}

	if group.DefGrp {
		return exception.NewBusinessError("not allowed to delete this group!", "")
	}
	defGrp, _ := u.GetDefault(userId)
	if len(group.Jobs) > 0 {
		for _, job := range group.Jobs {
			u.jobUseCase.UpdateGroupId(job.ID, defGrp.ID)
		}
	}
	u.repo.Delete(id)
	return nil
}
