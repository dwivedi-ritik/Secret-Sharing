package globals

import (
	"github.com/dwivedi-ritik/text-share-be/lib"
	"gorm.io/gorm"
)

var RedisCache *lib.RedisCache
var DB *gorm.DB
