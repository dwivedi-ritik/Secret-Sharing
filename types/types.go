package types

import (
	"errors"
	"time"
)

var (
	ErrMessageCacheNotFound error = errors.New("key isn't present or expired")
	ErrKeyNotFound          error = errors.New("key not found in cache")
)

type MessageCache struct {
	Key      string
	Value    interface{}
	Duration time.Duration
}

// Get new message cache type
// Key has to be string, Value could be anytype
// Duration is seconds
func NewMessageCahe(key string, value interface{}, duration int64) *MessageCache {
	return &MessageCache{
		Key:      key,
		Value:    value,
		Duration: time.Duration(duration * 1000000000),
	}
}
