package gormprovider

import "time"

type AbstractEntity struct {
	Id        string    `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
