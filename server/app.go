package server

import (
	"context"
	"log"
	"log/slog"
	"net/http"

	"github.com/dwivedi-ritik/text-share-be/apps/public"
	"github.com/dwivedi-ritik/text-share-be/middleware"
	"github.com/dwivedi-ritik/text-share-be/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func initializeDB() *gorm.DB {
	dbUrl := "postgres://localhost:5432/textshare"
	slog.Info("Initializing DB Connection")
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("DB Connection Established")

	db.AutoMigrate(&models.Message{}, &models.UniqueId{})
	SyncUpUniqueIds(db)

	return db

}

func DBContextMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "DB", DB)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CreateServer() *http.ServeMux {
	DB = initializeDB()
	mainRouter := http.NewServeMux()
	mainRouter.Handle("/public/", DBContextMiddleware(middleware.Logger((public.PublicRouter("/public/")))))

	return mainRouter
}
