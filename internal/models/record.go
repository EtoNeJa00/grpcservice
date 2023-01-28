package models

import (
	"github.com/google/uuid"
)

type Record struct {
	Id   uuid.UUID
	Data string
}
