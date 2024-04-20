package models

import (
	"errors"
	"log"
	"log/slog"

	lib "github.com/dwivedi-ritik/text-share-be/libs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

var RedisAlternative = make(map[uint32]uint64)

var UniqueIdsDeque = lib.Deque{}

func Init() {
	dbUrl := "postgres://localhost:5432/textshare"
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&Message{}, &UniqueId{})

	//upsertUniqueIdsRecords(db)
	updateUniqueIdsDeque(db)
	DB = db

}

// Run to update deque range
// func upsertUniqueIdsRecords(db *gorm.DB) {
// 	const MAX_SIZE uint16 = 1000
// 	const START_RANGE uint16 = 1000
// 	const LAST_RANGE uint16 = START_RANGE + MAX_SIZE

// 	slog.Info("Started Unique Ids Update")
// 	var uniqueIds [MAX_SIZE]UniqueId
// 	for element := START_RANGE; element <= LAST_RANGE; element++ {
// 		index := element % MAX_SIZE
// 		uniqueIds[index].Identity = uint32(element)
// 		uniqueIds[index].Available = true
// 		uniqueIds[index].Queued = true

// 		RedisAlternative[uint32(element)] = 0

// 		UniqueIdsDeque.AddFront(uint32(element))

// 	}

// 	db.Clauses(clause.OnConflict{OnConstraint: "uni_unique_ids_identity", DoNothing: true}).CreateInBatches(&uniqueIds, 500)

// 	slog.Info("Unique ids updation completed")
// }

// Necessary to make deque sync with DB value while server restart
func updateUniqueIdsDeque(db *gorm.DB) {
	slog.Info("Deque Update Started")
	var (
		missed  uint8 = 0
		updated uint8 = 0
	)
	var availableUniqueIds []UniqueId
	var uniqueIdsToSync []UniqueId
	db.Model(&UniqueId{}).Where("available = ?", true).Find(&availableUniqueIds)
	slog.Info("Queue addition of available ids started from front")

	for _, uniqueId := range availableUniqueIds {
		RedisAlternative[uniqueId.Identity] = 0
		UniqueIdsDeque.AddFront(uniqueId.Identity)
	}

	slog.Info("Deque Update completed for available ids")

	db.Model(&UniqueId{}).Where("available = ?", false).Find(&uniqueIdsToSync)
	for _, uniqueId := range uniqueIdsToSync {
		var message Message

		err := db.Model(&message).Where("unique_identifier = ?", uniqueId.Identity).First(&message).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Info("Message", "Id", message.Id, " lost, Updating map value to", 0)
			RedisAlternative[uniqueId.Identity] = 0
			missed++
		} else {
			slog.Info("Message", "Id", message.Id, " Found, Updating map value", message.Id)
			RedisAlternative[uniqueId.Identity] = message.Id
			updated++
		}
		UniqueIdsDeque.AddRear(uniqueId.Identity)
	}

	slog.Info("Deque Update completed for sync up")
	slog.Info("Total Map value", "Missed", missed)
	slog.Info("Total Map value", "Updated", updated)
}
