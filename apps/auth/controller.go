package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dwivedi-ritik/text-share-be/lib"
	"github.com/dwivedi-ritik/text-share-be/models"
	"gorm.io/gorm"
)

type UserDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"Email"`
}

func (userDto *UserDTO) Validate() error {
	if len(userDto.Username) == 0 && len(userDto.Email) == 0 {
		return ErrLoginCredsInvalid
	}
	return nil
}

func AddUserController(w http.ResponseWriter, r *http.Request) {

	DB := r.Context().Value("DB")
	var user models.User
	var userSerivce = UserService{DB: DB.(*gorm.DB)}
	json.NewDecoder(r.Body).Decode(&user)
	user.Password = lib.CreatePasswordHash(user.Password)
	err := userSerivce.AddUser(user)

	if err != nil {
		AuthErrorHandler(err, w, r)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(http.StatusText(http.StatusCreated)))
}

func LoginUserController(w http.ResponseWriter, r *http.Request) {
	DB := r.Context().Value("DB")
	var userDto UserDTO
	var userSerivce = UserService{DB: DB.(*gorm.DB)}
	json.NewDecoder(r.Body).Decode(&userDto)

	err := userDto.Validate()

	if err != nil {
		AuthErrorHandler(err, w, r)
		return
	}

	user, err := userSerivce.FetchUser(models.User{Username: userDto.Username, Email: userDto.Email})

	if err != nil {
		AuthErrorHandler(err, w, r)
		return
	}

	if !lib.ComparePasswordHash(user.Password, userDto.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
		return
	}

	jwtToken, err := lib.CreateNewToken(lib.UserToken{Username: userDto.Username, Email: user.Email})

	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("content-type", "application/json")

	jsonData, err := json.Marshal(struct {
		Token string `json:"token"`
	}{Token: jwtToken})
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)

}
