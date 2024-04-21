package auth

import (
	"errors"
	"log/slog"
	"net/http"

	"gorm.io/gorm"
)

func AuthErrorHandler(err error, w http.ResponseWriter, r *http.Request) {
	switch {
	case errors.Is(err, ErrLoginCredsInvalid):
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(ErrLoginCredsInvalid.Error()))
	case errors.Is(err, gorm.ErrRecordNotFound):
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(gorm.ErrRecordNotFound.Error()))
	case errors.Is(err, gorm.ErrDuplicatedKey):
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(gorm.ErrDuplicatedKey.Error()))
	default:
		slog.Info("Error", err.Error(), "occurred")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong"))
	}
}
