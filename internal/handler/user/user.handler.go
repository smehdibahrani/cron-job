package user

import (
	"cron_job/internal/usecase"
	"cron_job/pkg"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type Handler struct {
	useCase *usecase.UserUseCase
}

func NewUserHandler() *Handler {
	return &Handler{useCase: usecase.NewUserUseCase()}
}

func (h *Handler) Update(c *gin.Context) {
	var dto UpdateRequest
	appErr := pkg.MapToDto(c, &dto)
	if appErr != nil {
		pkg.HandleError(c, appErr)
		return
	}
	user := pkg.GetUserFromReq(c)

	user.FirstName = dto.FirstName
	user.LastName = dto.LastName
	userUpdated, err := h.useCase.Update(user.ID, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"exception": "Update user failed"})
		return
	}
	HandleResponse(*userUpdated, c)
}

func (h *Handler) GetInfo(c *gin.Context) {
	user := pkg.GetUserFromReq(c)
	HandleResponse(user, c)
}

func (h *Handler) ChangePassword(c *gin.Context) {
	var dto ChangePasswordRequest
	if err := pkg.MapToDto(c, &dto); err != nil {
		pkg.HandleError(c, err)
		return
	}

	user := pkg.GetUserFromReq(c)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.CurrentPassword)); err != nil {
		c.JSON(http.StatusPreconditionFailed, gin.H{"exception": "Invalid current password"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(dto.NewPassword), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	userUpdated, err := h.useCase.Update(user.ID, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"exception": "Update user failed"})
		return
	}
	HandleResponse(*userUpdated, c)

}
