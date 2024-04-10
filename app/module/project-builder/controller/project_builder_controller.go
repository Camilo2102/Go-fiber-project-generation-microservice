package controller

import (
	"github.com/bangadam/go-fiber-starter/app/module/project-builder/request"
	"github.com/bangadam/go-fiber-starter/app/module/project-builder/service"
	"github.com/bangadam/go-fiber-starter/utils/response"
	"github.com/gofiber/fiber/v2"
)

type projectBuilderController struct {
	projectBuilderService service.ProjectBuildService
}

type ProjectBuilderController interface {
	CreateAutoCrudProject(c *fiber.Ctx) error
}

func NewProjectBuilderController(service service.ProjectBuildService) ProjectBuilderController {
	return &projectBuilderController{
		projectBuilderService: service,
	}
}

func (_i *projectBuilderController) CreateAutoCrudProject(c *fiber.Ctx) error {
	req := new(request.ProjectInfo)
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	res, err := _i.projectBuilderService.CreateAutoCrudProject(*req)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Data:     res,
		Messages: response.Messages{"Created Successfully"},
		Code:     fiber.StatusOK,
	})
}
