package group

import (
	"cron_job/internal/entity"
	"cron_job/internal/handler/job"
	"cron_job/internal/usecase"
	"cron_job/pkg"
	"cron_job/pkg/exception"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	useCase *usecase.GroupUseCase
}

func NewGroupHandler() *Handler {
	return &Handler{useCase: usecase.NewGroupUseCase()}
}

func (h *Handler) Create(c *gin.Context) {
	var dto CreateGroupDTO
	if err := pkg.MapToDto(c, &dto); err != nil {
		pkg.HandleError(c, err)
		return
	}
	var jobGroup entity.Group
	jobGroup.UserId = pkg.GetUserIdFromReq(c)
	jobGroup.Name = dto.Name
	jobGroup.TagName = dto.TagName
	jobGroup.Description = dto.Description

	groupCreated, err := h.useCase.CreateGroup(&jobGroup)
	if err != nil {
		pkg.HandleError(c, exception.NewInternal("error in create group: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, generateResponse(groupCreated))
}

func (h *Handler) GetAll(c *gin.Context) {
	userId := pkg.GetUserIdFromReq(c)
	groups := h.useCase.GetAllGroupsByUserId(userId)
	c.JSON(http.StatusOK, generateResponses(groups))
}

func (h *Handler) GetById(c *gin.Context) {
	userId := pkg.GetUserIdFromReq(c)
	id := pkg.GetIntParam(c, "id")

	group, err := h.useCase.GetByIdAndUserId(id, userId)
	if err != nil {
		pkg.HandleError(c, err)
		return
	}
	fmt.Println(group)
	c.JSON(http.StatusOK, generateResponse(group))
}

func (h *Handler) GetJobsByGrpId(c *gin.Context) {
	userId := pkg.GetUserIdFromReq(c)
	id := pkg.GetIntParam(c, "id")

	group, err := h.useCase.GetByIdAndUserIdWithFullRelation(id, userId)
	if err != nil {
		pkg.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, job.GenerateJobResponses(group.Jobs))
}

func (h *Handler) Update(c *gin.Context) {
	var group entity.Group
	if err := pkg.MapToDto(c, &group); err != nil {
		pkg.HandleError(c, err)
		return
	}

	id := pkg.GetIntParam(c, "id")
	group.UserId = pkg.GetUserIdFromReq(c)
	groupUpdated, err := h.useCase.Update(id, group)
	if err != nil {
		pkg.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, generateResponse(groupUpdated))
}

func (h *Handler) Delete(c *gin.Context) {
	id := pkg.GetIntParam(c, "id")
	userId := pkg.GetUserIdFromReq(c)
	if err := h.useCase.Delete(id, userId); err != nil {
		pkg.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
