package server

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"example.com/fiber-m3o-validator/internal/entity"
	"example.com/fiber-m3o-validator/internal/service"
)

type taskController struct {
	taskService service.TaskService
	validate    *validator.Validate
}

func newTaskController(baseRouter fiber.Router, taskService service.TaskService) *taskController {
	ctrl := &taskController{
		taskService: taskService,
		validate:    validator.New(),
	}

	tasksRouter := baseRouter.Group("/tasks")
	tasksRouter.Post("/", ctrl.CreateTask)
	tasksRouter.Get("/", ctrl.ListTasks)

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
		return sendErrorResponse(c, fiber.StatusUnprocessableEntity, err)
	}

	// Validate request
	err = tc.validate.Struct(req)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, err)
	}

	// Create task
	task, err := tc.taskService.CreateTask(entity.Task{Title: req.Task.Title})
	if err != nil {
		return err
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
		return sendErrorResponse(c, fiber.StatusUnprocessableEntity, err)
	}

	// Validate request
	err = tc.validate.Struct(req)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, err)
	}

	// Build options
	opts := service.ListTasksOpts{PageId: req.PageId}

	// List tasks
	tasks, err := tc.taskService.ListTasks(req.PageSize, &opts)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(ListTasksResponse{Tasks: tasks})
}
