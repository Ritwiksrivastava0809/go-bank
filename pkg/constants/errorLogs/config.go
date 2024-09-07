package errorLogs

const (
	ParsingError          = "got error %v on parsing configuration file"
	BindingJsonError      = "got error while binding the request: %s"
	GetCreateAccountError = "got error while creating the account: %s"
	GetAccountError       = "got error while getting the account: %s"
	UpdateAccountError    = "got error while updating the account: %s"
	ValidationError       = "got error while validating the request body: %v"
	InvalidUserError      = "got error while getting the user: %s"
	AddAmountError        = "got error while adding the amount: %s"
)
