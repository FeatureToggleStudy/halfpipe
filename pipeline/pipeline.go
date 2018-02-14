package pipeline

import (
	"fmt"
	"strings"

	"github.com/concourse/atc"
	"github.com/ghodss/yaml"
	"github.com/springernature/halfpipe/model"
)

type Config atc.Config

type Renderer interface {
	Render(manifest model.Manifest) Config
}

type Pipeline struct{}

func (Pipeline) gitResource(repo model.Repo) atc.ResourceConfig {
	sources := atc.Source{
		"uri": repo.Uri,
	}

	if repo.PrivateKey != "" {
		sources["private_key"] = repo.PrivateKey
	}

	return atc.ResourceConfig{
		Name:   repo.GetName(),
		Type:   "git",
		Source: sources,
	}
}

func (p Pipeline) makeImageResource(image string) *atc.ImageResource {
	repo, tag := image, "latest"
	if strings.Contains(image, ":") {
		split := strings.Split(image, ":")
		repo = split[0]
		tag = split[1]
	}
	return &atc.ImageResource{
		Type: "docker-image",
		Source: atc.Source{
			"repository": repo,
			"tag":        tag,
		},
	}
}

func (p Pipeline) makeRunJob(task model.Run, repo model.Repo) atc.JobConfig {
	return atc.JobConfig{
		Name:   task.Script,
		Serial: true,
		Plan: atc.PlanSequence{
			atc.PlanConfig{Get: repo.GetName(), Trigger: true},
			atc.PlanConfig{
				Task: task.Script,
				TaskConfig: &atc.TaskConfig{
					Platform:      "linux",
					Params:        task.Vars,
					ImageResource: p.makeImageResource(task.Image),
					Run: atc.TaskRunConfig{
						Path: "sh",
						Dir:  repo.GetName(),
						Args: []string{"-exc", fmt.Sprintf("%s", task.Script)},
					},
					Inputs: []atc.TaskInputConfig{
						{Name: repo.GetName()},
					},
				}}}}
}

func (p Pipeline) Render(manifest model.Manifest) Config {
	config := Config{
		Resources: atc.ResourceConfigs{
			p.gitResource(manifest.Repo),
		},
	}

	for _, t := range manifest.Tasks {
		switch task := t.(type) {
		case model.Run:
			config.Jobs = append(config.Jobs, p.makeRunJob(task, manifest.Repo))
		}
	}
	return config
}

func (c Config) ToString() (string, error) {
	renderedPipeline, err := yaml.Marshal(c)
	return string(renderedPipeline), err
}
