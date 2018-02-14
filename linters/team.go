package linters

import (
	"github.com/springernature/halfpipe/errors"
	"github.com/springernature/halfpipe/model"
)

type TeamLinter struct{}

func (TeamLinter) Lint(manifest model.Manifest) []error {
	if manifest.Team == "" {
		return []error{
			errors.NewMissingField("team"),
		}
	}
	return nil
}