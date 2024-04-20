package public

import (
	"log"

	"github.com/dwivedi-ritik/text-share-be/models"
	"gorm.io/gorm"
)

type MessageService struct {
	DB *gorm.DB
}

func (messageService *MessageService) AddMessage(message *models.Message) *models.Message {

	uniqueId := models.UniqueIdsDeque.RemoveFront()

	message.UniqueIdentifier = uniqueId.(uint32) //type assertion

	id := messageService.DB.Create(&message)

	models.RedisAlternative[uint32(uniqueId.(uint32))] = uint64(message.Id) //type asseration

	messageService.DB.Model(&models.UniqueId{}).Where("identity = ?", uniqueId).Update("available", false).Update("queued", false)

	if id.Error != nil {
		log.Fatal(id.Error)
	}
	return message
}

func (messageService *MessageService) GetAllMessage() []models.Message {

	var messages []models.Message
	messageService.DB.Find(&messages)
	return messages
}

func (messageService *MessageService) GetMessageById(identity uint32) (models.Message, error) {
	var message models.Message

	messageId := models.RedisAlternative[identity]
	err := messageService.DB.First(&message, messageId).Error
	if err != nil {
		return message, err
	}
	_, err = message.IsMessageExpired()

	if err != nil {
		messageService.ExpireUniqueIds(identity)
		return message, err
	}

	return message, nil
}

func (messageService *MessageService) ExpireUniqueIds(identity uint32) {
	messageService.DB.Model(&models.UniqueId{}).Where("identity = ?", identity).Update("available", true).Update("queued", true)
	models.RedisAlternative[identity] = 0
	models.UniqueIdsDeque.AddRear(identity)
}
