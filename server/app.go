package server

import (
	"context"
	"log"
	"log/slog"
	"net/http"

	"github.com/dwivedi-ritik/text-share-be/apps/auth"
	"github.com/dwivedi-ritik/text-share-be/apps/background_process"
	"github.com/dwivedi-ritik/text-share-be/apps/private"
	"github.com/dwivedi-ritik/text-share-be/apps/public"
	"github.com/dwivedi-ritik/text-share-be/middleware"
	"github.com/dwivedi-ritik/text-share-be/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var RedisClient *redis.Client

func initializeRedisCache() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

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
	RedisClient = initializeRedisCache()
	ctx := context.Background()
	err := RedisClient.Set(ctx, "Hello", "World", 0).Err()
	if err != nil {
		panic(err)
	}
	mainRouter := http.NewServeMux()
	mainRouter.Handle("/api/public/", dBContextMiddleware(middleware.Logger((public.PublicRouter("/api/public/")))))
	mainRouter.Handle("/api/user/", dBContextMiddleware(middleware.Logger((auth.UserRouter("/api/user/")))))
	mainRouter.Handle("/api/private/", dBContextMiddleware(middleware.Logger(private.PrivateRouter("/api/private/"))))
	mainRouter.Handle("/api/process/", dBContextMiddleware(middleware.Logger(background_process.BackgroundProcessRouter("/api/process/"))))
	return mainRouter
}
