package accounts

type CreateAccount struct {
	Currency string `json:"currency" binding:"required,currency"`
}

type UpdateBalance struct {
	Balance int64 `json:"balance" binding:"required" validate:"gt=0"`
}
