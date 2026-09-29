package middleware

import (
	"time"

	echov4 "github.com/labstack/echo/v4"

	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/logger"
)

func RequestLogger() echov4.MiddlewareFunc {
	return func(next echov4.HandlerFunc) echov4.HandlerFunc {
		return func(c echov4.Context) error {
			// Step 1: time the whole request
			start := time.Now()
			err := next(c)

			// Step 2: write the error response now, otherwise the status we log would still be 200
			if err != nil {
				c.Error(err)
			}

			// Step 3: log method, path and how long it took
			logger.Info("http request",
				"request_id", c.Response().Header().Get(echov4.HeaderXRequestID),
				"method", c.Request().Method,
				"path", c.Request().URL.Path,
				"status", c.Response().Status,
				"duration", time.Since(start).String(),
			)
			return nil
		}
	}
}
