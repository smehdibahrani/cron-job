package group

import (
	"cron_job/internal/entity"
	"time"
)

type GroupResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	TagName     string    `json:"tagName"`
	Description string    `json:"description"`
	DefGrp      bool      `json:"defGrp"`
	JobCount    uint      `json:"jobCount"`
	CreatedAt   time.Time `json:"createdAt"`
}

func generateResponse(jobGroup entity.Group) GroupResponse {
	return GroupResponse{
		ID:          jobGroup.ID,
		Name:        jobGroup.Name,
		TagName:     jobGroup.TagName,
		Description: jobGroup.Description,
		DefGrp:      jobGroup.DefGrp,
		JobCount:    jobGroup.JobCount,
		CreatedAt:   jobGroup.CreatedAt,
	}
}

func generateResponses(groups []entity.Group) []GroupResponse {
	res := make([]GroupResponse, 0, len(groups))
	for _, grp := range groups {
		res = append(res, generateResponse(grp))
	}
	return res
}
