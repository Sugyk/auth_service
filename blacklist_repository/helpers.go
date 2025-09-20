package blacklist_repository

import "time"

func getMaxAttemptsCount() int {
	return 10
}

func getBlockExpiration() time.Duration {
	return time.Minute * 5
}
