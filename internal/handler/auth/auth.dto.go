package auth

type RegisterRequest struct {
	Firstname string `json:"firstName" binding:"required,min=3"`
	Lastname  string `json:"lastName" binding:"required,min=3"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}
