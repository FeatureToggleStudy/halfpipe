package linters

import (
	"fmt"

	"github.com/springernature/halfpipe/linters/errors"
	"github.com/springernature/halfpipe/parser"
)

type artifactsLinter struct {
}

func NewArtifactsLinter() artifactsLinter {
	return artifactsLinter{}
}

func (linter artifactsLinter) Lint(man parser.Manifest) (result LintResult) {
	result.Linter = "Artifacts"

	var artifacts int
	var artifact string
	for _, t := range man.Tasks {
		switch task := t.(type) {
		case parser.Run:
			if len(task.SaveArtifacts) > 0 {
				artifacts += 1
				artifact = task.SaveArtifacts[0]
				if artifacts > 1 {
					result.AddError(errors.NewInvalidField("run.save_artifact", "Found multiple 'save_artifact', currently halfpipe only supports saving artifacts from on task"))
					return
				}
				if len(task.SaveArtifacts) > 1 {
					result.AddError(errors.NewInvalidField("run.save_artifact", "Found multiple artifacts in 'save_artifact', currently halfpipe only supports saving one artifacts"))
					return
				}
			}

		case parser.DeployCF:
			if task.DeployArtifact != "" && task.DeployArtifact != artifact {
				var errorStr string
				if artifact == "" {
					errorStr = fmt.Sprintf("No previous tasks have saved the artifact '%s'", task.DeployArtifact)
				} else {
					errorStr = fmt.Sprintf("No previous tasks have saved the artifact '%s', but I found a previous job that saves the artifact '%s'.", task.DeployArtifact, artifact)
				}
				result.AddError(errors.NewInvalidField("deploy-cf.deploy_artifact", errorStr))
			}
		}
	}

	return
}
