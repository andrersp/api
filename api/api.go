package api

import (
	"fmt"
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

type HTTPErrorHandler func(err error, c Context)

type Api struct {
	server       *http.Server
	middlewares  []MiddlewareFunc
	Routes       []*Router
	Binder       Binder
	pool         sync.Pool
	ErrorHandler HTTPErrorHandler
}

func (a *Api) Use(middleware ...MiddlewareFunc) {
	a.middlewares = append(a.middlewares, middleware...)
}

func (a *Api) DefaultErrorHandler(err error, c Context) {
	c.Json(400, map[string]string{"err": err.Error()})
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

	for _, route := range a.Routes {
		path := fmt.Sprintf("%s %s", route.Method, route.Path)
		http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			ctx := &context{
				request:  r,
				response: NewResponse(w, a),
			}
			handler := applyMiddleware(route.Handler, a.middlewares...)
			if err := handler(ctx); err != nil {
				fmt.Println("Error")

			}
		})
		log.Printf("Registered routes with method: %s and path %s\n", route.Method, route.Path)
	}

	return a.server.Serve(listner)

}

// func (a *Api) ServeHTTP(w http.ResponseWriter, r *http.Request) {

// 	c := a.pool.Get().(*context)
// 	c.Reset(r, w)
// 	var handler HandlerFunc
// 	for _, router := range a.Routes {
// 		if r.Method == router.Method && r.URL.Path == router.Path {
// 			handler = applyMiddleware(router.Handler, a.middlewares...)
// 			break
// 		}
// 	}

// 	if handler != nil {
// 		if err := handler(c); err != nil {
// 			a.ErrorHandler(err, c)
// 		}
// 	} else {
// 		http.NotFound(w, r)
// 	}

// 	a.pool.Put(c)

// 	//

// }

func New() *Api {
	middlewares := make([]MiddlewareFunc, 0)
	middlewares = append(middlewares, loggerMiddeware)

	a := &Api{middlewares: middlewares}
	a.pool.New = func() any {
		return a.NewContext(nil, nil)

	}
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
