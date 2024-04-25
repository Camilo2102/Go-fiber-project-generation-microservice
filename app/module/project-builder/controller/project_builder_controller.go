package controller

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/bangadam/go-fiber-starter/app/module/project-builder/request"
	response2 "github.com/bangadam/go-fiber-starter/app/module/project-builder/response"
	"github.com/bangadam/go-fiber-starter/app/module/project-builder/service"
	"github.com/bangadam/go-fiber-starter/utils/response"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
	"sync"
	"time"
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

	msgChan := make(chan response2.ProjectCreateInfo)

	functionMap := map[string]func(module *request.Module) (*response2.ProjectCreateInfo, error){
		"Springboot-kotlin-autogenerable-crud": func(module *request.Module) (*response2.ProjectCreateInfo, error) {
			return _i.projectBuilderService.CreateAutoCrudProject(*module, req.UserId, msgChan)
		},
	}

	var wg sync.WaitGroup
	wg.Add(len(req.Modules))

	go func() {
		defer close(msgChan)

		for _, module := range req.Modules {
			module := module
			go func(moduleName string) {
				defer wg.Done()
				if fn, ok := functionMap[moduleName]; ok {
					_, err := fn(&module)
					if err != nil {
						return
					}
				} else {
					log.Info().Msg("errr")
				}
			}(module.ModuleName)
		}

		wg.Wait()

		msgChan <- response2.ProjectCreateInfo{
			Status:  true,
			Phase:   4,
			Message: "Project generated successfully",
		}
	}()

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		for {
			partialResponse, _ := <-msgChan

			if partialResponse.Status == false {
				return
			}

			jsonResponse, _ := json.Marshal(partialResponse)

			fmt.Fprintf(w, "data: %s\n\n", jsonResponse)

			w.Flush()

			time.Sleep(1 * time.Second)
		}
	}))

	return nil
}
