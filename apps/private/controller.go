package private

import (
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"

	"github.com/dwivedi-ritik/text-share-be/lib"
	"github.com/dwivedi-ritik/text-share-be/models"
	"gorm.io/gorm"
)

type MessageEncryptionDto struct {
	ReceiverUserId uint64 `json:"receiverUserId"`
	Message        string `json:"message"`
}

type EncryptionKeyDto struct {
	Privatekey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
}

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

	var messageEncryption MessageEncryptionDto
	loggedInUser := r.Context().Value(UserKey{}).(models.User)
	DB := r.Context().Value("DB").(*gorm.DB)
	privateService := PrivateService{DB: DB}
	err := json.NewDecoder(r.Body).Decode(&messageEncryption)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		log.Fatal(err)
		return
	}

	receiverUser, err := privateService.GetUserEncryptionByUserId(messageEncryption.ReceiverUserId)

	if err != nil {
		log.Fatal(err)
		return
	}

	privateKeyBytes, err := hex.DecodeString(receiverUser.Encryption.PrivateKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		log.Fatal(err)
		return
	}

	publicKeyBytes, err := hex.DecodeString(receiverUser.Encryption.PublicKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		log.Fatal(err)
		return
	}

	currentEncrytion := lib.RSAEncryption{
		Keys: lib.EncryptionKey{
			PrivateKey: privateKeyBytes,
			PublicKey:  publicKeyBytes,
		},
	}

	cipher := currentEncrytion.EncryptMessage([]byte(messageEncryption.Message))

	storedEncryptedMessage := privateService.AddEncryptedMessage(hex.EncodeToString(cipher), loggedInUser.Id, "RSA")

	jsonBytes, err := json.Marshal(&storedEncryptedMessage.UniqueIdentifier)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		log.Fatal(err)
		return
	}
	w.Header().Add("content-type", "application/json")
	w.Write(jsonBytes)

}

// func DecryptMessageController(w http.ResponseWriter, r *http.Request) {
// 	loggedInUser := r.Context().Value(UserKey{}).(models.User)
// 	DB := r.Context().Value("DB").(*gorm.DB)
// 	privateService := PrivateService{DB: DB}
// 	messageService := public.MessageService{DB: DB}
// 	receiverUser, err := privateService.GetUserEncryptionByUserId(loggedInUser.Id)
// 	if err != nil {
// 		log.Fatal(err)
// 		return
// 	}

// 	privateKeyBytes, err := hex.DecodeString(receiverUser.Encryption.PrivateKey)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
// 		log.Fatal(err)
// 		return
// 	}

// 	publicKeyBytes, err := hex.DecodeString(receiverUser.Encryption.PublicKey)

// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
// 		log.Fatal(err)
// 		return
// 	}

// 	currentEncrytion := lib.RSAEncryption{
// 		Keys: lib.EncryptionKey{
// 			PrivateKey: privateKeyBytes,
// 			PublicKey:  publicKeyBytes,
// 		},
// 	}

// 	queryId := r.URL.Query().Get("id")
// 	expired := r.URL.Query().Get("expired")

// 	w.Header().Add("content-type", "application/json")

// 	if len(queryId) != 0 {
// 		queryIdint, err := strconv.Atoi(queryId)
// 		if err != nil {
// 			log.Fatal(err)
// 			return
// 		}
// 		message, err := messageService.GetMessageById(uint32(queryIdint))

// 		if err != nil {
// 			PrivateErrorHandler(err, w, r)
// 			return
// 		}

// 		messageBytes, err := hex.DecodeString(message.Content)
// 		if err != nil {
// 			log.Fatal(err)
// 			return
// 		}

// 		message.Content = string(currentEncrytion.DecryptMessage(messageBytes))

// 		if len(expired) != 0 && expired == "true" {
// 			messageService.ExpireUniqueIds(uint32(queryIdint))
// 		}
// 		jsonData, err := json.Marshal(message)
// 		if err != nil {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			log.Fatal(err)
// 		}
// 		w.Write(jsonData)
// 		return
// 	}

// }

func GenerateUserEncryptionKeyController(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(UserKey{}).(models.User)
	DB := r.Context().Value("DB").(*gorm.DB)
	rsaEncryption := lib.NewRSAEncryption()
	privateService := PrivateService{DB: DB}
	err := privateService.AddEncryptionKeysOfUser(rsaEncryption, &user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		log.Fatal(err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(http.StatusText(http.StatusCreated)))
}
