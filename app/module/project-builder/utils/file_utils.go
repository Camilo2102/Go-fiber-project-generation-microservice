package utils

import (
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
	CreateProjectCopyForUser(projectName string, userId string)
	CreateModelsInAutoCrudProject(models []request.ModelInfo, userId string)
	CreateUserFolder(userId string)
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

func (f *fileUtils) CreateProjectCopyForUser(projectName string, userId string) {
	originPath := f.getClonedProjectPath(projectName)
	destinyPath := f.getUserProjectPathToCopy(projectName, userId)

	if f.existFolder(destinyPath) {
		return
	}

	err := gorecurcopy.CopyDirectory(originPath, destinyPath)

	if err != nil {
		log.Warn().Msg("Unable to copy repository")
		return
	}

	log.Info().Msg("Copy the folder " + originPath + " to destiny " + destinyPath)
}

func (f *fileUtils) CreateModelsInAutoCrudProject(models []request.ModelInfo, userId string) {
	userProjectCopy := f.getUserProjectPathToCopy(f.cfg.Github.AutoCrudUrl, userId)

	modelPath := filepath.Join(userProjectCopy, f.cfg.Folders.AutoCrudModelFolder, "DemoModel.kt")
	repositoryPath := filepath.Join(userProjectCopy, f.cfg.Folders.AutoCrudRepositoryFolder, "DemoRepository.kt")

	modelTemplate, err1 := os.ReadFile(modelPath)
	repositoryTemplate, err2 := os.ReadFile(repositoryPath)

	if err1 != nil || err2 != nil {
		log.Warn().Msg("Model template not found")
		return
	}

	f.replaceInDemoModelTemplate(models, string(modelTemplate), userProjectCopy)
	f.replaceInDemoRepositoryTemplate(models, string(repositoryTemplate), userProjectCopy)
}

func (f *fileUtils) replaceInfoThreaded(models []request.ModelInfo, modelTemplate string, repositoryTemplate string, userCopyPath string) {
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		f.replaceInDemoModelTemplate(models, modelTemplate, userCopyPath)
		defer wg.Done()
	}()

	go func() {
		f.replaceInDemoRepositoryTemplate(models, repositoryTemplate, userCopyPath)
		defer wg.Done()
	}()

	wg.Wait()
}

func (f *fileUtils) replaceInDemoRepositoryTemplate(models []request.ModelInfo, template string, userCopyPath string) {
	var wg sync.WaitGroup
	templateRepositoryName := "Demo"

	for _, model := range models {
		wg.Add(1)

		go func(m request.ModelInfo) {
			repositoryReplaced := strings.Replace(template, templateRepositoryName, m.ModelName, -1)

			newRepositoryPath := filepath.Join(userCopyPath, f.cfg.Folders.AutoCrudRepositoryFolder, m.ModelName+".kt")

			err := os.WriteFile(newRepositoryPath, []byte(repositoryReplaced), 0644)

			if err != nil {
				log.Warn().Msg("Failed creating repository")
				wg.Done()
				return
			}
			wg.Done()
			log.Info().Msg("Created file " + newRepositoryPath)
		}(model)
	}

}

func (f *fileUtils) replaceInDemoModelTemplate(models []request.ModelInfo, template string, userCopyPath string) {
	var wg sync.WaitGroup
	templateModelName := "Demo"
	templateModelAttributes := "val replace: String = \"\""

	for _, model := range models {
		wg.Add(1)

		go func(m request.ModelInfo) {
			classNameReplaced := strings.Replace(template, templateModelName, m.ModelName, -1)

			atrString := ""
			for _, atr := range m.AttributeList {
				atrString += "val " + atr.Name + ": " + atr.Type + " = " + atr.DefaultValue + ","
			}

			attributesAdded := strings.Replace(classNameReplaced, templateModelAttributes, atrString, -1)
			newModelPath := filepath.Join(userCopyPath, f.cfg.Folders.AutoCrudModelFolder, m.ModelName+".kt")

			err := os.WriteFile(newModelPath, []byte(attributesAdded), 0644)
			if err != nil {
				log.Warn().Msg("Failed creating model")
				wg.Done()
				return
			}

			wg.Done()
			log.Info().Msg("Created file " + newModelPath)
		}(model)
	}

}

func (f *fileUtils) CreateUserFolder(userId string) {
	f.createFolderFromRoot("users", userId)
}

func (f *fileUtils) createFolderFromRoot(folders ...string) {
	folderPath := filepath.Join(f.cfg.Folders.RootFolder, filepath.Join(folders...))

	if f.existFolder(folderPath) {
		return
	}

	err := os.MkdirAll(folderPath, os.ModePerm)

	if err != nil {
		log.Warn().Err(err)
		return
	}

	log.Info().Msg("Folder '" + folderPath + "' created")
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

func (f *fileUtils) getUserProjectPathToCopy(projectName string, userId string) string {
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
	}

	log.Info().Msg("Successfully cloned " + gitRepository)
}
