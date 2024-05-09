package main

import (
	"api/api"
	"context"
	"fmt"
)

func PrintAlgo(c context.Context) {
	fmt.Println(c)
}

type Payload struct {
	Name string `json:"name"`
}

type Response struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func Home(a api.Context) error {
	// fmt.Println(a.Request().PathValue())
	PrintAlgo(a.Request().Context())

	response := Response{Name: "Testando", Age: 21}
	return a.Json(201, response)
}

func Save(a api.Context) error {
	var payload Payload

	err := a.Bind(&payload)
	if err != nil {
		return err
	}
	return a.Json(201, payload)
}

func main() {

	server := api.New()
	server.Add("GET", "/home/{nome}", Home, api.JsonMiddleware)
	server.Add("POST", "/home", Save)

	server.Start(":8080")
}
