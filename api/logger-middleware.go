package api

import (
	"fmt"
	"log/slog"
	"time"
)

func loggerMiddeware(next HandlerFunc) HandlerFunc {
	return func(c Context) error {
		start := time.Now()
		path := c.Request().URL.Path
		method := c.Request().Method
		c.Set("teste", 123)
		next(c)
		status := c.Response().Status
		end := time.Since(start)

		slog.Info("request", "duration", end.Seconds(), "method", method, "path", path, "status", status)

		return nil
	}

}

func JsonMiddleware(next HandlerFunc) HandlerFunc {
	return func(c Context) error {
		fmt.Println("Start JSon")
		r := next(c)

		fmt.Println("End Json")
		return r
	}
}
