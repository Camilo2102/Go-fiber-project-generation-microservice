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
	CreateAutoCrudProject(req request.ProjectInfo) (res *response.ProjectCreateInfo, err error)
}

func NewProjectBuildService(githubUtils utils.GithubUtils) ProjectBuildService {
	return &projectBuilderService{
		githubUtils: githubUtils,
	}
}

func (p *projectBuilderService) CreateAutoCrudProject(req request.ProjectInfo) (res *response.ProjectCreateInfo, err error) {
	p.githubUtils.CopyAutoCrudProjectToUserFolder(req)

	return &response.ProjectCreateInfo{
		Status: true,
	}, nil
}
