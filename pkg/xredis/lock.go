package xredis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/xybor/todennus-backend/pkg/xcontext"
)

type Locker struct {
	key        string
	client     *redis.Client
	exipration time.Duration
}

func NewLock(redis *redis.Client, key string, exipration time.Duration) *Locker {
	return &Locker{
		client:     redis,
		key:        key,
		exipration: exipration,
	}
}

func (l *Locker) LockFunc(ctx context.Context, f func() error) error {
	if err := l.Lock(ctx); err != nil {
		return err
	}

	defer l.Unlock(ctx)

	if err := f(); err != nil {
		return err
	}

	return nil
}

func (l *Locker) Lock(ctx context.Context) error {
	for {
		ok, err := l.client.SetNX(ctx, l.key, "", l.exipration).Result()
		if err != nil {
			return fmt.Errorf("cannot lock by redis: %w", err)
		}

		if ok {
			break
		}

		time.Sleep(time.Second)
		continue
	}

	return nil
}

func (l *Locker) Unlock(ctx context.Context) {
	_, err := l.client.Del(ctx, l.key).Result()
	if err != nil {
		xcontext.Logger(ctx).Warn("failed-to-release-redis-lock", "key", l.key, "err", err)
	}
}
