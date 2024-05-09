package api

type Router struct {
	Method  string
	Path    string
	Handler HandlerFunc
}

func (a *Api) Add(method, path string, handler HandlerFunc, middlewares ...MiddlewareFunc) *Router {
	router := a.add(method, path, handler, middlewares...)
	a.Routes = append(a.Routes, router)
	return router
}

func (a *Api) add(method, path string, handler HandlerFunc, middlewares ...MiddlewareFunc) *Router {
	h := applyMiddleware(handler, middlewares...)
	return &Router{Method: method, Path: path, Handler: h}
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
