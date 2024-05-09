package api

import (
	"log"
	"net"
	"net/http"
	"sync"
)

const (
	HeaderContentType = "Content-Type"
)

const (
	MIMEApplicationJSON = "application/json"
)

type MiddlewareFunc func(HandlerFunc) HandlerFunc

type HandlerFunc func(c Context) error

type Api struct {
	server      *http.Server
	middlewares []MiddlewareFunc
	Routes      []*Router
	Binder      Binder
	pool        sync.Pool
}

func (a *Api) Use(middleware ...MiddlewareFunc) {
	a.middlewares = append(a.middlewares, middleware...)
}

func (a *Api) Start(address string) error {
	a.server = &http.Server{
		Addr:    address,
		Handler: a,
	}
	listner, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	defer listner.Close()

	for _, route := range a.Routes {
		log.Printf("Registered routes with method: %s and path %s\n", route.Method, route.Path)
	}

	return a.server.Serve(listner)

}

func (a *Api) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	c := a.pool.Get().(*context)
	for _, router := range a.Routes {
		if r.Method == router.Method && r.URL.Path == router.Path {
			c.request = r
			c.response = w
			handler := applyMiddleware(router.Handler, a.middlewares...)
			handler(c)
			return
		}
	}

	a.pool.Put(c)

	http.NotFound(w, r)

}

func New() *Api {
	middlewares := make([]MiddlewareFunc, 0)
	middlewares = append(middlewares, loggerMiddeware)

	a := &Api{middlewares: middlewares}
	a.pool.New = func() any {
		return &context{
			request:  nil,
			response: nil,
			handler:  notFoundHandler,
		}

	}

	return a

}
