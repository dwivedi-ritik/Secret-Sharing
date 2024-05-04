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
	Id               uint64    `json:"-" gorm:"primaryKey"`
	Content          string    `json:"content" gorm:"not null"`
	Encrypted        bool      `json:"encrypted" gorm:"default:false"`
	EncryptedType    string    `json:"-"`
	Duration         int64     `json:"duration" gorm:"default:1800000"`
	CreatedBy        uint64    `json:"-" gorm:"default:null"`
	UniqueIdentifier uint32    `json:"identifier" gorm:"not null"`
	Expired          bool      `json:"expired" gorm:"default:false"`
	CreatedAt        time.Time `json:"createdAt" gorm:"default:CURRENT_TIMESTAMP"`
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
