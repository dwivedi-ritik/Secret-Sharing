package private

import (
	"errors"
	"log/slog"
	"net/http"

	"gorm.io/gorm"
)

func PrivateErrorHandler(err error, w http.ResponseWriter, r *http.Request) {
	switch {
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
