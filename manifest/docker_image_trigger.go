package manifest

import "strings"

type DockerTrigger struct {
	Type     string
	Image    string `json:"image,omitempty" yaml:"image,omitempty"`
	Username string `json:"username,omitempty" yaml:"username,omitempty" secretAllowed:"true"`
	Password string `json:"password,omitempty" yaml:"password,omitempty" secretAllowed:"true"`
}

func (d DockerTrigger) GetTriggerAttempts() int {
	return 2
}

func (d DockerTrigger) MarshalYAML() (interface{}, error) {
	d.Type = "docker"
	return d, nil
}

func (d DockerTrigger) GetTriggerName() string {
	imageName := d.Image
	parts := strings.Split(imageName, "/")
	return parts[len(parts)-1]
}
