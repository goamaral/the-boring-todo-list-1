package gorm_provider

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"

	"gorm.io/gorm"
)

type Query[T AbstractEntity] struct {
	DB *gorm.DB
}

func (q Query[T]) Debug() Query[T] {
	return Query[T]{DB: q.DB.Debug()}
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

func (q Query[T]) First() (T, error) {
	var record T
	return record, q.DB.Take(&record).Error
}

func (q Query[T]) FindOne() (T, bool, error) {
	var record T
	err := q.DB.Take(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return record, false, nil
		}
		return record, false, err
	}
	return record, true, nil
}

// Update can be a struct or a map. If struct, it will be converted to a map
// Optional fields should implement AbstractOptionalField interface
func (q Query[T]) Update(update any) error {
	updateMap := map[string]any{}

	kind := reflect.Indirect(reflect.ValueOf(update)).Type().Kind()
	if kind == reflect.Struct {
		var err error
		updateMap, err = structToMap(q.DB.Statement, update)
		if err != nil {
			return fmt.Errorf("failed to convert to map: %w", err)
		}
	} else if kind != reflect.Map {
		return fmt.Errorf("update must be a map")
	}

	var record T
	return q.DB.Model(&record).Updates(updateMap).Error
}

func (q Query[T]) Delete() error {
	var record T
	return q.DB.Delete(&record).Error
}

func structToMap(stmt *gorm.Statement, record any) (map[string]any, error) {
	res := map[string]any{}
	if record == nil {
		return res, nil
	}

	reflectType := reflect.TypeOf(record)
	tableName := stmt.NamingStrategy.TableName(reflectType.Name())
	reflectValue := reflect.Indirect(reflect.ValueOf(record))

	if reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
	}
	for i := 0; i < reflectType.NumField(); i++ {
		fieldName := reflectType.Field(i).Name
		col := stmt.NamingStrategy.ColumnName(tableName, fieldName)
		field := reflectValue.Field(i).Interface()

		if of, ok := field.(AbstractOptionalField); ok && !of.Defined() {
			continue
		}

		if valuer, ok := field.(driver.Valuer); ok {
			val, err := valuer.Value()
			if err != nil {
				return nil, fmt.Errorf("failed to call Value method: %w", err)
			}
			res[col] = val

		} else {
			res[col] = field
		}
	}
	return res, nil
}
