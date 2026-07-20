package ratelimiter

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Limiter struct {
	client        *redis.Client
	maxAttempts   int
	blockDuration time.Duration
}

func New(client *redis.Client, maxAttempts int, blockDuration time.Duration) *Limiter {
	return &Limiter{
		client:        client,
		maxAttempts:   maxAttempts,
		blockDuration: blockDuration,
	}
}

// CheckAndIncrement reports whether login is currently blocked. If it is not,
// it counts this call as a failed attempt: the counter's TTL is set to
// blockDuration on the first attempt so the window decays on its own, and
// crossing maxAttempts blocks the login for blockDuration.
func (l *Limiter) CheckAndIncrement(ctx context.Context, login string) (bool, error) {
	blockedKey := "blocked:" + login

	blocked, err := l.client.Exists(ctx, blockedKey).Result()
	if err != nil {
		return false, fmt.Errorf("checking block status: %w", err)
	}
	if blocked > 0 {
		return true, nil
	}

	attemptKey := "attempt:" + login

	count, err := l.client.Incr(ctx, attemptKey).Result()
	if err != nil {
		return false, fmt.Errorf("incrementing attempt count: %w", err)
	}
	if count == 1 {
		if err := l.client.Expire(ctx, attemptKey, l.blockDuration).Err(); err != nil {
			return false, fmt.Errorf("setting attempt window: %w", err)
		}
	}

	if count > int64(l.maxAttempts) {
		if err := l.client.Set(ctx, blockedKey, "blocked", l.blockDuration).Err(); err != nil {
			return false, fmt.Errorf("setting block: %w", err)
		}
		return true, nil
	}

	return false, nil
}

// Reset clears both the attempt counter and any active block for login.
func (l *Limiter) Reset(ctx context.Context, login string) error {
	if err := l.client.Del(ctx, "blocked:"+login, "attempt:"+login).Err(); err != nil {
		return fmt.Errorf("resetting throttle state: %w", err)
	}
	return nil
}
