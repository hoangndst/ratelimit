package tokenbucket

import (
	"context"
	"github.com/hoangndst/ratelimit/drivers"
	"github.com/redis/go-redis/v9"
	"testing"
)

func TestNewTokenBucket(t *testing.T) {
	// Create rediser mock object
	//rbs := &mockRediser{}
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
	tb := NewTokenBucket(10, 5, 2, 10, rbs)

	sleepTime, err := tb.RateLimit(context.Background(), "testKey")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if sleepTime != 1.0 {
		t.Errorf("Expected sleepTime to be 1.0, got %f", sleepTime)
	}
}
