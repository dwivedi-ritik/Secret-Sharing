package server

import (
	"net/http"

	"github.com/dwivedi-ritik/text-share-be/apps/public"
	"github.com/dwivedi-ritik/text-share-be/middleware"
	"github.com/dwivedi-ritik/text-share-be/models"
)

func CreateServer() *http.ServeMux {
	models.Init()
	mainRouter := http.NewServeMux()
	mainRouter.Handle("/public/", middleware.Logger(public.PublicRouter("/public/")))
	return mainRouter
}
