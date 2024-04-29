package private

import (
	"encoding/json"
	"net/http"

	"github.com/dwivedi-ritik/text-share-be/models"
	"gorm.io/gorm"
)

func GetUserInfoController(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(UserKey{}).(models.User)
	db := r.Context().Value("DB").(*gorm.DB)
	var privateService = PrivateService{DB: db}

	userDto, err := privateService.GetUserInfo(&user)
	if err != nil {
		PrivateErrorHandler(err, w, r)
		return
	}
	jsonData, err := json.Marshal(userDto)

	if err != nil {
		PrivateErrorHandler(err, w, r)
		return
	}
	w.Header().Add("content-type", "application/json")
	w.Write(jsonData)
}

func GetUserMessagesController(w http.ResponseWriter, r *http.Request) {

	DB := r.Context().Value("DB").(*gorm.DB)
	privateService := PrivateService{DB: DB}
	userMessages := privateService.GetUserMessages(r)

	serializedMessages, err := json.Marshal(userMessages)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}
	w.Header().Add("content-type", "application/json")
	w.Write(serializedMessages)
}

func EncryptMessageController(w http.ResponseWriter, r *http.Request) {

}

func DecryptMessageController(w http.ResponseWriter, r *http.Request) {

}
