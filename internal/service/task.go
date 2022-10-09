package service

import (
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
	"go.m3o.com/db"

	"example.com/fiber-m3o-validator/internal/entity"
	"example.com/fiber-m3o-validator/pkg/errors"
	"example.com/fiber-m3o-validator/pkg/provider"
)

type taskService struct {
	m3oDbClient provider.M3ODb
}

type TaskService interface {
	CreateTask(task entity.Task) (entity.Task, error)
	ListTasks(pageId string, pageSize uint) ([]entity.Task, error)
}

func NewTaskService(m3oDbClient provider.M3ODb) *taskService {
	return &taskService{
		m3oDbClient: m3oDbClient,
	}
}

/* PUBLIC */
func (ts taskService) CreateTask(task entity.Task) (entity.Task, error) {
	task.CreatedAt = time.Now()
	task.Id = ulid.Make().String()

	rsp, err := ts.m3oDbClient.Create(&db.CreateRequest{
		Id:     task.Id,
		Record: entity.TaskToMap(task),
		Table:  entity.TasksTableName,
	})
	if err != nil {
		return entity.Task{}, err
	}

	task.Id = rsp.Id
	return task, nil
}

func (ts taskService) ListTasks(pageId string, pageSize uint) ([]entity.Task, error) {
	var tasks []entity.Task

	rsp, err := ts.m3oDbClient.Read(&db.ReadRequest{
		Limit:   int32(pageSize),
		OrderBy: "id",
		Query:   fmt.Sprintf("id > '%s'", pageId),
		Table:   entity.TasksTableName,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get tasks from m3o db")
	}

	for _, r := range rsp.Records {
		task, err := entity.TaskFromMap(r)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to convert task (%+v) from map", r)
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}
