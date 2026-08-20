package auth

import (
	"cron_job/internal/entity"
	"cron_job/internal/usecase"
	"cron_job/pkg"
	"cron_job/pkg/exception"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type Handler struct {
	userUseCase  *usecase.UserUseCase
	groupUseCase *usecase.GroupUseCase
}

type GoogleAuthData struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

func NewAuthHandler() *Handler {
	return &Handler{
		userUseCase:  usecase.NewUserUseCase(),
		groupUseCase: usecase.NewGroupUseCase(),
	}
}

func (h *Handler) Register(c *gin.Context) {
	var dto RegisterRequest
	if err := pkg.MapToDto(c, &dto); err != nil {
		pkg.HandleError(c, err)
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)

	user := &entity.User{
		Email:     dto.Email,
		Password:  string(hashedPassword),
		FirstName: dto.Firstname,
		LastName:  dto.Lastname,
	}

	if err := h.userUseCase.Register(user); err != nil {
		pkg.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user register success"})
}

func (h *Handler) Login(c *gin.Context) {
	var dto LoginRequest
	if err := pkg.MapToDto(c, &dto); err != nil {
		pkg.HandleError(c, err)
		return
	}

	user := h.userUseCase.GetByEmail(dto.Email)
	if user == nil {
		pkg.HandleError(c, exception.NewUnauthorized("email or password incorrect"))
		return
	}
	if !user.IsActive {
		pkg.HandleError(c, exception.NewUnauthorized("user is not active"))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password)); err != nil {
		pkg.HandleError(c, exception.NewUnauthorized("email or password incorrect"))
		return
	}

	HandleResponse(*user, c)
}

func (h *Handler) RefreshToken(c *gin.Context) {
	user := pkg.GetUserFromReq(c)
	HandleResponse(user, c)
}
