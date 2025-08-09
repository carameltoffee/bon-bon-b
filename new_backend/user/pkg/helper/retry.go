package helper

import (
	"context"
	"fmt"
	"time"
)

func Retry(ctx context.Context, attempts int, delay time.Duration, fn func() error) error {
	for i := 0; i < attempts; i++ {
		err := fn()
		if err == nil {
			return nil
		}

		select {
		case <-time.After(delay * time.Duration(i+1)):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return fmt.Errorf("all %d retry attempts failed", attempts)
}
