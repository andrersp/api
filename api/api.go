package api

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

const (
	HeaderContentType = "Content-Type"
)

const (
	MIMEApplicationJSON = "application/json"
)

type MiddlewareFunc func(HandlerFunc) HandlerFunc

type HandlerFunc func(c Context) error

type HTTPErrorHandler func(err error, c Context) error

type Api struct {
	server      *http.Server
	middlewares []MiddlewareFunc
	routes      []*Route
	Binder      Binder

	ErrorHandler HTTPErrorHandler
}

func (a *Api) Use(middleware ...MiddlewareFunc) {
	a.middlewares = append(a.middlewares, middleware...)
}

func (a *Api) DefaultErrorHandler(err error, c Context) error {
	fmt.Println("aaaa")
	return c.Json(400, map[string]string{"err": err.Error()})
}

func (a *Api) Start(address string) error {

	a.server = &http.Server{
		Addr: address,
	}
	listner, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	defer listner.Close()

	for _, route := range a.routes {
		path := fmt.Sprintf("%s %s", route.Method, route.Path)
		http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Aa")
			ctx := &context{
				request:  r,
				response: NewResponse(w, a),
				api:      a,
			}
			handler := applyMiddleware(route.Handler, a.middlewares...)
			err := handler(ctx)
			fmt.Println(err)
		})
		log.Printf("Registered routes with method: %s and path %s\n", route.Method, route.Path)
	}

	return a.server.Serve(listner)

}

func New() *Api {
	middlewares := make([]MiddlewareFunc, 0)
	middlewares = append(middlewares, loggerMiddeware)

	a := &Api{middlewares: middlewares}
	a.ErrorHandler = a.DefaultErrorHandler
	a.Binder = &DefaultBinder{}

	return a

}

func (a *Api) NewContext(w http.ResponseWriter, r *http.Request) Context {
	return &context{
		request:  r,
		response: NewResponse(w, a),
		handler:  notFoundHandler,
		api:      a,
	}

}
