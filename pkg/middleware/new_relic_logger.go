package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// New Relic API Key
const newRelicAPIKey = "8e816dd73b01fe07b7376bed16bbb0acFFFFNRAL"

// LogPayload represents the New Relic log structure
type LogPayload struct {
	Common struct {
		Attributes map[string]string `json:"attributes"`
	} `json:"common"`
	Logs []struct {
		Timestamp  int64                  `json:"timestamp"`
		Message    string                 `json:"message"`
		Attributes map[string]interface{} `json:"attributes"`
	} `json:"logs"`
}

// NewRelicLoggerMiddleware automatically sends logs to New Relic for every API request
func NewRelicLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Process the request
		c.Next()

		// Collect request and response info
		latency := time.Since(startTime)
		statusCode := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		clientIP := c.ClientIP()

		// Prepare log message
		logMessage := fmt.Sprintf("Handled request: method=%s path=%s status=%d latency=%v", method, path, statusCode, latency)

		// Create log payload
		logData := LogPayload{}
		logData.Common.Attributes = map[string]string{
			"service":  "go-bank",
			"hostname": "server-01",
			"logtype":  "api-request",
		}
		logData.Logs = append(logData.Logs, struct {
			Timestamp  int64                  `json:"timestamp"`
			Message    string                 `json:"message"`
			Attributes map[string]interface{} `json:"attributes"`
		}{
			Timestamp: time.Now().UnixMilli(),
			Message:   logMessage,
			Attributes: map[string]interface{}{
				"method":    method,
				"path":      path,
				"status":    statusCode,
				"latency":   latency.String(),
				"client_ip": clientIP,
			},
		})

		// Send log to New Relic
		go SendLogToNewRelic(logData)
	}
}

// SendLogToNewRelic sends log data to New Relic
func SendLogToNewRelic(logData LogPayload) {
	// Convert to JSON
	jsonData, err := json.Marshal(logData)
	if err != nil {
		fmt.Println("Failed to marshal log data:", err)
		return
	}

	// Print log payload to terminal (for debugging)
	fmt.Println("Sending log to New Relic:")
	fmt.Println(string(jsonData))

	// Make HTTP POST request
	url := "https://log-api.newrelic.com/log/v1"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", newRelicAPIKey)
	req.Header.Set("Accept", "*/*")

	// Send request
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to send log:", err)
		return
	}
	defer resp.Body.Close()

	// Read response
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Response Status: %d\n", resp.StatusCode)
	fmt.Println("Response Body:", string(body))

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		fmt.Println("Error: Failed to send log to New Relic")
	} else {
		fmt.Println("Log sent successfully to New Relic")
	}
}
