package defaults

import (
	"fmt"
	"strings"

	"github.com/springernature/halfpipe/config"
	"github.com/springernature/halfpipe/manifest"
	"github.com/springernature/halfpipe/project"
)

type Defaulter func(manifest.Manifest, project.Data) manifest.Manifest

type Defaults struct {
	RepoPrivateKey       string
	CfUsername           string
	CfPassword           string
	CfUsernameSnPaas     string
	CfPasswordSnPaas     string
	CfOrgSnPaas          string
	CfAPISnPaas          string
	CfManifest           string
	CfTestDomains        map[string]string
	DockerUsername       string
	DockerPassword       string
	Project              project.Data
	DockerComposeService string
	ArtifactoryUsername  string
	ArtifactoryPassword  string
	ArtifactoryURL       string
	ConcourseURL         string
	ConcourseUsername    string
	ConcoursePassword    string
	Timeout              string
}

func NewDefaulter(project project.Data) Defaults {
	d := DefaultValues
	d.Project = project
	return d
}

func (d Defaults) getUniqueName(name string, previousNames []string, counter int) string {
	candidate := name
	if counter > 0 {
		candidate = fmt.Sprintf("%s (%v)", name, counter)
	}

	for _, previousName := range previousNames {
		if previousName == candidate {
			return d.getUniqueName(name, previousNames, counter+1)
		}
	}

	return candidate

}

func (d Defaults) uniqueName(name string, previousNames []string) string {
	return d.getUniqueName(name, previousNames, 0)
}

func (d Defaults) uniqueifyNames(tasks manifest.TaskList) manifest.TaskList {
	var previousNames []string
	var taskSwitcher func(tasks manifest.TaskList) manifest.TaskList
	taskSwitcher = func(tasks manifest.TaskList) manifest.TaskList {
		var updatedTasks manifest.TaskList
		for _, task := range tasks {
			switch task := task.(type) {
			case manifest.DeployCF:
				task.Name = d.uniqueName(task.GetName(), previousNames)
				task.PrePromote = d.uniqueifyNames(task.PrePromote)
				previousNames = append(previousNames, task.GetName())
				updatedTasks = append(updatedTasks, task)

			case manifest.Run:
				task.Name = d.uniqueName(task.GetName(), previousNames)
				previousNames = append(previousNames, task.GetName())
				updatedTasks = append(updatedTasks, task)

			case manifest.DockerPush:
				task.Name = d.uniqueName(task.GetName(), previousNames)
				previousNames = append(previousNames, task.GetName())
				updatedTasks = append(updatedTasks, task)

			case manifest.DockerCompose:
				task.Name = d.uniqueName(task.GetName(), previousNames)
				previousNames = append(previousNames, task.GetName())
				updatedTasks = append(updatedTasks, task)

			case manifest.ConsumerIntegrationTest:
				task.Name = d.uniqueName(task.GetName(), previousNames)
				previousNames = append(previousNames, task.GetName())
				updatedTasks = append(updatedTasks, task)

			case manifest.DeployMLModules:
				task.Name = d.uniqueName(task.GetName(), previousNames)
				previousNames = append(previousNames, task.GetName())
				updatedTasks = append(updatedTasks, task)

			case manifest.DeployMLZip:
				task.Name = d.uniqueName(task.GetName(), previousNames)
				previousNames = append(previousNames, task.GetName())
				updatedTasks = append(updatedTasks, task)

			case manifest.Update:
				previousNames = append(previousNames, task.GetName())
				updatedTasks = append(updatedTasks, task)

			case manifest.Parallel:
				task.Tasks = taskSwitcher(task.Tasks)
				updatedTasks = append(updatedTasks, task)
			case manifest.Sequence:
				task.Tasks = taskSwitcher(task.Tasks)
				updatedTasks = append(updatedTasks, task)
			default:
				panic("got a unknown task..")
			}
		}
		return updatedTasks
	}
	return taskSwitcher(tasks)
}

