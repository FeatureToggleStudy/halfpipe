package linters

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/springernature/halfpipe/manifest"
	"github.com/stretchr/testify/assert"
)

func testTaskLinter() taskLinter {
	return taskLinter{
		Fs: afero.Afero{Fs: afero.NewMemMapFs()},
	}
}

func TestAtLeastOneTaskExists(t *testing.T) {
	man := manifest.Manifest{}
	taskLinter := testTaskLinter()

	result := taskLinter.Lint(man)
	assert.Len(t, result.Errors, 1)
	assertMissingField(t, "tasks", result.Errors[0])
}

func TestRunTaskWithoutScriptAndImage(t *testing.T) {
	man := manifest.Manifest{}
	taskLinter := testTaskLinter()

	man.Tasks = []manifest.Task{
		manifest.Run{},
	}

	result := taskLinter.Lint(man)
	assert.Len(t, result.Errors, 2)
	assertMissingField(t, "run.script", result.Errors[0])
	assertMissingField(t, "run.docker.image", result.Errors[1])
}

func TestRunTaskWithScriptAndImage(t *testing.T) {
	taskLinter := testTaskLinter()
	man := manifest.Manifest{}
	man.Tasks = []manifest.Task{
		manifest.Run{
			Script: "./build.sh",
			Docker: manifest.Docker{
				Image: "alpine",
			},
		},
	}

	result := taskLinter.Lint(man)
	assert.Len(t, result.Errors, 1)
	assertFileError(t, "./build.sh", result.Errors[0])
}

func TestRunTaskWithScriptAndImageWithPasswordAndUsername(t *testing.T) {
	taskLinter := testTaskLinter()
	taskLinter.Fs.WriteFile("build.sh", []byte("foo"), 0777)
	man := manifest.Manifest{}
	man.Tasks = []manifest.Task{
		manifest.Run{
			Script: "./build.sh",
			Docker: manifest.Docker{
				Image:    "alpine",
				Password: "secret",
				Username: "Michiel",
			},
		},
	}

	result := taskLinter.Lint(man)
	assert.Len(t, result.Errors, 0)
}

func TestRunTaskWithScriptAndImageAndOnlyPassword(t *testing.T) {
	taskLinter := testTaskLinter()
	taskLinter.Fs.WriteFile("build.sh", []byte("foo"), 0777)
	man := manifest.Manifest{}
	man.Tasks = []manifest.Task{
		manifest.Run{
			Script: "./build.sh",
			Docker: manifest.Docker{
				Image:    "alpine",
				Password: "secret",
			},
		},
	}

	result := taskLinter.Lint(man)
	assert.Len(t, result.Errors, 1)
	assertMissingField(t, "run.docker.username", result.Errors[0])
}
func TestRunTaskWithScriptAndImageAndOnlyUsername(t *testing.T) {
	taskLinter := testTaskLinter()
	taskLinter.Fs.WriteFile("build.sh", []byte("foo"), 0777)
	man := manifest.Manifest{}
	man.Tasks = []manifest.Task{
		manifest.Run{
			Script: "./build.sh",
			Docker: manifest.Docker{
				Image:    "alpine",
				Username: "Michiel",
			},
		},
	}

	result := taskLinter.Lint(man)
	assert.Len(t, result.Errors, 1)
	assertMissingField(t, "run.docker.password", result.Errors[0])
}

func TestRunTaskScriptFileExists(t *testing.T) {
	taskLinter := testTaskLinter()
	taskLinter.Fs.WriteFile("build.sh", []byte("foo"), 0777)

	man := manifest.Manifest{}
	man.Tasks = []manifest.Task{
		manifest.Run{
			Script: "./build.sh",
			Docker: manifest.Docker{
				Image: "alpine",
			},
		},
	}

	result := taskLinter.Lint(man)
	assert.Len(t, result.Errors, 0)
}

