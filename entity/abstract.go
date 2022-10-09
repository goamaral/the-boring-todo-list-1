package entity

import (
	"time"
)

type AbstractEntity struct {
	Id        string
	CreatedAt time.Time
}
