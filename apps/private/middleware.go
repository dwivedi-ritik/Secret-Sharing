package private

import (
	"context"
	"net/http"
	"strings"

	"github.com/dwivedi-ritik/text-share-be/apps/auth"
	"github.com/dwivedi-ritik/text-share-be/lib"
	"gorm.io/gorm"
)

type UserKey struct{}

func Authorization(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearerToken := r.Header.Get("Authorization")
		if len(bearerToken) < 7 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
			return
		}
		token := strings.TrimPrefix(bearerToken, "Bearer ")

		if !lib.ValidateToken(token) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
			return
		}
		claims, err := lib.GetUnverifiedClaims(token)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
			return
		}
		dB := r.Context().Value("DB")
		userService := auth.UserService{DB: dB.(*gorm.DB)}

		user, err := userService.FetchUserByUserName(claims.Username)
		if err != nil {
			return
		}

		ctx := context.WithValue(r.Context(), UserKey{}, user)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
