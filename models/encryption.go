package models

type Encryption struct {
	PrivateKey string
	PublicKey  string
	Active     bool   `gorm:"default:true"`
	Type       string `gorm:"default:RSA"`
	CreatedBy  uint64
}
