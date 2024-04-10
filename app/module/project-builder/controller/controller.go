package controller

import "github.com/bangadam/go-fiber-starter/app/module/project-builder/service"

type Controller struct {
	ProjectBuilder ProjectBuilderController
}

func NewController(projectBuilderService service.ProjectBuildService) *Controller {
	return &Controller{
		ProjectBuilder: NewProjectBuilderController(projectBuilderService),
	}
}
