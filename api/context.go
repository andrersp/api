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
	Param(string) string
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

func (c *context) Param(name string) string {
	return c.request.PathValue(name)
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
