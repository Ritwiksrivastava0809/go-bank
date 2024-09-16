package utils

import "time"

func GetAccessTokenDuration() time.Duration {
	return 15 * time.Minute // Ensure this matches your config
}
