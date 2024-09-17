package utils

import (
	"log"
	"time"

	"github.com/Ritwiksrivastava0809/go-bank/pkg/config"
)

func GetAccessTokenDuration() time.Duration {
	// Read the duration from a configuration file or environment variable
	durationStr := config.GetAccessTokenDuration()
	if durationStr == "" {
		durationStr = "15m" // Default to 15 minutes if not configured
	}
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		// Handle the error, e.g., log it and return a default value
		log.Printf("Failed to parse ACCESS_TOKEN_DURATION: %v", err)
		return 15 * time.Minute
	}
	return duration
}
