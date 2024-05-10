package api

type Route struct {
	Method  string
	Path    string
	Handler HandlerFunc
}

func (a *Api) Add(method, path string, handler HandlerFunc, middlewares ...MiddlewareFunc) *Route {
	router := a.add(method, path, handler, middlewares...)
	a.routes = append(a.routes, router)
	return router
}

func (a *Api) add(method, path string, handler HandlerFunc, middlewares ...MiddlewareFunc) *Route {
	h := applyMiddleware(handler, middlewares...)
	return &Route{Method: method, Path: path, Handler: h}
}

func applyMiddleware(h HandlerFunc, middlewares ...MiddlewareFunc) HandlerFunc {
	if len(middlewares) < 1 {
		return h
	}
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h

}
