package server

import (
	"context"
	"log"
	"log/slog"
	"net/http"

	"github.com/dwivedi-ritik/text-share-be/apps/auth"
	"github.com/dwivedi-ritik/text-share-be/apps/private"
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
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{TranslateError: true})
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("DB Connection Established")

	db.AutoMigrate(&models.Message{}, &models.UniqueId{}, &models.User{}, &models.Encryption{})
	SyncUpUniqueIds(db)

	return db

}

func dBContextMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "DB", DB)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CreateServer() *http.ServeMux {
	DB = initializeDB()
	mainRouter := http.NewServeMux()
	mainRouter.Handle("/api/public/", dBContextMiddleware(middleware.Logger((public.PublicRouter("/api/public/")))))
	mainRouter.Handle("/api/user/", dBContextMiddleware(middleware.Logger((auth.UserRouter("/api/user/")))))
	mainRouter.Handle("/api/private/", dBContextMiddleware(middleware.Logger(private.PrivateRouter("/api/private/"))))

	return mainRouter
}
