package api

import (
	"fmt"
)

func loggerMiddeware(next HandlerFunc) HandlerFunc {
	return func(c Context) error {
		fmt.Println("Start Middleware")
		r := next(c)
		fmt.Println("End Middleware")
		return r
	}

}

func JsonMiddleware(next HandlerFunc) HandlerFunc {
	return func(c Context) error {
		fmt.Println("Start JSon")
		r := next(c)
		fmt.Println("End Json")
		return r
	}
}
