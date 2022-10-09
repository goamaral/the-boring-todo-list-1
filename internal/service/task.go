package service

import (
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
	"go.m3o.com/db"

	"example.com/fiber-m3o-validator/internal/entity"
	"example.com/fiber-m3o-validator/pkg/provider"
)

type taskService struct {
	m3oDbClient provider.M3ODb
}

type TaskService interface {
	CreateTask(task entity.Task) (entity.Task, error)
	ListTasks(pageSize uint, opts *ListTasksOpts) ([]entity.Task, error)
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

type ListTasksOpts struct {
	PageId string
}

func (opts *ListTasksOpts) Apply(req *db.ReadRequest) *db.ReadRequest {
	if opts == nil {
		return req
	}

	if opts.PageId != "" {
		req.OrderBy = "id"
		req.Query = fmt.Sprintf("id > '%s'", opts.PageId)
	}

	return req
}

func (ts taskService) ListTasks(pageSize uint, opts *ListTasksOpts) ([]entity.Task, error) {
	var tasks []entity.Task

	// Apply options
	req := opts.Apply(&db.ReadRequest{
		Table: entity.TasksTableName,
		Limit: int32(pageSize),
	})

	// List
	rsp, err := ts.m3oDbClient.Read(req)
	if err != nil {
		return nil, err
	}

	// Convert to entities
	for _, r := range rsp.Records {
		task, err := entity.TaskFromMap(r)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}
