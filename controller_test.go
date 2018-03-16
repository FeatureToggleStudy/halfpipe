package halfpipe

import (
	"testing"

	"github.com/concourse/atc"
	"github.com/spf13/afero"
	"github.com/springernature/halfpipe/defaults"
	"github.com/springernature/halfpipe/linters"
	"github.com/springernature/halfpipe/linters/errors"
	"github.com/springernature/halfpipe/manifest"
	"github.com/stretchr/testify/assert"
)

func testController() Controller {
	var fs = afero.Afero{Fs: afero.NewMemMapFs()}
	_ = fs.MkdirAll("/pwd/foo/.git", 0777)
	return Controller{
		Fs:         fs,
		CurrentDir: "/pwd/foo",
		Defaulter:  defaults.DefaultValues,
	}
}

func TestProcessDoesNothingWhenFileDoesNotExist(t *testing.T) {
	c := testController()
	pipeline, results := c.Process()

	assert.Empty(t, pipeline)
	assert.Len(t, results, 1)
	assert.IsType(t, errors.FileError{}, results[0].Errors[0])
}

func TestProcessDoesNothingWhenManifestIsEmpty(t *testing.T) {
	c := testController()
	c.Fs.WriteFile("/pwd/foo/.halfpipe.io", []byte(""), 0777)
	pipeline, results := c.Process()

	assert.Empty(t, pipeline)
	assert.Len(t, results, 1)
	assert.IsType(t, errors.FileError{}, results[0].Errors[0])
}

func TestProcessDoesNothingWhenParserFails(t *testing.T) {
	c := testController()
	c.Fs.WriteFile("/pwd/foo/.halfpipe.io", []byte("WrYyYyYy"), 0777)
	pipeline, results := c.Process()

	assert.Empty(t, pipeline)
	assert.Len(t, results, 1)
	assert.IsType(t, manifest.ParseError{}, results[0].Errors[0])
}

type fakeLinter struct {
	Error error
}

func (f fakeLinter) Lint(manifest manifest.Manifest) linters.LintResult {
	return linters.NewLintResult("fake", []error{f.Error})
}

func TestAppliesAllLinters(t *testing.T) {
	c := testController()
	c.Fs.WriteFile("/pwd/foo/.halfpipe.io", []byte("team: asd"), 0777)

	linter1 := fakeLinter{errors.NewFileError("file", "is missing")}
	linter2 := fakeLinter{errors.NewMissingField("field")}
	c.Linters = []linters.Linter{linter1, linter2}

	pipeline, results := c.Process()

	assert.Empty(t, pipeline)
	assert.Len(t, results, 2)
	assert.Equal(t, linter1.Error, results[0].Errors[0])
	assert.Equal(t, linter2.Error, results[1].Errors[0])
}

type FakeRenderer struct {
	Config atc.Config
}

func (f FakeRenderer) Render(manifest manifest.Manifest) atc.Config {
	return f.Config
}

func TestGivesBackAtcConfigWhenLinterPasses(t *testing.T) {
	c := testController()
	c.Fs.WriteFile("/pwd/foo/.halfpipe.io", []byte("team: asd"), 0777)

	config := atc.Config{
		Resources: atc.ResourceConfigs{
			atc.ResourceConfig{
				Name: "Yolo",
			},
		},
	}
	c.Renderer = FakeRenderer{Config: config}

	pipeline, results := c.Process()
	assert.Len(t, results, 0)
	assert.Equal(t, config, pipeline)
}