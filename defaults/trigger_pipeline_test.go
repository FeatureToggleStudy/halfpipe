package defaults

import (
	"github.com/springernature/halfpipe/manifest"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPipelineTrigger(t *testing.T) {
	team := "asd"

	t.Run("Empty trigger", func(t *testing.T) {
		expectedTrigger := manifest.PipelineTrigger{
			Team:         team,
			ConcourseURL: DefaultValues.ConcourseURL,
			Username:     DefaultValues.ConcourseUsername,
			Password:     DefaultValues.ConcoursePassword,
			Status:       "succeeded",
		}

		assert.Equal(t, expectedTrigger, defaultPipelineTrigger(manifest.PipelineTrigger{}, DefaultValues, manifest.Manifest{Team: team}))
	})

	t.Run("With already present values for the default values", func(t *testing.T) {
		trigger := manifest.PipelineTrigger{
			ConcourseURL: "url",
			Username:     "username",
			Password:     "password",
			Status:       "asdf",
		}

		expectedTrigger := manifest.PipelineTrigger{
			Team:         team,
			ConcourseURL: "url",
			Username:     "username",
			Password:     "password",
			Status:       "asdf",
		}

		assert.Equal(t, expectedTrigger, defaultPipelineTrigger(trigger, DefaultValues, manifest.Manifest{Team: team}))
	})
}
