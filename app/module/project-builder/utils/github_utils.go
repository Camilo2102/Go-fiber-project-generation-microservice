package utils

import (
	"github.com/bangadam/go-fiber-starter/app/module/project-builder/request"
	"github.com/bangadam/go-fiber-starter/app/module/project-builder/response"
	"github.com/bangadam/go-fiber-starter/utils/config"
)

type githubUtils struct {
	cfg                *config.Config
	fileUtils          FileUtils
	githubBase         string
	autoCrudRepository string
}

type GithubUtils interface {
	CopyAutoCrudProjectToUserFolder(module request.Module, userId string, msgChan chan response.ProjectCreateInfo)
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

func (g *githubUtils) CopyAutoCrudProjectToUserFolder(module request.Module, userId string, msgChan chan response.ProjectCreateInfo) {
	fileInitializeStatus := g.fileUtils.CreateUserFolder(userId) == nil

	folderInitializeResponse := response.ProjectCreateInfo{
		Phase:   1,
		Status:  fileInitializeStatus,
		Message: "Folders initialized",
	}

	msgChan <- folderInitializeResponse

	if !fileInitializeStatus {
		return
	}

	projectCopyStatus := g.fileUtils.CreateProjectCopyForUser(module.ModuleName, userId) == nil

	msgChan <- response.ProjectCreateInfo{
		Phase:   2,
		Status:  projectCopyStatus,
		Message: "Project copied",
	}

	if !projectCopyStatus {
		return
	}

	neededFilesCreationStatus := g.fileUtils.CreateModelsInAutoCrudProject(module.Models, userId) == nil

	msgChan <- response.ProjectCreateInfo{
		Phase:   3,
		Status:  neededFilesCreationStatus,
		Message: "Requiered files created",
	}

	if !neededFilesCreationStatus {
		return
	}

	dockerizeStatus := g.fileUtils.DockerizeProject(module.ModuleName, userId) == nil

	msgChan <- response.ProjectCreateInfo{
		Phase:   3,
		Status:  dockerizeStatus,
		Message: "Project dockerized and published",
	}

}
