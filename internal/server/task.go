package server

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	html "github.com/gofiber/template/html/v2"
	"github.com/google/uuid"
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
	viewEngine *html.Engine
	TaskRepo   repository.TaskRepository
	UserRepo   repository.UserRepository
}

func newTaskController(baseRouter fiber.Router, jwtProvider jwt_provider.Provider, gormProvider gorm_provider.AbstractProvider, viewEngine *html.Engine) *taskController {
	ctrl := &taskController{
		controller: newController(),
		viewEngine: viewEngine,
		TaskRepo:   repository.NewTaskRepository(gormProvider),
		UserRepo:   repository.NewUserRepository(gormProvider),
	}

	router := baseRouter.Group("/tasks")
	router.Use(NewJWTAuthMiddleware(jwtProvider))
	router.Get("/new", ctrl.NewTask)
	router.Post("/", ctrl.CreateTask)
	router.Get("/", ctrl.ListTasks)
	router.Get("/:uuid", ctrl.GetTask)
	router.Patch("/:uuid", ctrl.PatchTask)
	router.Delete("/:uuid", ctrl.DeleteTask)

	return ctrl
}

type TaskForm struct {
	UUID  uuid.UUID
	Title string
}

func (f TaskForm) Method() string {
	return lo.Ternary(f.UUID == uuid.Nil, fiber.MethodPost, fiber.MethodPatch)
}

func (f TaskForm) Action() string {
	return lo.Ternary(f.UUID == uuid.Nil, "/tasks", "/tasks/"+f.UUID.String())
}

func (tc *taskController) NewTask(c *fiber.Ctx) error {
	return c.Render("tasks/new", fiber.Map{"form": TaskForm{}}, "layouts/private")
}

type CreateTaskRequest struct {
	Title string `json:"title" validate:"required"`
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
		// TODO: Return view with validation errors
		return SendValidationErrorsResponse(c, err.(validator.ValidationErrors))
	}

	// Get Auth
	authUser, err := GetAuthUser(c, tc.UserRepo, gorm_provider.SelectOption("id"))
	if err != nil {
		return c.Redirect("/auth/logout")
	}

	// Create task
	task := entity.Task{
		Title:    req.Title,
		AuthorID: authUser.ID,
	}
	err = tc.TaskRepo.Create(c.Context(), &task)
	if err != nil {
		return err
	}

	return c.Redirect("/tasks/" + task.UUID.String())
}

type ListTasksRequest struct {
	Page uint   `query:"page"`
	Done string `query:"done"`
}

func (tc *taskController) ListTasks(c *fiber.Ctx) error {
	// Parse request
	req := ListTasksRequest{}
	err := c.QueryParser(&req)
	if err != nil {
		return errors.WithStack(err)
	}

	// Get Auth
	// TODO: Use join instead of getting user
	authUser, err := GetAuthUser(c, tc.UserRepo, gorm_provider.SelectOption("id"))
	if err != nil {
		return c.Redirect("/auth/logout")
	}

	// List tasks
	pageSize := 2
	opts := []any{
		clause.Eq{Column: "author_id", Value: authUser.ID},
		clause.Limit{Offset: int(req.Page) * pageSize, Limit: &pageSize},
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

	return c.Render(
		"tasks/index",
		fiber.Map{
			"tasks":    tasks,
			"prevPage": lo.Ternary(req.Page > 0, fmt.Sprint(req.Page-1), ""),
			"nextPage": lo.Ternary(len(tasks) == pageSize, fmt.Sprint(req.Page+1), ""),
		},
		"layouts/private",
	)
}

func (tc *taskController) GetTask(c *fiber.Ctx) error {
	// Get Auth
	authUser, err := GetAuthUser(c, tc.UserRepo, gorm_provider.SelectOption("id"))
	if err != nil {
		return Logout(c)
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

	return c.Render(
		"tasks/show",
		fiber.Map{"task": task},
		"layouts/private",
	)
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
