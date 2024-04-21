package private

import (
	"github.com/dwivedi-ritik/text-share-be/models"
	"gorm.io/gorm"
)

type UserDto struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type PrivateService struct {
	DB *gorm.DB
}

func (privateService *PrivateService) GetAllUserMessages(limit uint8, offset uint8) {

}

func (privateService *PrivateService) GetUserInfo(user *models.User) (UserDto, error) {
	return UserDto{Username: user.Username, Email: user.Email}, nil
}

func (privateService *PrivateService) GetMessageById(messageId uint32) {

}
