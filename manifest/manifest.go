package manifest

import (
	"fmt"
	"regexp"
	"strings"
)

type TaskList []Task

func (taskList TaskList) NotifiesOnSuccess() bool {
	for _, task := range taskList {
		if task.NotifiesOnSuccess() {
			return true
		}
	}
	return false
}

type Task interface {
	ReadsFromArtifacts() bool
	GetAttempts() int
	SavesArtifacts() bool
	SavesArtifactsOnFailure() bool
	IsManualTrigger() bool
	NotifiesOnSuccess() bool
	GetTimeout() string
}

type Manifest struct {
	Team            string
	Pipeline        string
	SlackChannel    string         `json:"slack_channel,omitempty" yaml:"slack_channel,omitempty"`
	TriggerInterval string         `json:"trigger_interval" yaml:"trigger_interval,omitempty"`
	CronTrigger     string         `json:"cron_trigger" yaml:"cron_trigger,omitempty"`
	Repo            Repo           `yaml:"repo,omitempty"`
	ArtifactConfig  ArtifactConfig `json:"artifact_config,omitempty" yaml:"artifact_config,omitempty"`
	FeatureToggles  FeatureToggles `json:"feature_toggles,omitempty" yaml:"feature_toggles,omitempty"`
	Tasks           TaskList
}

func (m Manifest) NotifiesOnFailure() bool {
	return m.SlackChannel != ""
}

func (m Manifest) PipelineName() (pipelineName string) {
	re := regexp.MustCompile(`[^A-Za-z0-9\-]`)
	sanitize := func(s string) string {
		return re.ReplaceAllString(strings.TrimSpace(s), "_")
	}

	pipelineName = m.Pipeline
	if m.Repo.Branch != "" && m.Repo.Branch != "master" {
		pipelineName = fmt.Sprintf("%s-%s", sanitize(m.Pipeline), sanitize(m.Repo.Branch))
	}
	return
}

type ArtifactConfig struct {
	Bucket  string `json:"bucket" yaml:"bucket,omitempty" secretAllowed:"true"`
	JSONKey string `json:"json_key" yaml:"json_key,omitempty" secretAllowed:"true"`
}

type Repo struct {
	URI          string   `json:"uri,omitempty" yaml:"uri,omitempty"`
	BasePath     string   `json:"-" yaml:"-"` //don't auto unmarshal
	PrivateKey   string   `json:"private_key,omitempty" yaml:"private_key,omitempty" secretAllowed:"true"`
	WatchedPaths []string `json:"watched_paths,omitempty" yaml:"watched_paths,omitempty"`
	IgnoredPaths []string `json:"ignored_paths,omitempty" yaml:"ignored_paths,omitempty"`
	GitCryptKey  string   `json:"git_crypt_key,omitempty" yaml:"git_crypt_key,omitempty" secretAllowed:"true"`
	Branch       string   `json:"branch,omitempty" yaml:"branch,omitempty"`
	Shallow      bool     `json:"shallow,omitempty" yaml:"shallow,omitempty"`
}

func (repo Repo) IsPublic() bool {
	return len(repo.URI) > 4 && repo.URI[:4] == "http"
}

type Docker struct {
	Image    string
	Username string `yaml:"username,omitempty" secretAllowed:"true"`
	Password string `yaml:"password,omitempty" secretAllowed:"true"`
}

type Run struct {
	Type                   string
	Name                   string
	ManualTrigger          bool `json:"manual_trigger" yaml:"manual_trigger,omitempty"`
	Script                 string
	Docker                 Docker
	Vars                   Vars     `yaml:"vars,omitempty" secretAllowed:"true"`
	SaveArtifacts          []string `json:"save_artifacts" yaml:"save_artifacts,omitempty"`
	RestoreArtifacts       bool     `json:"restore_artifacts" yaml:"restore_artifacts,omitempty"`
	SaveArtifactsOnFailure []string `json:"save_artifacts_on_failure" yaml:"save_artifacts_on_failure,omitempty"`
	Parallel               bool     `yaml:"parallel,omitempty"`
	Retries                int      `yaml:"retries,omitempty"`
	NotifyOnSuccess        bool     `json:"notify_on_success,omitempty" yaml:"notify_on_success,omitempty"`
	Timeout                string   `json:"timeout,omitempty" yaml:"timeout,omitempty"`
}

func (r Run) GetTimeout() string {
	return r.Timeout
}

func (r Run) NotifiesOnSuccess() bool {
	return r.NotifyOnSuccess
}

func (r Run) SavesArtifactsOnFailure() bool {
	return len(r.SaveArtifactsOnFailure) > 0
}

func (r Run) IsManualTrigger() bool {
	return r.ManualTrigger
}

func (r Run) SavesArtifacts() bool {
	return len(r.SaveArtifacts) > 0
}

