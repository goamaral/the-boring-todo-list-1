package server

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"example.com/the-boring-to-do-list-1/internal/entity"
	"example.com/the-boring-to-do-list-1/internal/repository"
	gormprovider "example.com/the-boring-to-do-list-1/pkg/gormprovider"
)

type taskController struct {
	absctractController
	taskRepo repository.TaskRepository
}

func newTaskController(baseRouter fiber.Router, taskRepo repository.TaskRepository) *taskController {
	ctrl := &taskController{
		absctractController: newAbstractController(),
		taskRepo:            taskRepo,
	}

	router := baseRouter.Group("/tasks")
	router.Post("/", ctrl.CreateTask)
	router.Get("/", ctrl.ListTasks)
	router.Get("/:id", ctrl.GetTask)
	router.Put("/:id", ctrl.UpdateTask)
	router.Patch("/:id", ctrl.PatchTask)
	router.Delete("/:id", ctrl.DeleteTask)

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
		return sendValidationErrorsResponse(c, err.(validator.ValidationErrors))
	}

	// Create task
	task := entity.Task{
		Title: req.Task.Title,
	}
	err = tc.taskRepo.Create(c.Context(), &task)
	if err != nil {
		return err
	}

	return sendCreatedResponse(c, task.Id)
}

type ListTasksRequest struct {
	PageId     string `json:"pageId"`
	PageSize   int    `json:"pageSize" validate:"gte=0,lte=10"`
	IsComplete *bool  `json:"isComplete"`
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
		return sendValidationErrorsResponse(c, err.(validator.ValidationErrors))
	}

	// List tasks
	tasks, err := tc.taskRepo.List(
		c.Context(),
		gormprovider.PaginationOption{PageId: req.PageId, PageSize: req.PageSize},
		repository.TaskFilter{IsComplete: req.IsComplete},
	)
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
	task, found, err := tc.taskRepo.Get(c.Context(), repository.TaskFilter{Id: gormprovider.OptionalValue(c.Params("id"))})
	if err != nil {
		return err
	}
	if !found {
		return sendDefaultStatusResponse(c, http.StatusNotFound)
	}

	return c.Status(fiber.StatusOK).JSON(GetTaskResponse{Task: task})
}

type UpdateTaskRequest struct {
	Task entity.Task `json:"task"`
}

func (tc taskController) UpdateTask(c *fiber.Ctx) error {
	// Parse request
	req := UpdateTaskRequest{}
	err := c.BodyParser(&req)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnprocessableEntity, err)
	}
	req.Task.AbstractEntity.Id = c.Params("id")

	// Update task
	err = tc.taskRepo.Update(c.Context(), &req.Task, repository.TaskFilter{Id: gormprovider.OptionalValue(c.Params("id"))})
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

type PatchTaskRequest struct {
	Task entity.Task `json:"task"`
}

func (tc taskController) PatchTask(c *fiber.Ctx) error {
	// Parse request
	req := PatchTaskRequest{}
	err := c.BodyParser(&req)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnprocessableEntity, err)
	}
	req.Task.AbstractEntity.Id = c.Params("id")

	// Patch task
	err = tc.taskRepo.Patch(c.Context(), &req.Task, repository.TaskFilter{Id: gormprovider.OptionalValue(c.Params("id"))})
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

func (tc taskController) DeleteTask(c *fiber.Ctx) error {
	// Delete task
	err := tc.taskRepo.Delete(c.Context(), repository.TaskFilter{Id: gormprovider.OptionalValue(c.Params("id"))})
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}
