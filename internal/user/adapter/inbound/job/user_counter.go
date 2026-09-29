package job

import (
	"context"
	"time"

	"github.com/Genedevelop/7solutions-backend-challenge/internal/user/port/inbound"
	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/logger"
)

func RunUserCounter(ctx context.Context, counter inbound.ICountUsersUseCase, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		// Step 1: stop cleanly on shutdown
		case <-ctx.Done():
			logger.Info("user counter stopped")
			return
		// Step 2: every tick, count with a short timeout so a slow database cannot pile up ticks
		case <-ticker.C:
			countCtx, cancel := context.WithTimeout(ctx, interval/2)
			n, err := counter.Count(countCtx)
			cancel()
			if err != nil {
				logger.Error("count users failed", err)
				continue
			}
			logger.Info("total users", "count", n)
		}
	}
}
