package users

type CreateUser struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required" validate:"oneof=USD EUR INR"`
}
