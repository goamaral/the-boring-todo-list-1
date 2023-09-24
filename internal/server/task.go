package server

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"example.com/the-boring-to-do-list-1/internal/entity"
	"example.com/the-boring-to-do-list-1/internal/repository"
	gorm_provider "example.com/the-boring-to-do-list-1/pkg/gorm_provider"
)

type taskController struct {
	controller
	TaskRepo repository.AbstractTaskRepository
}

func newTaskController(baseRouter fiber.Router, gormProvider gorm_provider.AbstractProvider) *taskController {
	ctrl := &taskController{
		controller: newController(),
		TaskRepo:   repository.NewTaskRepository(gormProvider),
	}

	router := baseRouter.Group("/tasks")
	router.Post("/", ctrl.CreateTask)
	router.Get("/", ctrl.ListTasks)
	router.Get("/:uuid", ctrl.GetTask)
	router.Patch("/:uuid", ctrl.PatchTask)
	router.Delete("/:uuid", ctrl.DeleteTask)

	return ctrl
}

type NewTask struct {
	Title string `json:"title" validate:"required"`
}
type CreateTaskRequest struct {
	Task NewTask `json:"task" validate:"required"`
}

func (tc *taskController) CreateTask(c *fiber.Ctx) error {
	// Parse request
	req := CreateTaskRequest{}
	err := c.BodyParser(&req)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnprocessableEntity, err)
	}

	// Validate request
	err = tc.validate.Struct(req)
	if err != nil {
		return sendValidationErrorsResponse(c, err.(validator.ValidationErrors))
	}

	// Create task
	task := entity.Task{
		Title: req.Task.Title,
	}
	err = tc.TaskRepo.Create(c.Context(), &task)
	if err != nil {
		return err
	}

	return sendCreatedResponse(c, task.UUID)
}

type ListTasksRequest struct {
	PaginationToken string `json:"paginationToken"`
	Done            bool   `json:"done"`
}
type ListTasksResponse struct {
	Tasks []entity.Task `json:"tasks"`
}

func (tc *taskController) ListTasks(c *fiber.Ctx) error {
	// Parse request
	req := ListTasksRequest{}
	err := c.BodyParser(&req)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnprocessableEntity, err)
	}

	// Validate request
	err = tc.validate.Struct(req)
	if err != nil {
		return sendValidationErrorsResponse(c, err.(validator.ValidationErrors))
	}

	// Get last task fetched
	var lastId uint = 0
	if req.PaginationToken != "" {
		task, err := tc.TaskRepo.FindOne(c.Context(), repository.TaskFilter{UUID: gorm_provider.NewQueryFieldFilter(req.PaginationToken)})
		if err != nil {
			return err
		}
		lastId = task.ID
	}

	// List tasks
	tasks, err := tc.TaskRepo.Find(
		c.Context(),
		repository.TaskFilter{IDGt: gorm_provider.NewQueryFieldFilter(lastId)},
		repository.TaskFilter{Done: gorm_provider.NewQueryFieldFilter(req.Done)},
	)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(ListTasksResponse{Tasks: tasks})
}

type GetTaskResponse struct {
	Task entity.Task `json:"task"`
}

func (tc *taskController) GetTask(c *fiber.Ctx) error {
	// Get task
	task, found, err := tc.TaskRepo.First(c.Context(), repository.TaskFilter{UUID: gorm_provider.NewQueryFieldFilter(c.Params("uuid"))})
	if err != nil {
		return err
	}
	if !found {
		return sendDefaultStatusResponse(c, http.StatusNotFound)
	}

	return c.Status(fiber.StatusOK).JSON(GetTaskResponse{Task: task})
}

type PatchTaskRequest struct {
	Patch repository.TaskPatch `json:"patch"`
}

func (tc *taskController) PatchTask(c *fiber.Ctx) error {
	// Parse request
	req := PatchTaskRequest{}
	err := c.BodyParser(&req)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnprocessableEntity, err)
	}

	// Patch task
	err = tc.TaskRepo.Update(c.Context(), req.Patch, repository.TaskFilter{UUID: gorm_provider.NewQueryFieldFilter(c.Params("uuid"))})
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

func (tc *taskController) DeleteTask(c *fiber.Ctx) error {
	// Delete task
	err := tc.TaskRepo.Delete(c.Context(), repository.TaskFilter{UUID: gorm_provider.NewQueryFieldFilter(c.Params("uuid"))})
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}
