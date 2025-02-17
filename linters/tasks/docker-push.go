package tasks

import (
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/afero"
	"github.com/springernature/halfpipe/linters/filechecker"
	"github.com/springernature/halfpipe/linters/linterrors"
	"github.com/springernature/halfpipe/manifest"
)

func LintDockerPushTask(docker manifest.DockerPush, fs afero.Afero) (errs []error, warnings []error) {
	if docker.Username == "" {
		errs = append(errs, linterrors.NewMissingField("username"))
	}
	if docker.Password == "" {
		errs = append(errs, linterrors.NewMissingField("password"))
	}
	if docker.Image == "" {
		errs = append(errs, linterrors.NewMissingField("image"))
	} else {
		matched, _ := regexp.Match(`^(.*)/(.*)$`, []byte(docker.Image))
		if !matched {
			errs = append(errs, linterrors.NewInvalidField("image", "must be specified as 'user/image' or 'registry/user/image'"))
		}
	}

	if docker.Retries < 0 || docker.Retries > 5 {
		errs = append(errs, linterrors.NewInvalidField("retries", "must be between 0 and 5"))
	}

	if docker.DockerfilePath == "" {
		errs = append(errs, linterrors.NewInvalidField("dockerfile_path", "must not be empty"))
	}

	if err := filechecker.CheckFile(fs, docker.DockerfilePath, false); err != nil {
		errs = append(errs, err)
	}

	if docker.BuildPath != "" {
		isDir, err := fs.IsDir(docker.BuildPath)
		if err != nil {
			if os.IsNotExist(err) {
				errs = append(errs, linterrors.NewInvalidField("build_path", fmt.Sprintf("'%s' does not exist", docker.BuildPath)))
			} else {
				errs = append(errs, err)
			}
		} else if !isDir {
			errs = append(errs, linterrors.NewInvalidField("build_path", fmt.Sprintf("'%s' must be a directory but is a file ", docker.BuildPath)))
		}
	}

	return errs, warnings
}
