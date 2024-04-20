package public

import (
	"net/http"
)

// Register to register the controller with corresponding routes
func PublicRouter(pathPrefix string) *http.ServeMux {
	publicRouter := http.NewServeMux()
	publicRouter.HandleFunc("POST "+pathPrefix+"add/message", AddMessageController)
	publicRouter.HandleFunc("GET "+pathPrefix+"get/message", GetMessageController)
	return publicRouter
}
