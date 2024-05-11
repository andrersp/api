package api

import (
	"log/slog"
	"time"
)

func loggerMiddeware(next HandlerFunc) HandlerFunc {
	return func(c Context) error {
		start := time.Now()
		path := c.Request().URL.Path
		method := c.Request().Method
		next(c)
		status := c.Response().Status
		end := time.Since(start)

		slog.Info("request", "duration", end.Seconds(), "method", method, "path", path, "status", status)

		return nil
	}

}
