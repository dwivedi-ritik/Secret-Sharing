package public

import (
	"net/http"

	"github.com/dwivedi-ritik/text-share-be/globals"
)

func PublicRouter(pathPrefix string) *http.ServeMux {
	publicRouter := http.NewServeMux()

	controller := PublicMessageController{
		messageService: &MessageService{
			redisCache: globals.RedisCache,
			dB:         globals.DB,
		},
	}
	publicRouter.HandleFunc("GET "+pathPrefix+"get/message", controller.GetMessageController)
	publicRouter.HandleFunc("POST "+pathPrefix+"add/message", controller.AddMessageController)
	return publicRouter
}
