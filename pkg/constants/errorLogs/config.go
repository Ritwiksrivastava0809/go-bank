package errorLogs

const (
	ParsingError          = "got error %v on parsing configuration file"
	BindingJsonError      = "got error while binding the request. error :: %s"
	GetCreateAccountError = "got error while creating the account. error :: %s"
	GetAccountError       = "got error while getting the account. error :: %s"
	UpdateAccountError    = "got error while updating the account. error :: %s"
	ValidationError       = "got error while validating the request body. error :: %v"
	InvalidUserError      = "got error while getting the user. error :: %s"
	AddAmountError        = "got error while adding the amount. error :: %s"
	ListAccountsError     = "got error while getting list of the accounts. error :: %s"
	TransferTxError       = "got error while generating the transaction. error :: %s"
)
