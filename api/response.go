// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: © 2015 LabStack LLC and api contributors

package api

import (
	"bufio"
	"errors"
	"net"
	"net/http"
)

// Response wraps an http.ResponseWriter and implements its interface to be used
// by an HTTP handler to construct an HTTP response.
// See: https://golang.org/pkg/net/http/#ResponseWriter
type Response struct {
	api       *Api
	Writer    http.ResponseWriter
	Status    int
	Size      int64
	Committed bool
}

// NewResponse creates a new instance of Response.
func NewResponse(w http.ResponseWriter, e *Api) (r *Response) {
	return &Response{Writer: w, api: e}
}

func (r *Response) Header() http.Header {
	return r.Writer.Header()
}

func (r *Response) WriteHeader(code int) {
	r.Status = code
	r.Writer.WriteHeader(r.Status)
	r.Committed = true
}

func (r *Response) Write(b []byte) (n int, err error) {
	n, err = r.Writer.Write(b)
	r.Size += int64(n)

	return
}

// See [http.Flusher](https://golang.org/pkg/net/http/#Flusher)
func (r *Response) Flush() {
	err := responseControllerFlush(r.Writer)
	if err != nil && errors.Is(err, http.ErrNotSupported) {
		panic(errors.New("response writer flushing is not supported"))
	}
}

// See [http.Hijacker](https://golang.org/pkg/net/http/#Hijacker)
func (r *Response) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return responseControllerHijack(r.Writer)
}

// Unwrap returns the original http.ResponseWriter.
// ResponseController can be used to access the original http.ResponseWriter.
// See [https://go.dev/blog/go1.20]
func (r *Response) Unwrap() http.ResponseWriter {
	return r.Writer
}

func (r *Response) reset(w http.ResponseWriter) {
	r.Writer = w
	r.Size = 0
	r.Status = http.StatusOK
	r.Committed = false
}

func responseControllerHijack(rw http.ResponseWriter) (net.Conn, *bufio.ReadWriter, error) {
	return http.NewResponseController(rw).Hijack()
}

func responseControllerFlush(rw http.ResponseWriter) error {
	return http.NewResponseController(rw).Flush()
}
