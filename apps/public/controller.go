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

	serverResponse := struct {
		UniqueId uint32
	}{UniqueId: message.UniqueIdentifier}

	jsonData, err := json.Marshal(serverResponse)

	if err != nil {
		log.Fatal(err)
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)
}

func GetMessageController(w http.ResponseWriter, r *http.Request) {
	messageService := MessageService{DB: models.DB}

	queryId := r.URL.Query().Get("id")
	expired := r.URL.Query().Get("expired")

	w.Header().Add("content-type", "application/json")

	if len(queryId) != 0 {
		queryIdint, err := strconv.Atoi(queryId)
		if err != nil {
			MessageServiceErrorHandler(ErrIdParamIsMissing, w, r)
			return
		}

		message, err := messageService.GetMessageById(uint32(queryIdint))

		if err != nil {
			MessageServiceErrorHandler(err, w, r)
			return
		}

		if len(expired) != 0 && expired == "true" {
			messageService.ExpireUniqueIds(uint32(queryIdint))
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
