package public

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/dwivedi-ritik/text-share-be/models"
)

type PublicMessageController struct {
	messageService *MessageService
}

func (messageController *PublicMessageController) AddMessageController(w http.ResponseWriter, r *http.Request) {
	messageService := messageController.messageService
	var message models.Message

	json.NewDecoder(r.Body).Decode(&message)

	_, err := message.ValidateModel()
	if err != nil {
		MessageServiceErrorHandler(err, w, r)
		return
	}
	messageService.AddMessage(&message)
	serverResponse := struct {
		UniqueId uint32 `json:"uniqueId"`
	}{UniqueId: message.UniqueIdentifier}

	jsonData, err := json.Marshal(serverResponse)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)
}

func (messageController *PublicMessageController) GetMessageController(w http.ResponseWriter, r *http.Request) {
	messageService := messageController.messageService
	queryId := r.URL.Query().Get("id")
	expired := r.URL.Query().Get("expired")
	expireMessage := expired == "true"

	_, err := strconv.Atoi(queryId)
	if err != nil {
		messageController.renderErrorResponse(w, http.StatusBadRequest, nil)
		return
	}

	message, err := messageService.GetMessageById(queryId, expireMessage)

	if err != nil {
		messageController.renderErrorResponse(w, http.StatusBadRequest, ErrInvalidUniqueId)
		return
	}
	jsonData, err := json.Marshal(message)
	if err != nil {
		messageController.renderErrorResponse(w, http.StatusBadRequest, ErrInvalidUniqueId)
		return
	}
	messageController.renderResponse(w, &jsonData, http.StatusOK)

}

func (messageController *PublicMessageController) renderResponse(w http.ResponseWriter, data *[]byte, status int) {
	w.Header().Set("Content-Type", "application/json")
	if status != 0 {
		w.WriteHeader(status)
	}
	w.Write(*data)
}

func (messageController *PublicMessageController) renderErrorResponse(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	w.Write([]byte(err.Error()))
}
