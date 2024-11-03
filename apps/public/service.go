package public

import (
	"strconv"

	"github.com/dwivedi-ritik/text-share-be/lib"
	"github.com/dwivedi-ritik/text-share-be/models"
	"github.com/dwivedi-ritik/text-share-be/types"
	"gorm.io/gorm"
)

type MessageService struct {
	redisCache *lib.RedisCache
	dB         *gorm.DB
}

func (messageService *MessageService) AddMessage(message *models.Message) bool {

	message.Save()
	digitToMap := 0
	for {
		newDigit := lib.RandomDigitGenerate(1000, 9999)
		cacheHits := messageService.redisCache.IsExists(strconv.Itoa(newDigit))
		if !cacheHits {
			digitToMap = newDigit
			break
		}
	}
	message.UniqueIdentifier = uint32(digitToMap)
	messageCache := types.MessageCache{
		Key:   strconv.Itoa(digitToMap),
		Value: message.Id,
	}
	messageService.redisCache.AddValue(&messageCache)
	return true
}

func (messageService *MessageService) GetMessageById(identity string, expire bool) (models.Message, error) {
	var message models.Message

	isExists := messageService.redisCache.IsExists(identity)
	if !isExists {
		return message, ErrInvalidUniqueId
	}
	actualId := messageService.redisCache.FetchValue(&types.MessageCache{Key: identity})
	err := messageService.dB.First(&message, actualId).Error

	if err != nil {
		return message, err
	}

	if !expire {
		return message, nil
	}

	messageService.dB.Model(&message).Update("expired", true)
	messageService.redisCache.DeleteValue(&types.MessageCache{Key: identity})
	return message, nil
}
