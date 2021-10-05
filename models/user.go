package models

import (
	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	TelegramID string
	Address    string
}
