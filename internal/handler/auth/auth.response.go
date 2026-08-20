package auth

import (
	"cron_job/internal/entity"
	"cron_job/pkg/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleResponse(user entity.User, c *gin.Context) {
	accessToken, _ := jwt.GenerateToken(user.ID, false)
	refreshToken, _ := jwt.GenerateToken(user.ID, true)

	c.JSON(http.StatusOK, gin.H{
		"firstName":    user.FirstName,
		"lastName":     user.LastName,
		"email":        user.Email,
		"createdAt":    user.CreatedAt,
		"accessToken":  accessToken,
		"refreshToken": refreshToken})
}
