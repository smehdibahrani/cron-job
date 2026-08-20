package job

import (
	"cron_job/internal/entity"
	"cron_job/internal/scheduler"
	"cron_job/internal/usecase"
	"cron_job/pkg"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	jobUseCase *usecase.JobUseCase
	logUseCase *usecase.LogUseCase
	grpUseCase *usecase.GroupUseCase
}

func NewJobHandler() *Handler {
	return &Handler{
		jobUseCase: usecase.NewJobUseCase(),
		logUseCase: usecase.NewLogUseCase(),
		grpUseCase: usecase.NewGroupUseCase(),
	}
}

func (h *Handler) Create(c *gin.Context) {
	var dto CreateJobDTO
	if err := pkg.MapToDto(c, &dto); err != nil {
		pkg.HandleError(c, err)
		return
	}
	var job entity.Job
	var notification *entity.Notification
	var reqHttp entity.RequestHttp

	job.UserId = pkg.GetUserIdFromReq(c)
	dto.JobData.MapDtoToJob(&job)
	if dto.Notification != nil {
		notification = &entity.Notification{}
		dto.Notification.MapDtoToNotification(notification)
	}
	dto.HttpRequest.MapDtoToReqHttp(&reqHttp)

	if job.GroupId == 0 {
		grpDefault, _ := h.grpUseCase.GetDefault(job.UserId)
		job.GroupId = grpDefault.ID
	} else {
		_, err := h.grpUseCase.GetByIdAndUserId(job.GroupId, job.UserId)
		if err != nil {
			pkg.HandleError(c, err)
			return
		}
	}

	job.UserId = pkg.GetUserIdFromReq(c)

	jobCreated, err := h.jobUseCase.CreateJob(&job, &reqHttp, notification)
	if err != nil {
		pkg.HandleError(c, err)
		return
	}
	if *job.IsActive == true {
		go scheduler.AddToRedisQueue(job)
	}

	c.JSON(http.StatusOK, generateJobResponse(jobCreated))
}

func (h *Handler) Update(c *gin.Context) {
	var job entity.Job
	var notification *entity.Notification
	var reqHttp entity.RequestHttp
	var dto UpdateJobDTO

	job.UserId = pkg.GetUserIdFromReq(c)

	id := pkg.GetIntParam(c, "id")
	if err := pkg.MapToDto(c, &dto); err != nil {
		pkg.HandleError(c, err)
		return
	}

	dto.JobData.MapDtoToJob(&job)
	if dto.Notification != nil {
		notification = &entity.Notification{}
		dto.Notification.MapDtoToNotification(notification)
	}
	dto.HttpRequest.MapDtoToReqHttp(&reqHttp)

	jobUpdated, canActiveJob, err := h.jobUseCase.Update(id, &job, &reqHttp, notification)
	if err != nil {
		pkg.HandleError(c, err)
		return
	}
	if *jobUpdated.IsActive == false {
		scheduler.RemoveFromRedisQueue(jobUpdated.ID)
	} else if canActiveJob == true || dto.JobData.Schedule != "" {
		scheduler.RemoveFromRedisQueue(jobUpdated.ID)
		go scheduler.AddToRedisQueue(*jobUpdated)
	}

	c.JSON(http.StatusOK, generateJobResponse(*jobUpdated))
}

func (h *Handler) GetById(c *gin.Context) {
	userId := pkg.GetUserIdFromReq(c)
	id := pkg.GetIntParam(c, "id")
	job, err := h.jobUseCase.GetByIdAndUserId(id, userId)
	if err != nil {
		pkg.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, generateJobResponse(job))
}

func (h *Handler) GetLogsByJobId(c *gin.Context) {
	userId := pkg.GetUserIdFromReq(c)
	id := pkg.GetIntParam(c, "id")
	logs, err := h.logUseCase.GetByJobIdAndUserId(id, userId)
	if err != nil {
		pkg.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, generateLogResponses(logs))
}

func (h *Handler) GetAll(c *gin.Context) {
	userId := pkg.GetUserIdFromReq(c)
	jobs, _ := h.jobUseCase.GetAllJobsByUserId(userId)
	c.JSON(http.StatusOK, GenerateJobResponses(jobs))
}

func (h *Handler) Delete(c *gin.Context) {
	userId := pkg.GetUserIdFromReq(c)
	id := pkg.GetIntParam(c, "id")
	if err := h.jobUseCase.DeleteJob(id, userId); err != nil {
		pkg.HandleError(c, err)
		return
	}
	scheduler.RemoveFromRedisQueue(id)
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
