package server

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"gorm.io/gorm/clause"

	"example.com/the-boring-to-do-list-1/internal/entity"
	"example.com/the-boring-to-do-list-1/internal/repository"
	gorm_provider "example.com/the-boring-to-do-list-1/pkg/gorm_provider"
	"example.com/the-boring-to-do-list-1/pkg/jwt_provider"
)

type taskController struct {
	controller
	TaskRepo repository.TaskRepository
	UserRepo repository.UserRepository
}

func newTaskController(baseRouter fiber.Router, jwtProvider jwt_provider.Provider, gormProvider gorm_provider.AbstractProvider) *taskController {
	ctrl := &taskController{
		controller: newController(),
		TaskRepo:   repository.NewTaskRepository(gormProvider),
		UserRepo:   repository.NewUserRepository(gormProvider),
	}

	router := baseRouter.Group("/tasks")
	router.Use(NewJWTAuthMiddleware(jwtProvider))
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
type CreateTaskResponse struct {
	UUID string `json:"uuid"`
}

func (tc *taskController) CreateTask(c *fiber.Ctx) error {
	// Parse request
	req := CreateTaskRequest{}
	err := c.BodyParser(&req)
	if err != nil {
		return err
	}

	// Validate request
	err = tc.validate.Struct(req)
	if err != nil {
		return SendValidationErrorsResponse(c, err.(validator.ValidationErrors))
	}

	// Get Auth
	authUser, err := GetAuthUser(c, tc.UserRepo, gorm_provider.SelectOption("id"))
	if err != nil {
		return c.Redirect("/auth/logout")
	}

	// Create task
	task := entity.Task{
		Title:    req.Task.Title,
		AuthorID: authUser.ID,
	}
	err = tc.TaskRepo.Create(c.Context(), &task)
	if err != nil {
		return err
	}

	return c.Status(http.StatusCreated).JSON(CreateTaskResponse{UUID: task.UUID.String()})
}

type ListTasksRequest struct {
	PaginationToken string `query:"paginationToken"`
	Done            string `query:"done"`
}

func (tc *taskController) ListTasks(c *fiber.Ctx) error {
	// Parse request
	req := ListTasksRequest{}
	err := c.QueryParser(&req)
	if err != nil {
		return errors.WithStack(err)
	}

	// Get last task fetched
	var lastId uint = 0
	if req.PaginationToken != "" {
		task, err := tc.TaskRepo.First(c.Context(), clause.Eq{Column: "uuid", Value: req.PaginationToken})
		if err != nil {
			return errors.WithStack(err)
		}
		lastId = task.ID
	}

	// Get Auth
	// TODO: Use join instead of getting user
	authUser, err := GetAuthUser(c, tc.UserRepo, gorm_provider.SelectOption("id"))
	if err != nil {
		return c.Redirect("/auth/logout")
	}

	// List tasks
	opts := []any{
		clause.Eq{Column: "author_id", Value: authUser.ID},
		clause.Gt{Column: "id", Value: lastId},
	}
	if req.Done != "" {
		opts = append(opts, lo.Ternary[clause.Expression](
			req.Done == "true",
			clause.Neq{Column: "done_at", Value: nil},
			clause.Eq{Column: "done_at", Value: nil},
		))
	}
	tasks, err := tc.TaskRepo.Find(c.Context(), opts)
	if err != nil {
		return errors.WithStack(err)
	}

	return c.Render("tasks/index", fiber.Map{"tasks": tasks}, "layouts/private")
}

type GetTaskResponse struct {
	Task entity.Task `json:"task"`
}

func (tc *taskController) GetTask(c *fiber.Ctx) error {
	// Get Auth
	authUser, err := GetAuthUser(c, tc.UserRepo, gorm_provider.SelectOption("id"))
	if err != nil {
		return c.Redirect("/auth/logout")
	}

	// Get task
	task, found, err := tc.TaskRepo.FindOne(
		c.Context(),
		clause.Eq{Column: "author_id", Value: authUser.ID},
		clause.Eq{Column: "uuid", Value: c.Params("uuid")},
	)
	if err != nil {
		return err
	}
	if !found {
		return c.SendStatus(fiber.StatusNotFound)
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
		return err
	}

	// Get Auth
	authUser, err := GetAuthUser(c, tc.UserRepo, gorm_provider.SelectOption("id"))
	if err != nil {
		return c.Redirect("/auth/logout")
	}

	// Patch task
	err = tc.TaskRepo.Update(
		c.Context(),
		req.Patch,
		clause.Eq{Column: "author_id", Value: authUser.ID},
		clause.Eq{Column: "uuid", Value: c.Params("uuid")},
	)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

func (tc *taskController) DeleteTask(c *fiber.Ctx) error {
	// Get Auth
	authUser, err := GetAuthUser(c, tc.UserRepo, gorm_provider.SelectOption("id"))
	if err != nil {
		return c.Redirect("/auth/logout")
	}

	// Delete task
	err = tc.TaskRepo.Delete(
		c.Context(),
		clause.Eq{Column: "author_id", Value: authUser.ID},
		clause.Eq{Column: "uuid", Value: c.Params("uuid")},
	)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}
