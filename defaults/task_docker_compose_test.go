package defaults

import (
	"github.com/springernature/halfpipe/manifest"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetsDefaultDockerComposeService(t *testing.T) {
	assert.Equal(t, DefaultValues.DockerComposeService, dockerComposeDefaulter(manifest.DockerCompose{}, DefaultValues).Service)
}

func TestDoesntOverrideService(t *testing.T) {
	service := "asdf"
	assert.Equal(t, service, dockerComposeDefaulter(manifest.DockerCompose{Service: service}, DefaultValues).Service)
}

func TestSetsDefaultDockerComposeFile(t *testing.T) {
	assert.Equal(t, DefaultValues.DockerComposeFile, dockerComposeDefaulter(manifest.DockerCompose{}, DefaultValues).ComposeFile)
}

func TestDoesntOverrideComposeFile(t *testing.T) {
	file := "docker-compose-asdf.yml"
	assert.Equal(t, file, dockerComposeDefaulter(manifest.DockerCompose{ComposeFile: file}, DefaultValues).ComposeFile)
}
