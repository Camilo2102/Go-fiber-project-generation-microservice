package utils

import (
	"github.com/bangadam/go-fiber-starter/app/module/project-builder/request"
	"github.com/bangadam/go-fiber-starter/utils/config"
)

type githubUtils struct {
	cfg                *config.Config
	fileUtils          FileUtils
	githubBase         string
	autoCrudRepository string
}

type GithubUtils interface {
	CopyAutoCrudProjectToUserFolder(req request.ProjectInfo)
}

func NewGithubUtils(cfg *config.Config, fileUtils FileUtils) GithubUtils {
	fileUtils.InitializeFolders()
	fileUtils.CloneGithubProjects()

	return &githubUtils{
		cfg:                cfg,
		fileUtils:          fileUtils,
		autoCrudRepository: cfg.Github.AutoCrudUrl,
	}
}

func (g *githubUtils) CopyAutoCrudProjectToUserFolder(req request.ProjectInfo) {
	g.fileUtils.CreateUserFolder(req.UserId)
	g.fileUtils.CreateProjectCopyForUser(req.ProjectType, req.UserId)
	g.fileUtils.CreateModelsInAutoCrudProject(req.Models, req.UserId)
}
