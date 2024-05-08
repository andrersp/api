package api

import (
	"fmt"
	"log/slog"
	"net/http"
)

const (
	HeaderContentType = "Content-Type"
)

const (
	MIMEApplicationJSON = "application/json"
)

type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

type HttpContext interface {
	Request() *http.Request
	Response() http.ResponseWriter
	Json(code int, i interface{}) error
	Bind(i interface{}) error
}

type HandlerFunc func(c HttpContext) error

type Api struct {
	server      *http.Server
	middlewares []MiddlewareFunc
	Binder      Binder
}

func (a *Api) Add(method, path string, handler HandlerFunc) {

	http.HandleFunc(fmt.Sprintf("%s %s", method, path), func(w http.ResponseWriter, r *http.Request) {
		ctx := NewContext(w, r)
		if err := handler(ctx); err != nil {
			fmt.Println("Ok")
		}

	})

}

func (a *Api) Start(address string) error {

	a.server.Addr = address
	slog.Info("start http server", "address", address)

	return a.server.ListenAndServe()
}

func New() *Api {
	server := new(http.Server)
	middlewares := make([]MiddlewareFunc, 0)
	middlewares = append(middlewares, loggerMiddeware)

	return &Api{server: server, middlewares: middlewares}

}

// func getEnv(value, defaultValue string) string {
// 	value = os.Getenv(value)
// 	if value == "" {
// 		return defaultValue
// 	}
// 	return value
// }
