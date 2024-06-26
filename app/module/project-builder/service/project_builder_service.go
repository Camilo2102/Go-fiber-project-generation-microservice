package service

import (
	"github.com/bangadam/go-fiber-starter/app/module/project-builder/request"
	"github.com/bangadam/go-fiber-starter/app/module/project-builder/response"
	"github.com/bangadam/go-fiber-starter/app/module/project-builder/utils"
)

type projectBuilderService struct {
	githubUtils utils.GithubUtils
}

type ProjectBuildService interface {
	CreateAutoCrudProject(module request.Module, userId string, msgChan chan response.ProjectCreateInfo) (res *response.ProjectCreateInfo, err error)
}

func NewProjectBuildService(githubUtils utils.GithubUtils) ProjectBuildService {
	return &projectBuilderService{
		githubUtils: githubUtils,
	}
}

func (p *projectBuilderService) CreateAutoCrudProject(module request.Module, userId string, msgChan chan response.ProjectCreateInfo) (res *response.ProjectCreateInfo, err error) {
	p.githubUtils.CopyAutoCrudProjectToUserFolder(module, userId, msgChan)

	return &response.ProjectCreateInfo{
		Status: true,
	}, nil
}
