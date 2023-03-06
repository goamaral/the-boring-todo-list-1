package server

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"example.com/the-boring-to-do-list-1/internal/entity"
	"example.com/the-boring-to-do-list-1/internal/repository"
	gormprovider "example.com/the-boring-to-do-list-1/pkg/provider/gorm"
)

type taskController struct {
	taskRepo repository.TaskRepository
	validate *validator.Validate
}

func newTaskController(baseRouter fiber.Router, taskRepo repository.TaskRepository) *taskController {
	ctrl := &taskController{
		taskRepo: taskRepo,
		validate: validator.New(),
	}

	tasksRouter := baseRouter.Group("/tasks")
	tasksRouter.Post("/", ctrl.CreateTask)
	tasksRouter.Get("/", ctrl.ListTasks)
	tasksRouter.Get("/:id", ctrl.GetTask)

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
	task := entity.Task{Title: req.Task.Title}
	err = tc.taskRepo.CreateTask(c.Context(), &task)
	if err != nil {
		return err
	}

	return sendCreatedResponse(c, task.Id)
}

type ListTasksRequest struct {
	PageId   string `json:"pageId"`
	PageSize int    `json:"pageSize" validate:"gte=0,lte=10"`
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

	// List tasks
	tasks, err := tc.taskRepo.ListTasks(c.Context(), gormprovider.PaginationOption{PageId: req.PageId, PageSize: req.PageSize})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(ListTasksResponse{Tasks: tasks})
}

type GetTaskResponse struct {
	Task entity.Task `json:"task"`
}

func (tc taskController) GetTask(c *fiber.Ctx) error {
	// Get task
	task, err := tc.taskRepo.GetTask(c.Context(), repository.TaskFilter{Id: c.Params("id")})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(GetTaskResponse{Task: task})
}
