package constants

const DefaultConfigurationType = "yaml"
const DefaultConfigurationPath = "environment"

// db related constants
const (
	DBDriver = "postgres"
)

// server related constants
const (
	Origin        = "origin"
	ContentLength = "Content-Length"
	ContentType   = "Content-Type"
	Authorization = "Authorization"
)

// handler related constants
const (
	ConstantDB    = "db"
	UserID        = "X-UserID"
	InternalToken = "X-Internal-Token"
	PageLimit     = "limit"
	PageOffset    = "offset"
	SortBy        = "sortBy"
	ID            = "id"
	SortDirection = "sortDirection"
	Ascending     = "asc"
	Descending    = "desc"
)

// currency related constants
const (
	USD  = "USD"
	EUR  = "EUR"
	INR  = "INR"
	CAD  = "CAD"
	YEN  = "YEN"
	RUB  = "RUB"
	YUAN = "YUAN"
	PESO = "PESO"
)
