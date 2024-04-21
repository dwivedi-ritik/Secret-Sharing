package models

import (
	"errors"
	"time"
)

var (
	ErrMessageBodyIsEmpty     error = errors.New("message body is empty")
	ErrMessageBodySizeInvalid error = errors.New("message body size is invalid")
	ErrMessageExpired         error = errors.New("message exceeded it duration")
)

type Message struct {
	Id               uint64 `gorm:"primaryKey"`
	Content          string `gorm:"not null"`
	Encrypted        bool   `gorm:"default:false"`
	EncryptedType    string
	Duration         int64     `gorm:"default:1800000"`
	CreatedBy        uint64    `gorm:"default:null"`
	UniqueIdentifier uint32    `gorm:"not null"`
	Expired          bool      `gorm:"default:false"`
	CreatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func (m *Message) ValidateModel() (bool, error) {
	if len(m.Content) <= 5 {
		return false, ErrMessageBodySizeInvalid
	}

	return true, nil
}

func (m *Message) IsMessageExpired() (bool, error) {
	currentTime := time.Now().UnixMilli()
	messageCreatedAt := m.CreatedAt.UnixMilli()
	if m.Duration < (currentTime - messageCreatedAt) {
		return true, ErrMessageExpired
	}
	return false, nil
}

//TBD Users
