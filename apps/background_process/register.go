package background_process

import "net/http"

func BackgroundProcessRouter(pathPrefix string) *http.ServeMux {
	backgroundProcessRouter := http.NewServeMux()
	backgroundProcessRouter.HandleFunc("GET "+pathPrefix+"status", GetProcessStatus)
	return backgroundProcessRouter
}
