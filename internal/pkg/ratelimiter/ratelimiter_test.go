package ratelimiter

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func newTestLimiter(t *testing.T, maxAttempts int, blockDuration time.Duration) (*Limiter, *miniredis.Miniredis) {
	t.Helper()

	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})

	return New(client, maxAttempts, blockDuration), mr
}

func TestLimiter_CheckAndIncrement_AllowsUnderThreshold(t *testing.T) {
	limiter, _ := newTestLimiter(t, 3, time.Minute)
	ctx := context.Background()

	for i := range 3 {
		blocked, err := limiter.CheckAndIncrement(ctx, "john")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if blocked {
			t.Fatalf("attempt %d: expected not blocked", i+1)
		}
	}
}

func TestLimiter_CheckAndIncrement_BlocksAfterThreshold(t *testing.T) {
	limiter, _ := newTestLimiter(t, 3, time.Minute)
	ctx := context.Background()

	for range 3 {
		if _, err := limiter.CheckAndIncrement(ctx, "john"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	blocked, err := limiter.CheckAndIncrement(ctx, "john")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !blocked {
		t.Fatal("expected login to be blocked after exceeding threshold")
	}
}

func TestLimiter_CheckAndIncrement_AlreadyBlockedShortCircuits(t *testing.T) {
	limiter, mr := newTestLimiter(t, 1, time.Minute)
	ctx := context.Background()

	// First call: under threshold, not blocked yet.
	if _, err := limiter.CheckAndIncrement(ctx, "john"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Second call: crosses the threshold, gets blocked.
	blocked, err := limiter.CheckAndIncrement(ctx, "john")
	if err != nil || !blocked {
		t.Fatalf("expected blocked=true, err=nil, got blocked=%v err=%v", blocked, err)
	}

	attemptsAtBlock, err := mr.Get("attempt:john")
	if err != nil {
		t.Fatalf("unexpected error reading attempt key: %v", err)
	}

	// Third call: already blocked, must short-circuit before touching the counter.
	blocked, err = limiter.CheckAndIncrement(ctx, "john")
	if err != nil || !blocked {
		t.Fatalf("expected blocked=true, err=nil, got blocked=%v err=%v", blocked, err)
	}

	attemptsAfterBlock, err := mr.Get("attempt:john")
	if err != nil {
		t.Fatalf("unexpected error reading attempt key: %v", err)
	}
	if attemptsAfterBlock != attemptsAtBlock {
		t.Fatalf("expected attempt counter to stay at %s once blocked, got %s", attemptsAtBlock, attemptsAfterBlock)
	}
}

func TestLimiter_Reset_ClearsState(t *testing.T) {
	limiter, mr := newTestLimiter(t, 1, time.Minute)
	ctx := context.Background()

	if _, err := limiter.CheckAndIncrement(ctx, "john"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	blocked, err := limiter.CheckAndIncrement(ctx, "john")
	if err != nil || !blocked {
		t.Fatal("expected login to be blocked before reset")
	}

	if err := limiter.Reset(ctx, "john"); err != nil {
		t.Fatalf("unexpected error resetting: %v", err)
	}

	if mr.Exists("blocked:john") {
		t.Fatal("expected blocked key to be cleared after reset")
	}
	if mr.Exists("attempt:john") {
		t.Fatal("expected attempt key to be cleared after reset")
	}

	blocked, err = limiter.CheckAndIncrement(ctx, "john")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if blocked {
		t.Fatal("expected login to not be blocked after reset")
	}
}

func TestLimiter_AttemptCounterExpiresAfterWindow(t *testing.T) {
	limiter, mr := newTestLimiter(t, 2, time.Minute)
	ctx := context.Background()

	if _, err := limiter.CheckAndIncrement(ctx, "john"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mr.FastForward(time.Minute + time.Second)

	if mr.Exists("attempt:john") {
		t.Fatal("expected attempt counter to expire after the block window")
	}
}
