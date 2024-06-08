package test

import "github.com/oklog/ulid/v2"

func RandomString() string {
	return ulid.Make().String()
}
