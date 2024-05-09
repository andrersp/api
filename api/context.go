package api

import (
	"encoding/json"
	"net/http"
	"sync"
)

type Map map[string]interface{}

type Context interface {
	Request() *http.Request
	Response() *Response
	Json(code int, i interface{}) error
	Bind(i interface{}) error
	Set(key string, val interface{})
	Get(key string) interface{}
	Reset(r *http.Request, w http.ResponseWriter)
	Api() *Api
}

type context struct {
	request  *http.Request
	response *Response
	api      *Api
	handler  HandlerFunc
	store    Map
	lock     sync.RWMutex
}

func (c *context) Get(key string) interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.store[key]
}

// Set implements Context.
func (c *context) Set(key string, val interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.store == nil {
		c.store = make(Map)
	}
	c.store[key] = val
}

func (c *context) Api() *Api {
	return c.api
}

// Bind implements Context.
func (c *context) Bind(i interface{}) error {
	return c.api.Binder.Bind(i, c)
}

func (c *context) Request() *http.Request {
	return c.request
}

func (c *context) Response() *Response {
	return c.response
}

func (c *context) Json(code int, i interface{}) error {
	c.writeContentType(MIMEApplicationJSON)
	c.Response().WriteHeader(code)
	return json.NewEncoder(c.Response()).Encode(i)
}

func (c *context) writeContentType(value string) {
	header := c.Response().Header()
	if header.Get(HeaderContentType) == "" {
		header.Set(HeaderContentType, value)
	}
}

func (c *context) Reset(r *http.Request, w http.ResponseWriter) {
	c.request = r
	c.response.reset(w)
	// c.query = nil
	// c.handler = NotFoundHandler
	c.store = nil
	// c.path = ""
	// c.pnames = nil
	// c.logger = nil
	// NOTE: Don't reset because it has to have length c.echo.maxParam (or bigger) at all times

}
