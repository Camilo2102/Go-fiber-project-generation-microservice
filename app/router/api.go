package router

import (
	project_builder "github.com/bangadam/go-fiber-starter/app/module/project-builder"
	"github.com/bangadam/go-fiber-starter/utils/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

type Router struct {
	App fiber.Router
	Cfg *config.Config

	ProjectBuilderRouter *project_builder.ProjectBuilderRouter
}

func NewRouter(
	fiber *fiber.App,
	cfg *config.Config,

	projectBuilderRouter *project_builder.ProjectBuilderRouter,
) *Router {
	return &Router{
		App:                  fiber,
		Cfg:                  cfg,
		ProjectBuilderRouter: projectBuilderRouter,
	}
}

// Register routes
func (r *Router) Register() {
	// Test Routes
	r.App.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("Pong! 👋")
	})

	// Swagger Documentation
	r.App.Get("/swagger/*", swagger.HandlerDefault)

	// Register routes of modules
	r.ProjectBuilderRouter.RegisterProjectBuilderRoutes()
}