func (r Run) ReadsFromArtifacts() bool {
	return r.RestoreArtifacts
}

func (r Run) GetAttempts() int {
	return 1 + r.Retries
}

type DockerPush struct {
	Type             string
	Name             string
	ManualTrigger    bool   `json:"manual_trigger" yaml:"manual_trigger"`
	Username         string `secretAllowed:"true"`
	Password         string `secretAllowed:"true"`
	Image            string
	Vars             Vars   `secretAllowed:"true"`
	RestoreArtifacts bool   `json:"restore_artifacts" yaml:"restore_artifacts"`
	Parallel         bool   `yaml:"parallel,omitempty"`
	Retries          int    `yaml:"retries,omitempty"`
	NotifyOnSuccess  bool   `json:"notify_on_success,omitempty" yaml:"notify_on_success,omitempty"`
	Timeout          string `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	DockerfilePath   string `json:"dockerfile_path,omitempty" yaml:"dockerfile_path,omitempty"`
	BuildPath        string `json:"build_path,omitempty" yaml:"build_path,omitempty"`
}

func (r DockerPush) GetTimeout() string {
	return r.Timeout
}

func (r DockerPush) NotifiesOnSuccess() bool {
	return r.NotifyOnSuccess
}

func (r DockerPush) SavesArtifactsOnFailure() bool {
	return false
}

func (r DockerPush) IsManualTrigger() bool {
	return r.ManualTrigger
}

func (r DockerPush) SavesArtifacts() bool {
	return false
}

func (r DockerPush) ReadsFromArtifacts() bool {
	return r.RestoreArtifacts
}

func (r DockerPush) GetAttempts() int {
	return 1 + r.Retries
}

type DockerCompose struct {
	Type                   string
	Name                   string
	Command                string
	ManualTrigger          bool `json:"manual_trigger" yaml:"manual_trigger"`
	Vars                   Vars `secretAllowed:"true"`
	Service                string
	ComposeFile            string   `json:"compose_file"`
	SaveArtifacts          []string `json:"save_artifacts"`
	RestoreArtifacts       bool     `json:"restore_artifacts" yaml:"restore_artifacts"`
	SaveArtifactsOnFailure []string `json:"save_artifacts_on_failure" yaml:"save_artifacts_on_failure,omitempty"`
	Parallel               bool     `yaml:"parallel,omitempty"`
	Retries                int      `yaml:"retries,omitempty"`
	NotifyOnSuccess        bool     `json:"notify_on_success,omitempty" yaml:"notify_on_success,omitempty"`
	Timeout                string   `json:"timeout,omitempty" yaml:"timeout,omitempty"`
}

func (r DockerCompose) GetTimeout() string {
	return r.Timeout
}

func (r DockerCompose) NotifiesOnSuccess() bool {
	return r.NotifyOnSuccess
}

func (r DockerCompose) SavesArtifactsOnFailure() bool {
	return len(r.SaveArtifactsOnFailure) > 0
}

func (r DockerCompose) IsManualTrigger() bool {
	return r.ManualTrigger
}

func (r DockerCompose) SavesArtifacts() bool {
	return len(r.SaveArtifacts) > 0
}

func (r DockerCompose) ReadsFromArtifacts() bool {
	return r.RestoreArtifacts
}

func (r DockerCompose) GetAttempts() int {
	return 1 + r.Retries
}

type DeployCF struct {
	Type            string
	Name            string
	ManualTrigger   bool   `json:"manual_trigger" yaml:"manual_trigger"`
	API             string `secretAllowed:"true"`
	Space           string `secretAllowed:"true"`
	Org             string `secretAllowed:"true"`
	Username        string `secretAllowed:"true"`
	Password        string `secretAllowed:"true"`
	Manifest        string
	TestDomain      string   `json:"test_domain" yaml:"test_domain" secretAllowed:"true"`
	Vars            Vars     `secretAllowed:"true"`
	DeployArtifact  string   `json:"deploy_artifact"`
	PrePromote      TaskList `json:"pre_promote"`
	Parallel        bool     `yaml:"parallel,omitempty"`
	Timeout         string   `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Retries         int      `yaml:"retries,omitempty"`
	NotifyOnSuccess bool     `json:"notify_on_success,omitempty" yaml:"notify_on_success,omitempty"`
}

func (r DeployCF) GetTimeout() string {
	return r.Timeout
}

func (r DeployCF) NotifiesOnSuccess() bool {
	return r.NotifyOnSuccess || r.PrePromote.NotifiesOnSuccess()
}

func (r DeployCF) SavesArtifactsOnFailure() bool {
	for _, task := range r.PrePromote {
		if task.SavesArtifactsOnFailure() {
			return true
		}
	}
	return false
}

