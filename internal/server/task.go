package server

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"example.com/fiber-m3o-validator/internal/entity"
	"example.com/fiber-m3o-validator/internal/service"
	"example.com/fiber-m3o-validator/pkg/errors"
)

type taskController struct {
	taskService service.TaskService
	validate    *validator.Validate
}

func newTaskController(fiberGroup fiber.Router, taskService service.TaskService) *taskController {
	ctrl := &taskController{
		taskService: taskService,
		validate:    validator.New(),
	}

	fiberGroup.Post("/", ctrl.CreateTask)
	fiberGroup.Get("/", ctrl.ListTasks)

	return ctrl
}

/* PUBLIC */
type NewTask struct {
	Title string `json:"title" validate:"required"`
}
type CreateTaskRequest struct {
	Task NewTask `json:"task" validate:"required"`
}

func (tc taskController) CreateTask(c *fiber.Ctx) error {
	// Parse request
	req := CreateTaskRequest{}
	err := c.BodyParser(&req)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnprocessableEntity, errors.Wrap(err, "failed to parse request body"))
	}

	// Validate request
	err = tc.validate.Struct(req)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, errors.Wrap(err, "failed to validate request"))
	}

	// Create task
	task, err := tc.taskService.CreateTask(entity.Task{Title: req.Task.Title})
	if err != nil {
		return errors.Wrap(err, "failed to create task")
	}

	return sendCreateResponse(c, task.Id)
}

type ListTasksRequest struct {
	PageId   string `json:"page_id"`
	PageSize uint   `json:"page_size" validate:"gt=0,lte=10"`
}
type ListTasksResponse struct {
	Tasks []entity.Task `json:"tasks"`
}

func (tc taskController) ListTasks(c *fiber.Ctx) error {
	// Parse request
	req := ListTasksRequest{}
	err := c.BodyParser(&req)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnprocessableEntity, errors.Wrap(err, "failed to parse request body"))
	}

	// Validate request
	err = tc.validate.Struct(req)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, errors.Wrap(err, "failed to validate request"))
	}

	// List tasks
	tasks, err := tc.taskService.ListTasks(req.PageId, req.PageSize)
	if err != nil {
		return errors.Wrap(err, "failed to create task")
	}

	return c.Status(fiber.StatusOK).JSON(ListTasksResponse{Tasks: tasks})
}
