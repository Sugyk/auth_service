package blacklist_repository

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func (b *BLRepo) CheckAndIncrement(login string) (bool, error) {
	ctx := context.Background()
	var is_blocked string
	err := b.db.Get(ctx, "blocked:"+login).Scan(&is_blocked)
	if err != nil && err != redis.Nil {
		return false, fmt.Errorf("error getting blocked: %w", err)
	}
	if is_blocked == "blocked" {
		return true, nil
	}
	count, err := b.db.Incr(ctx, "attempt:"+login).Result()
	if err != nil {
		return false, fmt.Errorf("error increment attempts: %w", err)
	}
	if count > int64(getMaxAttemptsCount()) {
		b.db.Set(ctx, "blocked:"+login, "blocked", getBlockExpiration())
		return true, nil
	}
	return false, nil
}
