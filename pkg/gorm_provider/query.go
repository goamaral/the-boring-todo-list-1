package gorm_provider

import (
	"errors"

	"gorm.io/gorm"
)

type Query[T AbstractEntity] struct {
	DB *gorm.DB
}

func (q Query[T]) Create(record *T) error {
	return q.DB.Create(record).Error
}

func (q Query[T]) Find() ([]T, error) {
	var records []T
	return records, q.DB.Find(&records).Error
}

// TODO: Error if order option is used
func (q Query[T]) FindInBatches(bacthSize int, fn func([]T) error) error {
	var lastId uint = 0
	for {
		var records []T
		err := q.DB.Where("id > ?", lastId).Limit(bacthSize).Order("id").Find(&records).Error
		if err != nil {
			return err
		}
		if len(records) == 0 {
			break
		}

		err = fn(records)
		if err != nil {
			return err
		}
		lastId = records[len(records)-1].GetID()
	}
	return nil
}

func (q Query[T]) FindOne() (T, error) {
	var record T
	return record, q.DB.First(&record).Error
}

func (q Query[T]) First() (T, bool, error) {
	var record T
	err := q.DB.First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return record, false, nil
		}
		return record, false, err
	}
	return record, true, nil
}

func (q Query[T]) Update(update any) error {
	return q.DB.Updates(update).Error
}

func (q Query[T]) Delete() error {
	var record T
	return q.DB.Delete(&record).Error
}
