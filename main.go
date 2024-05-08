package main

import (
	"api/api"
)

func Home(a api.HttpContext) error {
	return a.Json(201, map[string]bool{"sucess": true})
}

func main() {

	api := api.New()
	api.Add("GET", "/home", Home)
	api.Add("GET", "/home2", Home)
	api.Add("GET", "/home3", Home)
	api.Add("GET", "/home4", Home)
	api.Add("GET", "/home5", Home)
	api.Add("GET", "/hom4e", Home)
	api.Start(":8080")
}
