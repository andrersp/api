package api

import "encoding/json"

type Binder interface {
	Bind(i interface{}, c Context)
}

type DefaultBinder struct {
}

func (b *DefaultBinder) BindBody(c Context, i interface{}) error {
	req := c.Request()

	if req.ContentLength == 0 {
		return nil
	}

	return json.NewDecoder(c.Request().Body).Decode(i)

}

func (b *DefaultBinder) Bind(i interface{}, c Context) error {
	return b.BindBody(c, i)

}
