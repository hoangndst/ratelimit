package tokenbucket

import (
	"context"
	"github.com/hoangndst/ratelimit/drivers"
	"github.com/redis/go-redis/v9"
	"testing"
)

func TestNewTokenBucket(t *testing.T) {
	rbs := &redis.Client{}
	tb := NewTokenBucket(10, 5, 2, rbs)

	if tb.capacity != 10 {
		t.Errorf("Expected capacity to be 10, got %d", tb.capacity)
	}

	if tb.refillAmount != 5 {
		t.Errorf("Expected refillAmount to be 5, got %d", tb.refillAmount)
	}

	if tb.timeBetweenSlots != 2 {
		t.Errorf("Expected timeBetweenSlots to be 2, got %d", tb.timeBetweenSlots)
	}

	if tb.rbs != rbs {
		t.Errorf("Expected rbs to be %v, got %v", rbs, tb.rbs)
	}
}

type mockRediser struct {
	drivers.Rediser
	RunFunc func(ctx context.Context, rbs drivers.Rediser, keys []string, args ...interface{}) redis.Cmder
}

func (m *mockRediser) Run(ctx context.Context, rbs drivers.Rediser, keys []string, args ...interface{}) redis.Cmder {
	return m.RunFunc(ctx, rbs, keys, args...)
}

func TestRateLimit(t *testing.T) {
	rbs := &mockRediser{
		RunFunc: func(ctx context.Context, rbs drivers.Rediser, keys []string, args ...interface{}) redis.Cmder {
			return redis.NewIntCmd(ctx, 1000) // Mock the return value of the Lua script
		},
	}
	tb := NewTokenBucket(10, 5, 2, rbs)

	sleepTime, err := tb.RateLimit(context.Background(), "testKey")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if sleepTime != 1.0 {
		t.Errorf("Expected sleepTime to be 1.0, got %f", sleepTime)
	}
}
