package tokenbucket

import (
	"context"
	"github.com/hoangndst/ratelimit/drivers"
	"time"
)

type TokenBucket struct {
	capacity         int
	refillAmount     int
	timeBetweenSlots int
	rbs              drivers.Rediser
}

func NewTokenBucket(capacity, refillAmount, timeBetweenSlots int, rbs drivers.Rediser) *TokenBucket {
	return &TokenBucket{
		capacity:         capacity,
		refillAmount:     refillAmount,
		timeBetweenSlots: timeBetweenSlots,
		rbs:              rbs,
	}
}

func (tb *TokenBucket) createRedisTime() (int64, int64) {
	now := time.Now()
	secondsPart := now.Unix()
	microsecondsPart := now.UnixNano()/1000 - secondsPart*1_000_000
	return secondsPart, microsecondsPart
}

func (tb *TokenBucket) parseTimestamp(timestamp int64) (float64, error) {
	wakeupTime := time.Unix(0, timestamp*1_000_000)
	now := time.Now()
	if wakeupTime.Before(now) {
		return 0, nil
	}
	sleepTime := wakeupTime.Sub(now).Seconds()
	return sleepTime, nil
}

// RateLimit checks if the rate limit is exceeded for the given key.
// If the rate limit is not exceeded, the function returns the time
// in seconds until the next request can be made.
// If the rate limit is exceeded, the function returns 0.
func (tb *TokenBucket) RateLimit(ctx context.Context, key string) (float64, error) {
	seconds, microseconds := tb.createRedisTime()
	args := []interface{}{
		tb.capacity,
		tb.refillAmount,
		tb.timeBetweenSlots,
		seconds,
		microseconds,
	}
	slot, err := tokenBucketScript.Run(ctx, tb.rbs, []string{key}, args...).Result()
	if err != nil {
		return 0, err
	}
	return tb.parseTimestamp(slot.(int64))
}
