package tasks

import (
	"github.com/springernature/halfpipe/linters/linterrors"
	"github.com/springernature/halfpipe/manifest"
)

func LintSequenceTask(seqTask manifest.Sequence, cameFromAParallel bool) (errs []error, warnings []error) {
	if !cameFromAParallel {
		errs = append(errs, linterrors.NewInvalidField("type", "You are only allowed to use a 'sequence' inside a 'parallel'"))
		return errs, warnings
	}

	if len(seqTask.Tasks) == 0 {
		errs = append(errs, linterrors.NewInvalidField("tasks", "You are not allowed to use a empty 'sequence'"))
		return errs, warnings
	}

	if len(seqTask.Tasks) == 1 {
		warnings = append(warnings, linterrors.NewInvalidField("tasks", "It seems unnecessary to have a single task in a sequence"))
		return errs, warnings
	}

	for _, task := range seqTask.Tasks {
		switch task.(type) {
		case manifest.Sequence:
			errs = append(errs, linterrors.NewInvalidField("tasks", "A sequence task cannot contain sequence tasks"))
		case manifest.Parallel:
			errs = append(errs, linterrors.NewInvalidField("tasks", "A sequence task cannot contain parallel tasks"))
		}
	}

	return errs, warnings
}