func (d Defaults) updateTasks(tasks manifest.TaskList, man manifest.Manifest) manifest.TaskList {
	var taskSwitcher func(tasks manifest.TaskList) manifest.TaskList
	taskSwitcher = func(tasks manifest.TaskList) manifest.TaskList {

		var updatedTasks manifest.TaskList

		for _, task := range tasks {
			switch task := task.(type) {
			case manifest.DeployCF:
				if task.API == d.CfAPISnPaas {
					if task.Org == "" {
						task.Org = d.CfOrgSnPaas
					}
					if task.Username == "" {
						task.Username = d.CfUsernameSnPaas
					}
					if task.Password == "" {
						task.Password = d.CfPasswordSnPaas
					}
				} else {
					if task.Org == "" {
						task.Org = man.Team
					}
					if task.Username == "" {
						task.Username = d.CfUsername
					}
					if task.Password == "" {
						task.Password = d.CfPassword
					}
				}

				if task.Manifest == "" {
					task.Manifest = d.CfManifest
				}
				if task.PrePromote != nil {
					task.PrePromote = d.updateTasks(task.PrePromote, man)
				}
				if task.TestDomain == "" {
					if domain, ok := d.CfTestDomains[task.API]; ok {
						task.TestDomain = domain
					}
				}

				if task.GetTimeout() == "" {
					task.Timeout = d.Timeout
				}

				updatedTasks = append(updatedTasks, task)

			case manifest.Run:
				if strings.HasPrefix(task.Docker.Image, config.DockerRegistry) {
					task.Docker.Username = d.DockerUsername
					task.Docker.Password = d.DockerPassword
				}
				task.Vars = d.addArtifactoryCredentialsToVars(task.Vars)

				if task.GetTimeout() == "" {
					task.Timeout = d.Timeout
				}

				updatedTasks = append(updatedTasks, task)

			case manifest.DockerPush:
				if strings.HasPrefix(task.Image, config.DockerRegistry) {
					task.Username = d.DockerUsername
					task.Password = d.DockerPassword
				}

				if task.DockerfilePath == "" {
					task.DockerfilePath = "Dockerfile"
				}

				task.Vars = d.addArtifactoryCredentialsToVars(task.Vars)

				if task.GetTimeout() == "" {
					task.Timeout = d.Timeout
				}

				updatedTasks = append(updatedTasks, task)

			case manifest.DockerCompose:
				if task.Service == "" {
					task.Service = d.DockerComposeService
				}

				task.Vars = d.addArtifactoryCredentialsToVars(task.Vars)

				if task.GetTimeout() == "" {
					task.Timeout = d.Timeout
				}

				updatedTasks = append(updatedTasks, task)

			case manifest.ConsumerIntegrationTest:
				task.Vars = d.addArtifactoryCredentialsToVars(task.Vars)

				if task.GetTimeout() == "" {
					task.Timeout = d.Timeout
				}

				updatedTasks = append(updatedTasks, task)

			case manifest.DeployMLModules:
				if task.GetTimeout() == "" {
					task.Timeout = d.Timeout
				}

				updatedTasks = append(updatedTasks, task)

			case manifest.DeployMLZip:
				if task.GetTimeout() == "" {
					task.Timeout = d.Timeout
				}

				updatedTasks = append(updatedTasks, task)

			case manifest.Update:
				task.Timeout = d.Timeout
				updatedTasks = append(updatedTasks, task)

			case manifest.Parallel:
				task.Tasks = taskSwitcher(task.Tasks)
				updatedTasks = append(updatedTasks, task)
			case manifest.Sequence:
				task.Tasks = taskSwitcher(task.Tasks)
				updatedTasks = append(updatedTasks, task)
			default:
				panic("got a unknown task..")
			}
		}
		return updatedTasks
	}
	return taskSwitcher(tasks)
}

