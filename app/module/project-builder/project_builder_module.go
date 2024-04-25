package project_builder

import (
	"github.com/bangadam/go-fiber-starter/app/module/project-builder/controller"
	"github.com/bangadam/go-fiber-starter/app/module/project-builder/service"
	"github.com/bangadam/go-fiber-starter/app/module/project-builder/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

type ProjectBuilderRouter struct {
	App        fiber.Router
	Controller *controller.Controller
}

func NewProjectBuilderRouter(fiber *fiber.App, controller *controller.Controller) *ProjectBuilderRouter {
	return &ProjectBuilderRouter{
		App:        fiber,
		Controller: controller,
	}
}

var NewProjectBuilderModule = fx.Options(
	fx.Provide(utils.NewFileUtils),
	fx.Provide(utils.NewGithubUtils),
	fx.Provide(service.NewProjectBuildService),
	fx.Provide(controller.NewController),
	fx.Provide(NewProjectBuilderRouter),
)

func (_i *ProjectBuilderRouter) RegisterProjectBuilderRoutes() {
	projectBuilderController := _i.Controller.ProjectBuilder

	_i.App.Route("/api/projectBuilder", func(router fiber.Router) {
		router.Post("/generateCrudProject", projectBuilderController.CreateAutoCrudProject)
	})
}
