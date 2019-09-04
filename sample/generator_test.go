package sample

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/springernature/halfpipe/project"
	"github.com/stretchr/testify/assert"
)

type fakeProjectResolver struct {
	p   project.Data
	err error
}

func (pr fakeProjectResolver) ShouldFindManifestPath() project.Project {
	panic("implement me")
}

func (pr fakeProjectResolver) ShouldParseManifest() project.Project {
	panic("implement me")
}

func (pr fakeProjectResolver) Parse(workingDir string) (p project.Data, err error) {
	return pr.p, pr.err
}

func TestFailsIfHalfpipeFileAlreadyExists(t *testing.T) {

	t.Run(".halfpipe.io", func(t *testing.T) {
		fs := afero.Afero{Fs: afero.NewMemMapFs()}
		fs.WriteFile(".halfpipe.io", []byte(""), 777)

		sampleGenerator := NewSampleGenerator(fs, fakeProjectResolver{p: project.Data{HalfpipeFilePath: ".halfpipe.io"}}, "/home/user/src/myApp")

		err := sampleGenerator.Generate()

		assert.Equal(t, err, ErrHalfpipeAlreadyExists)
	})

	t.Run(".halfpipe.io.yaml", func(t *testing.T) {
		fs := afero.Afero{Fs: afero.NewMemMapFs()}
		fs.WriteFile(".halfpipe.io.yaml", []byte(""), 777)

		sampleGenerator := NewSampleGenerator(fs, fakeProjectResolver{p: project.Data{HalfpipeFilePath: ".halfpipe.io.yaml"}}, "/home/user/src/myApp")

		err := sampleGenerator.Generate()

		assert.Equal(t, err, ErrHalfpipeAlreadyExists)
	})
}

func TestFailsIfProjectResolverErrorsOut(t *testing.T) {
	expectedError := errors.New("Oeh noes")

	fs := afero.Afero{Fs: afero.NewMemMapFs()}

	sampleGenerator := NewSampleGenerator(fs, fakeProjectResolver{err: expectedError}, "/home/user/src/myApp")

	err := sampleGenerator.Generate()

	assert.Equal(t, err, expectedError)
}

func TestWritesSample(t *testing.T) {
	fs := afero.Afero{Fs: afero.NewMemMapFs()}

	sampleGenerator := NewSampleGenerator(fs, fakeProjectResolver{p: project.Data{RootName: "myApp"}}, "/home/user/src/myApp")

	err := sampleGenerator.Generate()

	assert.Nil(t, err)

	bytes, err := fs.ReadFile(".halfpipe.io")
	assert.Nil(t, err)

	expected := `team: CHANGE-ME
pipeline: myApp
feature_toggles:
- update-pipeline
tasks:
- type: run
  name: CHANGE-ME OPTIONAL NAME IN CONCOURSE UI
  script: ./gradlew CHANGE-ME
  docker:
    image: CHANGE-ME:tag
`
	assert.Equal(t, string(bytes), expected)
}

func TestWritesSampleWhenExecutedInASubDirectory(t *testing.T) {
	fs := afero.Afero{Fs: afero.NewMemMapFs()}

	sampleGenerator := NewSampleGenerator(fs, fakeProjectResolver{p: project.Data{
		BasePath: "subApp",
		RootName: "myApp",
		GitURI:   "",
	}}, "/home/user/src/myApp/subApp")

	err := sampleGenerator.Generate()

	assert.Nil(t, err)

	bytes, err := fs.ReadFile(".halfpipe.io")
	assert.Nil(t, err)

	expected := `team: CHANGE-ME
pipeline: myApp-subApp
feature_toggles:
- update-pipeline
triggers:
- type: git
  watched_paths:
  - subApp
tasks:
- type: run
  name: CHANGE-ME OPTIONAL NAME IN CONCOURSE UI
  script: ./gradlew CHANGE-ME
  docker:
    image: CHANGE-ME:tag
`
	assert.Equal(t, expected, string(bytes))
}

func TestWritesSampleWhenExecutedInASubSubDirectory(t *testing.T) {
	fs := afero.Afero{Fs: afero.NewMemMapFs()}

	sampleGenerator := NewSampleGenerator(fs, fakeProjectResolver{p: project.Data{
		BasePath: "subFolder/subApp",
		RootName: "myApp",
		GitURI:   "",
	}}, "/home/user/src/myApp/subFolder/subApp")

	err := sampleGenerator.Generate()

	assert.Nil(t, err)

	bytes, err := fs.ReadFile(".halfpipe.io")
	assert.Nil(t, err)

	expected := `team: CHANGE-ME
pipeline: myApp-subFolder-subApp
feature_toggles:
- update-pipeline
triggers:
- type: git
  watched_paths:
  - subFolder/subApp
tasks:
- type: run
  name: CHANGE-ME OPTIONAL NAME IN CONCOURSE UI
  script: ./gradlew CHANGE-ME
  docker:
    image: CHANGE-ME:tag
`
	assert.Equal(t, string(bytes), expected)
}
