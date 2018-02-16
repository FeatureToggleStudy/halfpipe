package pipeline

import (
	"testing"

	"github.com/concourse/atc"
	"github.com/springernature/halfpipe/model"
	"github.com/stretchr/testify/assert"
)

func TestRendersCfDeployResources(t *testing.T) {
	manifest := model.Manifest{}
	manifest.Tasks = []model.Task{
		model.DeployCF{
			Api:      "dev-api",
			Space:    "space-station",
			Org:      "springer",
			Username: "rob",
			Password: "supersecret",
			Manifest: "manifest-dev.yml",
			Vars: model.Vars{
				"VAR1": "value1",
				"VAR2": "value2",
			},
		},
		model.DeployCF{
			Api:      "live-api",
			Space:    "space-station",
			Org:      "springer",
			Username: "rob",
			Password: "supersecret",
			Manifest: "manifest-live.yml",
		},
	}

	expectedDevResource := atc.ResourceConfig{
		Name: "Cloud Foundry",
		Type: "cf",
		Source: atc.Source{
			"api":          "dev-api",
			"space":        "space-station",
			"organization": "springer",
			"password":     "supersecret",
			"username":     "rob",
		},
	}

	expectedDevJob := atc.JobConfig{
		Name:   "deploy-cf",
		Serial: true,
		Plan: atc.PlanSequence{
			atc.PlanConfig{Get: manifest.Repo.GetName(), Trigger: true},
			atc.PlanConfig{
				Put: "Cloud Foundry",
				Params: atc.Params{
					"manifest": "manifest-dev.yml",
					"environment_variables": map[string]interface{}{
						"VAR1": "value1",
						"VAR2": "value2",
					},
				},
			},
		},
	}

	expectedLiveResource := atc.ResourceConfig{
		Name: "Cloud Foundry (1)",
		Type: "cf",
		Source: atc.Source{
			"api":          "live-api",
			"space":        "space-station",
			"organization": "springer",
			"password":     "supersecret",
			"username":     "rob",
		},
	}

	expectedLiveJob := atc.JobConfig{
		Name:   "deploy-cf (1)",
		Serial: true,
		Plan: atc.PlanSequence{
			atc.PlanConfig{Get: manifest.Repo.GetName(), Trigger: true, Passed: []string{"deploy-cf"}},
			atc.PlanConfig{
				Put: "Cloud Foundry (1)",
				Params: atc.Params{
					"manifest":              "manifest-live.yml",
					"environment_variables": map[string]interface{}{},
				},
			},
		},
	}

	config := testPipeline().Render(manifest)

	assert.Equal(t, expectedDevResource, config.Resources[1])
	assert.Equal(t, expectedDevJob, config.Jobs[0])

	assert.Equal(t, expectedLiveResource, config.Resources[2])
	assert.Equal(t, expectedLiveJob, config.Jobs[1])
}
