package test

import (
	"context"
	"testing"

	"example.com/the-boring-to-do-list-1/internal/entity"
	"example.com/the-boring-to-do-list-1/internal/repository"
	"example.com/the-boring-to-do-list-1/pkg/gorm_provider"
	"github.com/stretchr/testify/require"
)

func AddUser(t *testing.T, p gorm_provider.AbstractProvider, u *entity.User) *entity.User {
	if u == nil {
		u = &entity.User{}
	}
	if u.EncryptedPassword == nil {
		u.SetEncryptedPassword("password")
	}
	userRepo := repository.NewUserRepository(p)
	require.NoError(t, userRepo.Create(context.Background(), u))
	return u
}
