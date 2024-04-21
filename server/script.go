package server

import (
	"errors"
	"log/slog"

	"github.com/dwivedi-ritik/text-share-be/lib"
	"github.com/dwivedi-ritik/text-share-be/models"
	"gorm.io/gorm"
)

// func upsertUniqueIdsRecords(db *gorm.DB) {
// 	// const MAX_SIZE uint16 = 1000
// 	// const START_RANGE uint16 = 1000
// 	// const LAST_RANGE uint16 = START_RANGE + MAX_SIZE

// 	// slog.Info("Started Unique Ids Update")
// 	// var uniqueIds [MAX_SIZE]UniqueId
// 	// for element := START_RANGE; element <= LAST_RANGE; element++ {
// 	// 	index := element % MAX_SIZE
// 	// 	uniqueIds[index].Identity = uint32(element)
// 	// 	uniqueIds[index].Available = true
// 	// 	uniqueIds[index].Queued = true

// 	// 	RedisAlternative[uint32(element)] = 0

// 	// 	UniqueIdsDeque.AddFront(uint32(element))

// 	// }

// 	db.Clauses(clause.OnConflict{OnConstraint: "uni_unique_ids_identity", DoNothing: true}).CreateInBatches(&uniqueIds, 500)

// 	slog.Info("Unique ids updation completed")
// }

// Necessary to make deque sync with DB value while server restart
func SyncUpUniqueIds(db *gorm.DB) {
	slog.Info("Deque Update Started")
	var (
		missed  uint8 = 0
		updated uint8 = 0
	)
	var availableUniqueIds []models.UniqueId
	var uniqueIdsToSync []models.UniqueId
	db.Model(&models.UniqueId{}).Where("available = ?", true).Find(&availableUniqueIds)
	slog.Info("Queue addition of available ids started from front")

	for _, uniqueId := range availableUniqueIds {
		lib.RedisAlternative[uniqueId.Identity] = 0
		lib.UniqueIdsDeque.AddFront(uniqueId.Identity)
	}

	slog.Info("Deque Update completed for available ids")

	db.Model(&models.UniqueId{}).Where("available = ?", false).Find(&uniqueIdsToSync)
	for _, uniqueId := range uniqueIdsToSync {

		var message models.Message
		err := db.Model(&message).Where("unique_identifier = ?", uniqueId.Identity).First(&message).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Info("Message", "Id", uniqueId.Identity, "Lost, Updating map value to", 0)
			lib.RedisAlternative[uniqueId.Identity] = 0
			missed++
		} else {
			slog.Info("Message", "Id", uniqueId.Identity, "Found, Updating map value", message.Id)
			lib.RedisAlternative[uniqueId.Identity] = message.Id
			updated++
		}
		lib.UniqueIdsDeque.AddRear(uniqueId.Identity)
	}

	slog.Info("Deque Update completed for sync up")
	slog.Info("Total Map value", "Missed", missed)
	slog.Info("Total Map value", "Updated", updated)
}
