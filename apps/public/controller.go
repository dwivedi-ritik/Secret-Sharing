package public

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/dwivedi-ritik/text-share-be/models"
)

func AddMessageController(w http.ResponseWriter, r *http.Request) {
	messageService := MessageService{DB: models.DB}
	var message models.Message
	json.NewDecoder(r.Body).Decode(&message)

	_, err := message.ValidateModel()
	if err != nil {
		MessageServiceErrorHandler(err, w, r)
		return
	}
	messageService.AddMessage(&message)

	w.WriteHeader(http.StatusCreated)

	serverResponse := struct {
		UniqueId uint32
	}{UniqueId: message.UniqueIdentifier}

	jsonData, err := json.Marshal(serverResponse)

	if err != nil {
		log.Fatal(err)
	}
	w.Write(jsonData)
}

func GetMessageController(w http.ResponseWriter, r *http.Request) {
	messageService := MessageService{DB: models.DB}

	queryId := r.URL.Query().Get("id")
	w.Header().Add("content-type", "application/json")

	if len(queryId) != 0 {
		queryIdint, err := strconv.Atoi(queryId)
		if err != nil {
			MessageServiceErrorHandler(ErrIdParamIsMissing, w, r)
			return
		}

		message, gormError := messageService.GetMessageById(uint32(queryIdint))
		_, err = message.IsMessageExpired()
		if err != nil {
			MessageServiceErrorHandler(err, w, r)
			return
		}
		if gormError != nil {
			MessageServiceErrorHandler(gormError, w, r)
			return
		}

		jsonData, err := json.Marshal(message)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal(err)
		}
		w.Write(jsonData)
		return
	}
	messages := messageService.GetAllMessage()
	jsonData, err := json.Marshal(messages)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
	}
	w.Write(jsonData)

}
