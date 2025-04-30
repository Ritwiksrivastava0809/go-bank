package newRelic

// NewRelicLog represents the structure of the log payload
type NewRelicLog struct {
	Common struct {
		Attributes map[string]string `json:"attributes"`
	} `json:"common"`
	Logs []struct {
		Timestamp  int64                  `json:"timestamp"`
		Message    string                 `json:"message"`
		Attributes map[string]interface{} `json:"attributes,omitempty"`
	} `json:"logs"`
}
