package test

import (
	"context"
	"testing"

	"example.com/the-boring-to-do-list-1/internal/entity"
	"example.com/the-boring-to-do-list-1/internal/repository"
	"example.com/the-boring-to-do-list-1/pkg/gorm_provider"
	"github.com/stretchr/testify/require"
)

func AddUser(t *testing.T, p gorm_provider.AbstractProvider, e entity.User) entity.User {
	if e.Username == "" {
		e.Username = RandomString()
	}
	if e.EncryptedPassword == nil {
		e.SetEncryptedPassword("password")
	}
	require.NoError(t, repository.NewUserRepository(p).Create(context.Background(), &e))
	return e
}

func AddTask(t *testing.T, p gorm_provider.AbstractProvider, e entity.Task) entity.Task {
	require.NoError(t, repository.NewTaskRepository(p).Create(context.Background(), &e))
	return e
}
