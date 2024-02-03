package repository_test

import (
	"os"
	"testing"

	"example.com/the-boring-to-do-list-1/internal/test"
)

func TestMain(m *testing.M) {
	test.LoadEnv()
	os.Exit(m.Run())
}