func (r DeployCF) IsManualTrigger() bool {
	return r.ManualTrigger
}

func (r DeployCF) SavesArtifacts() bool {
	return false
}

func (r DeployCF) ReadsFromArtifacts() bool {
	if r.DeployArtifact != "" || strings.HasPrefix(r.Manifest, "../artifacts/") {
		return true
	}

	for _, pp := range r.PrePromote {
		if pp.ReadsFromArtifacts() {
			return true
		}
	}
	return false
}

func (r DeployCF) GetAttempts() int {
	return 2 + r.Retries
}

type ConsumerIntegrationTest struct {
	Type                 string
	Name                 string
	Consumer             string
	ConsumerHost         string `json:"consumer_host" yaml:"consumer_host"`
	GitCloneOptions      string `json:"git_clone_options,omitempty" yaml:"git_clone_options,omitempty"`
	ProviderHost         string `json:"provider_host" yaml:"provider_host"`
	Script               string
	DockerComposeService string `json:"docker_compose_service" yaml:"docker_compose_service"`
	Parallel             bool   `yaml:"parallel,omitempty"`
	Vars                 Vars   `secretAllowed:"true"`
	Retries              int    `yaml:"retries,omitempty"`
	NotifyOnSuccess      bool   `json:"notify_on_success,omitempty" yaml:"notify_on_success,omitempty"`
	Timeout              string `json:"timeout,omitempty" yaml:"timeout,omitempty"`
}

func (r ConsumerIntegrationTest) GetTimeout() string {
	return r.Timeout
}

func (r ConsumerIntegrationTest) NotifiesOnSuccess() bool {
	return r.NotifyOnSuccess
}

func (r ConsumerIntegrationTest) SavesArtifactsOnFailure() bool {
	return false
}

func (r ConsumerIntegrationTest) IsManualTrigger() bool {
	return false
}

func (r ConsumerIntegrationTest) SavesArtifacts() bool {
	return false
}

func (r ConsumerIntegrationTest) ReadsFromArtifacts() bool {
	return false
}

func (r ConsumerIntegrationTest) GetAttempts() int {
	return 1 + r.Retries
}

type DeployMLZip struct {
	Type            string
	Name            string
	Parallel        bool     `yaml:"parallel,omitempty"`
	DeployZip       string   `json:"deploy_zip"`
	AppName         string   `json:"app_name"`
	AppVersion      string   `json:"app_version"`
	Targets         []string `secretAllowed:"true"`
	ManualTrigger   bool     `json:"manual_trigger" yaml:"manual_trigger"`
	Retries         int      `yaml:"retries,omitempty"`
	NotifyOnSuccess bool     `json:"notify_on_success,omitempty" yaml:"notify_on_success,omitempty"`
	Timeout         string   `json:"timeout,omitempty" yaml:"timeout,omitempty"`
}

func (r DeployMLZip) GetTimeout() string {
	return r.Timeout
}

func (r DeployMLZip) NotifiesOnSuccess() bool {
	return r.NotifyOnSuccess
}

func (r DeployMLZip) SavesArtifactsOnFailure() bool {
	return false
}

func (r DeployMLZip) IsManualTrigger() bool {
	return r.ManualTrigger
}

func (r DeployMLZip) SavesArtifacts() bool {
	return false
}

func (r DeployMLZip) GetAttempts() int {
	return 1 + r.Retries
}

func (r DeployMLZip) ReadsFromArtifacts() bool {
	return true
}

type DeployMLModules struct {
	Type             string
	Name             string
	Parallel         bool     `yaml:"parallel,omitempty"`
	MLModulesVersion string   `json:"ml_modules_version"`
	AppName          string   `json:"app_name"`
	AppVersion       string   `json:"app_version"`
	Targets          []string `secretAllowed:"true"`
	ManualTrigger    bool     `json:"manual_trigger" yaml:"manual_trigger"`
	Retries          int      `yaml:"retries,omitempty"`
	NotifyOnSuccess  bool     `json:"notify_on_success,omitempty" yaml:"notify_on_success,omitempty"`
	Timeout          string   `json:"timeout,omitempty" yaml:"timeout,omitempty"`
}

func (r DeployMLModules) GetTimeout() string {
	return r.Timeout
}

func (r DeployMLModules) NotifiesOnSuccess() bool {
	return r.NotifyOnSuccess
}

func (r DeployMLModules) SavesArtifactsOnFailure() bool {
	return false
}

func (r DeployMLModules) IsManualTrigger() bool {
	return r.ManualTrigger
}

func (r DeployMLModules) SavesArtifacts() bool {
	return false
}

func (r DeployMLModules) ReadsFromArtifacts() bool {
	return false
}

func (r DeployMLModules) GetAttempts() int {
	return 1 + r.Retries
}

type Vars map[string]string
