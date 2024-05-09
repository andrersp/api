package api

func notFoundHandler(a Context) error {
	return a.Json(201, map[string]bool{"sucess": true})
}
