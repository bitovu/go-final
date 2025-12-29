package server

import (
	"final/pkg/api"
	"net/http"
)

func Run() error {
	api.Init()
	http.Handle("/", http.FileServer(http.Dir("web")))
	return http.ListenAndServe(":7540", nil)
}