func (d Defaults) updateGitTriggerWithDefaults(man manifest.Manifest) manifest.Manifest {
	// We assume that the first git trigger we find is the right one as we lint later that we only have trigger.

	updateGitTrigger := func(t manifest.GitTrigger) manifest.Trigger {
		t.BasePath = d.Project.BasePath

		if t.URI == "" {
			t.URI = d.Project.GitURI

			for from, to := range config.RewriteGitHTTPToSSH {
				if strings.Contains(t.URI, from) {
					t.URI = strings.Replace(t.URI, from, to, 1)
				}
			}
		}

		if t.URI != "" && !t.IsPublic() && t.PrivateKey == "" {
			t.PrivateKey = d.RepoPrivateKey
		}
		return t
	}

	var updatedTriggers manifest.TriggerList
	var gitTriggerFound bool

	for _, trigger := range man.Triggers {
		switch trigger := trigger.(type) {
		case manifest.GitTrigger:
			gitTriggerFound = true
			updatedTriggers = append(updatedTriggers, updateGitTrigger(trigger))
		default:
			updatedTriggers = append(updatedTriggers, trigger)
		}
	}

	if !gitTriggerFound {
		man.Triggers = updatedTriggers
		updatedTriggers = append(updatedTriggers, updateGitTrigger(manifest.GitTrigger{}))
	}

	man.Triggers = updatedTriggers

	return man
}
func (d Defaults) updatePipelineTriggerWithDefaults(man manifest.Manifest) manifest.Manifest {
	var updatedTriggers manifest.TriggerList
	for _, trigger := range man.Triggers {
		switch trigger := trigger.(type) {
		case manifest.PipelineTrigger:
			if trigger.Team == "" {
				trigger.Team = man.Team
			}

			if trigger.ConcourseURL == "" {
				trigger.ConcourseURL = d.ConcourseURL
			}

			if trigger.Username == "" {
				trigger.Username = d.ConcourseUsername
			}

			if trigger.Password == "" {
				trigger.Password = d.ConcoursePassword
			}

			if trigger.Status == "" {
				trigger.Status = "succeeded"
			}

			updatedTriggers = append(updatedTriggers, trigger)
		default:
			updatedTriggers = append(updatedTriggers, trigger)
		}
	}

	man.Triggers = updatedTriggers
	return man
}

func (d Defaults) updateDockerTriggerWithDefaults(man manifest.Manifest) manifest.Manifest {
	// We assume that the first docker trigger we find is the right one as we lint later that we only have trigger.

	var updatedTriggers manifest.TriggerList
	for _, trigger := range man.Triggers {
		switch trigger := trigger.(type) {
		case manifest.DockerTrigger:
			if strings.HasPrefix(trigger.Image, config.DockerRegistry) {
				trigger.Username = d.DockerUsername
				trigger.Password = d.DockerPassword
			}
			updatedTriggers = append(updatedTriggers, trigger)
		default:
			updatedTriggers = append(updatedTriggers, trigger)
		}
	}

	man.Triggers = updatedTriggers
	return man
}

func (d Defaults) updateTriggersWithDefaults(man manifest.Manifest) manifest.Manifest {
	man = d.updateGitTriggerWithDefaults(man)
	man = d.updateDockerTriggerWithDefaults(man)
	man = d.updatePipelineTriggerWithDefaults(man)
	return man
}

func (d Defaults) Update(man manifest.Manifest) manifest.Manifest {
	updated := d.updateTriggersWithDefaults(man)

	if updated.FeatureToggles.UpdatePipeline() {
		updated.Tasks = append(manifest.TaskList{manifest.Update{}}, updated.Tasks...)
	}

	updated.Tasks = d.uniqueifyNames(updated.Tasks)
	updated.Tasks = d.updateTasks(updated.Tasks, updated)

	return updated

}

func (d Defaults) addArtifactoryCredentialsToVars(vars manifest.Vars) manifest.Vars {
	updatedVars := map[string]string{}
	for key, value := range vars {
		updatedVars[key] = value
	}

	if _, found := updatedVars["ARTIFACTORY_USERNAME"]; !found {
		updatedVars["ARTIFACTORY_USERNAME"] = d.ArtifactoryUsername
	}

	if _, found := updatedVars["ARTIFACTORY_PASSWORD"]; !found {
		updatedVars["ARTIFACTORY_PASSWORD"] = d.ArtifactoryPassword
	}

	if _, found := updatedVars["ARTIFACTORY_URL"]; !found {
		updatedVars["ARTIFACTORY_URL"] = d.ArtifactoryURL
	}

	return updatedVars
}
