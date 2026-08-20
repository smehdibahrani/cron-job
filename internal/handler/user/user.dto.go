package user

type UpdateRequest struct {
	FirstName string `json:"firstName" binding:"required,min=3"`
	LastName  string `json:"lastName" binding:"required,min=3"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required,min=6"`
	NewPassword     string `json:"newPassword" binding:"required,min=8"`
}
