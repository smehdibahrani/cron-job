package group

type CreateGroupDTO struct {
	Name        string `json:"name" binding:"required"`
	TagName     string `json:"tagName" `
	Description string `json:"description"`
}
