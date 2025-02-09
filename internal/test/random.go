package test

import (
	"github.com/google/uuid"
)

func RandomString() string {
	return uuid.NewString() // TODO: Use std rand
}
