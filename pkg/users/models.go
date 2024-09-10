package users

type User struct {
	UserName string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
	Password string `json:"password" binding:"required" validate:"min=8,max=20"`
}
