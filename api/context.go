package api

import (
	"encoding/json"
	"net/http"
)

type context struct {
	request  *http.Request
	response http.ResponseWriter
	api      *Api
}

// Bind implements HttpContext.
func (c *context) Bind(i interface{}) error {
	panic("unimplemented")
}

func (c *context) Request() *http.Request {
	return c.request
}

func (c *context) Response() http.ResponseWriter {
	return c.response
}

func (c *context) Json(code int, i interface{}) error {
	c.writeContentType(MIMEApplicationJSON)
	c.response.WriteHeader(code)
	return json.NewEncoder(c.response).Encode(i)
}

func (c *context) writeContentType(value string) {
	header := c.Response().Header()
	if header.Get(HeaderContentType) == "" {
		header.Set(HeaderContentType, value)
	}
}

func NewContext(w http.ResponseWriter, r *http.Request) HttpContext {
	return &context{
		request:  r,
		response: w,
	}
}