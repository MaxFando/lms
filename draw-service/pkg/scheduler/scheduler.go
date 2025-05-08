package scheduler

import (
	"context"
	"fmt"
	"time"
)

func Schedule(ctx context.Context, f func(ctx context.Context) error, interval time.Duration) error {
	fmt.Println("aboba!")
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
