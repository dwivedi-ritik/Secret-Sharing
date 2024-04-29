package private

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/dwivedi-ritik/text-share-be/models"
	"gorm.io/gorm"
)

type UserDto struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func Paginate(r *http.Request) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		pageSize := 25
		pageNumber := 0
		pageNumber, err := strconv.Atoi(r.URL.Query().Get("pageNumber"))
		if err != nil {
			pageNumber = 0
		}

		pageSize, err = strconv.Atoi(r.URL.Query().Get("pageSize"))

		if err != nil {
			pageSize = 25
		}
		switch {
		case pageSize > 25:
			pageSize = 25
		case pageSize <= 0:
			pageSize = 0
		}

		offset := (pageNumber - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

type PrivateService struct {
	DB *gorm.DB
}

func (privateService *PrivateService) GetUserMessages(r *http.Request) []models.Message {
	var messages []models.Message
	loggedInUser := r.Context().Value(UserKey{}).(models.User)
	slog.Info("User context", loggedInUser.Email, "found")
	privateService.DB.Scopes(Paginate(r)).Where("created_by = ?", loggedInUser.Id).Find(&messages)
	return messages
}

func (privateService *PrivateService) GetUserInfo(user *models.User) (UserDto, error) {
	return UserDto{Username: user.Username, Email: user.Email}, nil
}

func (privateService *PrivateService) GetMessageById(messageId uint32) {

}
