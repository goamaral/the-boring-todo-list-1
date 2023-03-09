package server

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"

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
	tasksRouter.Put("/:id", ctrl.UpdateTask)
	tasksRouter.Patch("/:id", ctrl.PatchTask)

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
	err := parseBody(c, &req)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnprocessableEntity, err)
	}

	// Validate request
	err = tc.validate.Struct(req)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, err)
	}

	// Create task
	task := entity.Task{
		AbstractEntity: gormprovider.AbstractEntity{Id: ulid.Make().String()},
		Title:          req.Task.Title,
	}
	err = tc.taskRepo.Create(c.Context(), &task)
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
	err := parseBody(c, &req)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnprocessableEntity, err)
	}

	// Validate request
	err = tc.validate.Struct(req)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusBadRequest, err)
	}

	// List tasks
	tasks, err := tc.taskRepo.List(c.Context(), gormprovider.PaginationOption{PageId: req.PageId, PageSize: req.PageSize})
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
	task, err := tc.taskRepo.Get(c.Context(), repository.TaskFilter{Id: c.Params("id")})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(GetTaskResponse{Task: task})
}

type UpdateTaskRequest struct {
	Task entity.Task `json:"task"`
}

func (tc taskController) UpdateTask(c *fiber.Ctx) error {
	// Parse request
	req := UpdateTaskRequest{}
	err := parseBody(c, &req)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnprocessableEntity, err)
	}
	req.Task.AbstractEntity.Id = c.Params("id")

	// Update task
	err = tc.taskRepo.Update(c.Context(), &req.Task, repository.TaskFilter{Id: c.Params("id")})
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
	err := parseBody(c, &req)
	if err != nil {
		return sendErrorResponse(c, fiber.StatusUnprocessableEntity, err)
	}
	req.Task.AbstractEntity.Id = c.Params("id")

	// Patch task
	err = tc.taskRepo.Patch(c.Context(), &req.Task, repository.TaskFilter{Id: c.Params("id")})
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}
