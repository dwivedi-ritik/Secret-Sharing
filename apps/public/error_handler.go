package public

import (
	"errors"
	"net/http"

	"github.com/dwivedi-ritik/text-share-be/models"
	"gorm.io/gorm"
)

func MessageServiceErrorHandler(err error, w http.ResponseWriter, r *http.Request) {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(gorm.ErrRecordNotFound.Error()))

	case errors.Is(err, ErrIdParamIsMissing):
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(ErrIdParamIsMissing.Error()))

	case errors.Is(err, models.ErrMessageBodyIsEmpty):
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(models.ErrMessageBodyIsEmpty.Error()))

	case errors.Is(err, models.ErrMessageBodySizeInvalid):
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(models.ErrMessageBodySizeInvalid.Error()))

	case errors.Is(err, models.ErrMessageExpired):
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(models.ErrMessageExpired.Error()))

	case errors.Is(err, ErrInvalidUniqueId):
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(ErrInvalidUniqueId.Error()))
	}

}
