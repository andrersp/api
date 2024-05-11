package main

import (
	"api/api"
	"context"
	"fmt"
	"log/slog"
	"os"
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
	fmt.Println(a.Param("nome"))

	response := Response{Name: "Testando", Age: 21}
	return a.Json(201, response)
}

func Save(a api.Context) error {
	var payload Payload

	err := a.Bind(&payload)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return a.Json(201, payload)
}

func setDefaultLogger() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)
}

func main() {
	setDefaultLogger()

	server := api.New()
	server.Add("GET", "/home/{nome}", Home)
	server.Add("POST", "/home", Save)

	server.Start(":8080")
}
