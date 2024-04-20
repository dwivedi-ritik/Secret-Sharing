package models

type UniqueId struct {
	Id        uint32 `gorm:"primaryKey"`
	Identity  uint32 `gorm:"not null;unique"`
	Available bool   `gorm:"default:true"`
	Queued    bool   `gorm:"default:false"`
}
