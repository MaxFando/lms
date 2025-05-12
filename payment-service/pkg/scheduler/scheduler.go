package scheduler

import (
	"context"
	"time"
)

func Schedule(ctx context.Context, f func(ctx context.Context) error, interval time.Duration) error {
	if err := f(ctx); err != nil {
		return err
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := f(ctx); err != nil {
				return err
			}
			ticker.Reset(interval)
		case <-ctx.Done():
			return nil
		}
	}
}
