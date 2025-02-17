package defaults

import (
	"github.com/springernature/halfpipe/manifest"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testTriggersDefaulter struct {
	apply func(original manifest.TriggerList, defaults Defaults, man manifest.Manifest) (updated manifest.TriggerList)
}

func (t testTriggersDefaulter) Apply(original manifest.TriggerList, defaults Defaults, man manifest.Manifest) (updated manifest.TriggerList) {
	return t.apply(original, defaults, man)
}

type testTasksDefaulter struct {
	apply func(original manifest.TaskList, defaults Defaults, man manifest.Manifest) (updated manifest.TaskList)
}

func (t testTasksDefaulter) Apply(original manifest.TaskList, defaults Defaults, man manifest.Manifest) (updated manifest.TaskList) {
	return t.apply(original, defaults, man)
}

func (t testTasksArtifactoryVarsDefaulter) Apply(original manifest.TaskList, defaults Defaults) (updated manifest.TaskList) {
	return t.apply(original, defaults)
}

func TestUpdatePipeline(t *testing.T) {
	defaults := Defaults{
		triggersDefaulter: testTriggersDefaulter{apply: func(original manifest.TriggerList, defaults Defaults, man manifest.Manifest) (updated manifest.TriggerList) {
			return original
		}},
		tasksDefaulter: testTasksDefaulter{apply: func(original manifest.TaskList, defaults Defaults, man manifest.Manifest) (updated manifest.TaskList) {
			return original
		}},
	}

	t.Run("doesnt do anything when feature toggle is not enabled", func(t *testing.T) {
		assert.Equal(t, manifest.Manifest{}, defaults.Apply(manifest.Manifest{}))
	})

	t.Run("adds update job as first job if feature toggle is enabled", func(t *testing.T) {
		man := manifest.Manifest{
			FeatureToggles: manifest.FeatureToggles{
				manifest.FeatureUpdatePipeline,
			},
		}

		expected := manifest.Manifest{
			FeatureToggles: manifest.FeatureToggles{
				manifest.FeatureUpdatePipeline,
			},
			Tasks: manifest.TaskList{
				manifest.Update{},
			},
		}
		assert.Equal(t, expected, defaults.Apply(man))
	})
}

func TestCallsOutToDefaulters(t *testing.T) {
	expectedTriggers := manifest.TriggerList{
		manifest.TimerTrigger{},
		manifest.GitTrigger{},
	}

	expectedTasks := manifest.TaskList{
		manifest.Run{},
		manifest.DockerPush{},
	}

	defaults := Defaults{
		triggersDefaulter: testTriggersDefaulter{apply: func(original manifest.TriggerList, defaults Defaults, man manifest.Manifest) (updated manifest.TriggerList) {
			return expectedTriggers
		}},
		tasksDefaulter: testTasksDefaulter{apply: func(original manifest.TaskList, defaults Defaults, man manifest.Manifest) (updated manifest.TaskList) {
			return expectedTasks
		}},
	}

	assert.Equal(t, manifest.Manifest{Triggers: expectedTriggers, Tasks: expectedTasks}, defaults.Apply(manifest.Manifest{}))

}
