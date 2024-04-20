package models

import (
	"log"
	"log/slog"

	lib "github.com/dwivedi-ritik/text-share-be/libs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var DB *gorm.DB

var RedisAlternative = make(map[uint32]int64)

var UniqueIdsDeque = lib.Deque{}

func Init() {
	dbUrl := "postgres://localhost:5432/textshare"
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&Message{}, &UniqueId{})
	upsertUniqueIdsRecords(db)
	DB = db

}

func upsertUniqueIdsRecords(db *gorm.DB) {
	const MAX_SIZE uint16 = 1000
	const START_RANGE uint16 = 1000
	const LAST_RANGE uint16 = START_RANGE + MAX_SIZE

	slog.Info("Started Unique Ids Update")
	var uniqueIds [MAX_SIZE]UniqueId
	for element := START_RANGE; element <= LAST_RANGE; element++ {
		index := element % MAX_SIZE
		uniqueIds[index].Identity = uint32(element)
		uniqueIds[index].Available = true
		uniqueIds[index].Queued = true

		RedisAlternative[uint32(element)] = 0

		UniqueIdsDeque.AddFront(uint32(element))

	}

	db.Clauses(clause.OnConflict{OnConstraint: "uni_unique_ids_identity", DoNothing: true}).CreateInBatches(&uniqueIds, 500)

	slog.Info("Unique ids updation completed")
}
