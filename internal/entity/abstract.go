package entity

import (
	"time"
)

type AbstractEntity struct {
	Id        string `gorm:"primarykey"`
	CreatedAt time.Time
}

func AbstractEntityFromMap(entityMap map[string]interface{}) (AbstractEntity, error) {
	var entity AbstractEntity
	var err error

	entity.Id = entityMap["id"].(string)
	entity.CreatedAt, err = time.Parse(time.RFC3339, entityMap["created_at"].(string))
	if err != nil {
		return AbstractEntity{}, err
	}

	return entity, nil
}
