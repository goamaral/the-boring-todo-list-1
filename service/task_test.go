package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.m3o.com/db"

	"example.com/fiber-m3o-validator/entity"
	"example.com/fiber-m3o-validator/mocks"
	"example.com/fiber-m3o-validator/service"
)

func TestTaskService_CreateTask(t *testing.T) {
	task := entity.Task{Title: "Test Task"}
	expectedId := "task-id"

	mockM3ODbClient := mocks.NewM3ODb(t)
	mockM3ODbClient.On("Create", mock.Anything).Return(&db.CreateResponse{Id: expectedId}, nil)

	s := service.NewTaskService(mockM3ODbClient)
	err := s.CreateTask(&task)
	if assert.NoError(t, err) {
		assert.Equal(t, expectedId, task.Id)
		assert.NotZero(t, task.CreatedAt)
	}
}
