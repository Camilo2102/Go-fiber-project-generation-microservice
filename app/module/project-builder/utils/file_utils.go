package utils

import (
	"fmt"
	"github.com/bangadam/go-fiber-starter/app/module/project-builder/request"
	"github.com/bangadam/go-fiber-starter/utils/config"
	"github.com/plus3it/gorecurcopy"
	"github.com/rs/zerolog/log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

type fileUtils struct {
	cfg *config.Config
}

type FileUtils interface {
	InitializeFolders()
	CloneGithubProjects()
	CreateProjectCopyForUser(projectType string, userId string) (error error)
	CreateModelsInAutoCrudProject(models []request.Model, userId string) (error error)
	DockerizeProject(projectType string, userId string) (error error)
	CreateUserFolder(userId string) (error error)
}

func NewFileUtils(cfg *config.Config) FileUtils {
	return &fileUtils{
		cfg: cfg,
	}
}

func (f *fileUtils) InitializeFolders() {
	f.createFolderFromRoot("cloned-projects")
	f.createFolderFromRoot("users")
}

func (f *fileUtils) CloneGithubProjects() {
	f.cloneProject(f.cfg.Github.AutoCrudUrl)
}

func (f *fileUtils) CreateProjectCopyForUser(projectType string, userId string) (error error) {
	originPath := f.getClonedProjectPath(projectType)
	destinyPath := f.getUserProjectPathCopy(projectType, userId)

	if f.existFolder(destinyPath) {
		return nil
	}

	err := gorecurcopy.CopyDirectory(originPath, destinyPath)

	if err != nil {
		log.Warn().Msg("Unable to copy repository")
		return err
	}

	log.Info().Msg("Copy the folder " + originPath + " to destiny " + destinyPath)

	return nil
}

func (f *fileUtils) CreateModelsInAutoCrudProject(models []request.Model, userId string) (error error) {
	userProjectCopy := f.getUserProjectPathCopy(f.cfg.Github.AutoCrudUrl, userId)

	modelPath := filepath.Join(userProjectCopy, f.cfg.Folders.AutoCrudModelFolder, "DemoModel.kt")
	repositoryPath := filepath.Join(userProjectCopy, f.cfg.Folders.AutoCrudRepositoryFolder, "DemoRepository.kt")

	modelTemplate, err1 := os.ReadFile(modelPath)
	repositoryTemplate, err2 := os.ReadFile(repositoryPath)

	if err1 != nil || err2 != nil {
		log.Warn().Msg("Templates not found")
		customErr := fmt.Errorf("failed to find model templates: %v, %v", err1, err2)
		return customErr
	}

	f.createProjectFiles(models, string(modelTemplate), string(repositoryTemplate), userProjectCopy)

	return nil
}

func (f *fileUtils) DockerizeProject(projectType string, userId string) (error error) {
	projectPath := f.getUserProjectPathCopy(projectType, userId)

	err := f.buildDocker(projectType, userId, projectPath)

	if err != nil {
		return err
	}

	err = f.publishDocker(projectType, userId)

	if err != nil {
		return err
	}

	return nil
}

func (f *fileUtils) createProjectFiles(models []request.Model, modelTemplate string, repositoryTemplate string, userCopyPath string) {
	var wg sync.WaitGroup

	wg.Add(len(models) * 2)

	for _, model := range models {
		go func(m request.Model) {
			defer wg.Done()
			f.replaceInDemoModelTemplate(m, modelTemplate, userCopyPath)
		}(model)

		go func(m request.Model) {
			defer wg.Done()
			f.replaceInDemoRepositoryTemplate(m, repositoryTemplate, userCopyPath)
		}(model)
	}
	wg.Wait()
}

func (f *fileUtils) replaceInDemoRepositoryTemplate(model request.Model, template string, userCopyPath string) {
	templateRepositoryName := "Demo"

	repositoryReplaced := strings.Replace(template, templateRepositoryName, model.ModelName, -1)

	newRepositoryPath := filepath.Join(userCopyPath, f.cfg.Folders.AutoCrudRepositoryFolder, model.ModelName+".kt")

	err := os.WriteFile(newRepositoryPath, []byte(repositoryReplaced), 0644)
	if err != nil {
		log.Warn().Msgf("Failed creating repository for model %s: %v", model.ModelName, err)
		return
	}
	log.Info().Msgf("Created file %s", newRepositoryPath)
}

func (f *fileUtils) replaceInDemoModelTemplate(model request.Model, template string, userCopyPath string) {
	templateModelName := "Demo"
	templateModelAttributes := "val replace: String = \"\""

	classNameReplaced := strings.Replace(template, templateModelName, model.ModelName, -1)

	var atrString strings.Builder
	for _, atr := range model.AttributeList {
		atrString.WriteString(fmt.Sprintf("val %s: %s = %s,", atr.Name, atr.Type, atr.DefaultValue))
	}

	attributesAdded := strings.Replace(classNameReplaced, templateModelAttributes, atrString.String(), -1)
	newModelPath := filepath.Join(userCopyPath, f.cfg.Folders.AutoCrudModelFolder, model.ModelName+".kt")

	err := os.WriteFile(newModelPath, []byte(attributesAdded), 0644)
	if err != nil {
		log.Warn().Msgf("Failed creating model for %s: %v", model.ModelName, err)
		return
	}
	log.Info().Msgf("Created file %s", newModelPath)
}

func (f *fileUtils) CreateUserFolder(userId string) (err error) {
	return f.createFolderFromRoot("users", userId)
}

func (f *fileUtils) createFolderFromRoot(folders ...string) (error error) {
	folderPath := filepath.Join(f.cfg.Folders.RootFolder, filepath.Join(folders...))

	if f.existFolder(folderPath) {
		return nil
	}

	err := os.MkdirAll(folderPath, os.ModePerm)

	if err != nil {
		return err
	}

	log.Info().Msg("Folder '" + folderPath + "' created")
	return nil
}

func (f *fileUtils) existFolder(path string) bool {
	_, err := os.Stat(path)

	return !os.IsNotExist(err)
}

func (f *fileUtils) cleanRootFolder() {
	err := os.RemoveAll(f.cfg.Folders.RootFolder)

	if err != nil {
		log.Error().Msg("Failed to empty the root folder")
	}
}

func (f *fileUtils) getClonedProjectPath(projectName string) string {
	return filepath.Join(f.cfg.Folders.RootFolder, "cloned-projects", projectName)
}

func (f *fileUtils) getUserProjectPathCopy(projectName string, userId string) string {
	return filepath.Join(f.cfg.Folders.RootFolder, "users", userId, projectName+"-"+userId)
}

func (f *fileUtils) cloneProject(projectName string) {
	folderPath := f.getClonedProjectPath(projectName)

	if f.existFolder(folderPath) {
		return
	}

	gitRepository := f.cfg.Github.GithubUri + projectName

	cmd := exec.Command("git", "clone", gitRepository, folderPath)

	err := cmd.Run()

	if err != nil {
		log.Error().Msg("Unable to access to git repository " + gitRepository)
		return
	}

	log.Info().Msg("Successfully cloned " + gitRepository)
}

func (f *fileUtils) buildDocker(projectType string, userId string, dockerPath string) (error error) {
	dockerName := f.getDockerName(projectType, userId)

	log.Info().Msg(dockerName)

	cmd := exec.Command("docker", "buildx", "build", "--platform", "linux/amd64", "-t", dockerName, ".")
	cmd.Dir = dockerPath

	err := cmd.Run()

	if err != nil {
		log.Error().Msgf("Unable to dockerize %s", dockerName)
		return err
	}

	log.Info().Msgf("Successfully dockerize %s", dockerName)

	return nil
}

func (f *fileUtils) publishDocker(projectType string, userId string) (error error) {
	dockerName := f.getDockerName(projectType, userId)

	loginCmd := f.loginDocker()
	err := loginCmd.Run()
	if err != nil {
		log.Error().Msg("Failed to log in to Docker")
		return err
	}

	publishCmd := f.publishDockerCmd(dockerName)
	err = publishCmd.Run()
	if err != nil {
		log.Error().Msgf("Failed to publish Docker image %s", dockerName)
		return err
	}

	log.Info().Msgf("Successfully published Docker image %s", dockerName)

	return nil
}

func (f *fileUtils) getDockerName(projectType string, userId string) string {
	return "cammd21/webgen-project-" + strings.ToLower(projectType) + "-" + userId + ":1.0"
}

func (f *fileUtils) loginDocker() *exec.Cmd {
	cmd := exec.Command("docker", "login", "--username", f.cfg.Docker.User, "--password", f.cfg.Docker.Password)
	return cmd
}

func (f *fileUtils) publishDockerCmd(imageName string) *exec.Cmd {
	cmd := exec.Command("docker", "push", imageName)
	return cmd
}
