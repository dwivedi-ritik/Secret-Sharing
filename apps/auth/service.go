package auth

import (
	"github.com/dwivedi-ritik/text-share-be/models"
	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func (userService *UserService) AddUser(user models.User) error {
	id := userService.DB.Create(&user)
	if id.Error != nil {
		return id.Error
	}
	return nil
}

func (userSerivce *UserService) FetchUserByUserName(username string) (models.User, error) {
	var user models.User
	err := userSerivce.DB.Model(&models.User{}).Where("username = ?", username).First(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (userSerivce *UserService) FetchUser(user models.User) (models.User, error) {
	var fetchedUser models.User
	err := userSerivce.DB.Where("username = ?", user.Username).Or("email = ?", user.Email).First(&fetchedUser).Error
	if err != nil {
		return fetchedUser, err
	}
	return fetchedUser, nil
}
