package service

import (
	"time"

	"github.com/oklog/ulid/v2"
	"go.m3o.com/db"

	"example.com/fiber-m3o-validator/entity"
	"example.com/fiber-m3o-validator/provider/thirdparty"
)

type taskService struct {
	m3oDbClient thirdparty.M3ODb
}

type TaskService interface {
	CreateTask(task *entity.Task) error
}

func NewTaskService(m3oDbClient thirdparty.M3ODb) *taskService {
	return &taskService{
		m3oDbClient: m3oDbClient,
	}
}

/* PUBLIC */
func (ts taskService) CreateTask(task *entity.Task) error {
	task.CreatedAt = time.Now()
	task.Id = ulid.Make().String()

	rsp, err := ts.m3oDbClient.Create(&db.CreateRequest{
		Record: task.ToMap(),
		Table:  task.GetTableName(),
	})
	if err != nil {
		return err
	}

	task.Id = rsp.Id
	return nil
}
