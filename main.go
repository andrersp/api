package main

import (
	"api/api"
	"context"
	"fmt"
)

func PrintAlgo(c context.Context) {
	fmt.Println(c)
}

type Response struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func Home(a api.Context) error {
	PrintAlgo(a.Request().Context())

	response := Response{Name: "Testando", Age: 21}
	return a.Json(201, response)
}

func main() {

	server := api.New()
	server.Add("GET", "/home", Home, api.JsonMiddleware)
	server.Add("GET", "/home2", Home)
	server.Add("GET", "/home3", Home)
	server.Add("GET", "/home4", Home)
	server.Add("GET", "/home5", Home)
	server.Add("GET", "/hom4e", Home)
	server.Start(":8080")
}
