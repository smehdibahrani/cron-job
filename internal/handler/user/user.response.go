package user

import (
	"cron_job/internal/entity"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleResponse(user entity.User, c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"email":     user.Email,
		"createdAt": user.CreatedAt})
}
