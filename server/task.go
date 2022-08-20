package server

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"

	"example.com/fiber-m3o-validator/entity"
	"example.com/fiber-m3o-validator/provider/thirdparty"
	"example.com/fiber-m3o-validator/service"
)

type taskController struct {
	taskService service.TaskService
	validate    *validator.Validate
}

func newTaskController(fiberGroup fiber.Router, m3oDbClient thirdparty.M3ODb) *taskController {
	ctrl := &taskController{
		taskService: service.NewTaskService(m3oDbClient),
		validate:    validator.New(),
	}

	fiberGroup.Post("/", ctrl.CreateTask)

	return ctrl
}

/* PUBLIC */
type createTaskRequest struct {
	Task struct {
		Title string `json:"title" validate:"required"`
	} `json:"task" validate:"required"`
}

func (tc taskController) CreateTask(c *fiber.Ctx) error {
	// Parse request
	req := createTaskRequest{}
	err := c.BodyParser(&req)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnprocessableEntity, err, "failed to parse request body")
	}

	// Validate task
	err = tc.validate.Struct(req)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, err, "failed to validate task")
	}

	// Create task
	task := entity.Task{Title: req.Task.Title}
	err = tc.taskService.CreateTask(&task)
	if err != nil {
		return errors.Wrap(err, "failed to create task")
	}

	return sendCreatedResponse(c, task.Id)
}