func TestCFDeployTaskWithEmptyTask(t *testing.T) {
	taskLinter := testTaskLinter()
	man := manifest.Manifest{}
	man.Tasks = []manifest.Task{
		manifest.DeployCF{Manifest: "manifest.yml"},
	}

	result := taskLinter.Lint(man)
	assert.Len(t, result.Errors, 4)
	assertMissingField(t, "deploy-cf.api", result.Errors[0])
	assertMissingField(t, "deploy-cf.space", result.Errors[1])
	assertMissingField(t, "deploy-cf.org", result.Errors[2])
	assertFileError(t, "manifest.yml", result.Errors[3])
}

func TestDockerPushTaskWithEmptyTask(t *testing.T) {
	taskLinter := testTaskLinter()
	man := manifest.Manifest{
		Tasks: []manifest.Task{
			manifest.DockerPush{},
		},
	}

	result := taskLinter.Lint(man)
	assert.Len(t, result.Errors, 4)
	assertMissingField(t, "docker-push.username", result.Errors[0])
	assertMissingField(t, "docker-push.password", result.Errors[1])
	assertMissingField(t, "docker-push.image", result.Errors[2])
	assertFileError(t, "Dockerfile", result.Errors[3])

}

func TestDockerPushTaskWithBadRepo(t *testing.T) {
	taskLinter := testTaskLinter()
	man := manifest.Manifest{
		Tasks: []manifest.Task{
			manifest.DockerPush{
				Username: "asd",
				Password: "asd",
				Image:    "asd",
			},
		},
	}

	result := taskLinter.Lint(man)
	assert.Len(t, result.Errors, 2)
	assertInvalidField(t, "docker-push.image", result.Errors[0])
	assertFileError(t, "Dockerfile", result.Errors[1])

}

func TestDockerPushTaskWhenDockerfileIsMissing(t *testing.T) {
	taskLinter := testTaskLinter()
	man := manifest.Manifest{
		Tasks: []manifest.Task{
			manifest.DockerPush{
				Username: "asd",
				Password: "asd",
				Image:    "asd/asd",
			},
		},
	}

	result := taskLinter.Lint(man)
	assert.Len(t, result.Errors, 1)
	assertFileError(t, "Dockerfile", result.Errors[0])
}

func TestDockerPushTaskWithCorrectData(t *testing.T) {
	taskLinter := testTaskLinter()
	taskLinter.Fs.WriteFile("Dockerfile", []byte("FROM ubuntu"), 0777)

	man := manifest.Manifest{
		Tasks: []manifest.Task{
			manifest.DockerPush{
				Username: "asd",
				Password: "asd",
				Image:    "asd/asd",
			},
		},
	}

	result := taskLinter.Lint(man)
	assert.Len(t, result.Errors, 0)
}

func TestEnvVarsMustBeUpperCase(t *testing.T) {
	taskLinter := testTaskLinter()

	badKey1 := "KeHe"
	badKey2 := "b"
	badKey3 := "AAAAa"

	goodKey1 := "YO"
	goodKey2 := "A"
	goodKey3 := "AOIJASOID"

	man := manifest.Manifest{
		Tasks: []manifest.Task{
			manifest.Run{
				Vars: map[string]string{
					badKey1:  "a",
					goodKey1: "sup",
				},
			},

			manifest.DockerPush{
				Vars: map[string]string{
					goodKey2: "a",
					badKey2:  "B",
				},
			},

			manifest.DeployCF{
				Vars: map[string]string{
					badKey3:  "asd",
					goodKey3: "asd",
				},
			},
		},
	}

	result := taskLinter.Lint(man)
	assertInvalidFieldInErrors(t, badKey1, result.Errors)
	assertInvalidFieldInErrors(t, badKey2, result.Errors)
	assertInvalidFieldInErrors(t, badKey3, result.Errors)

	assertInvalidFieldShouldNotBeInErrors(t, goodKey1, result.Errors)
	assertInvalidFieldShouldNotBeInErrors(t, goodKey2, result.Errors)
	assertInvalidFieldShouldNotBeInErrors(t, goodKey3, result.Errors)
}