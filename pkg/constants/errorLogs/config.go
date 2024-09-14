package errorLogs

const (
	ParsingError            = "got error %v on parsing configuration file"
	BindingJsonError        = "got error while binding the request. error :: %s"
	GetCreateAccountError   = "got error while creating the account. error :: %s"
	GetAccountError         = "got error while getting the account. error :: %s"
	UpdateAccountError      = "got error while updating the account. error :: %s"
	ValidationError         = "got error while validating the request body. error :: %v"
	InvalidUserError        = "got error while getting the user. error :: %s"
	AddAmountError          = "got error while adding the amount. error :: %s"
	ListAccountsError       = "got error while getting list of the accounts. error :: %s"
	TransferTxError         = "got error while generating the transaction. error :: %s"
	HashPasswordError       = "got error while hashing the password. error :: %v"
	CreateUserError         = "got error while creating the user. error :: %s"
	InvalidEmailFormatError = "got error while parsing the email. error :: %s"
	CheckUserExistenceError = "got error while checking the user existence. error :: %s"
	TransactionError        = "got error while starting the transaction. error :: %s"
	CommitTransactionError  = "got error while committing the transaction. error :: %s"
	InvalidKeySize          = "invalid key size : must be atleast %d characters long"
	TokenError              = "got error while generating the token. error :: %s"
	GetUserError            = "got error while getting the user. error :: %s"
)
