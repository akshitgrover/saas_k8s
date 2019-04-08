package models

type ResourceQuota struct {
	ApiVersion string                       `yaml:"apiVersion"`
	Kind       string                       `yaml:"kind"`
	Metadata   map[string]string            `yaml:"metadata"`
	Spec       map[string]map[string]string `yaml:"spec"`
}

type Container struct {
	Name      string                       `yaml:"name"`
	Image     string                       `yaml:"image"`
	Resources map[string]map[string]string `yaml:"resources"`
	Ports     [](map[string]int)           `yaml:"ports"`
}

type PodMetadata struct {
	Name      string            `yaml:"name"`
	Namespace string            `yaml:"namespace"`
	Labels    map[string]string `yaml:"labels"`
}

type PodSpec struct {
	Containers []Container `yaml:"containers"`
}

type Pod struct {
	ApiVersion string      `yaml:"apiVersion"`
	Kind       string      `yaml:"kind"`
	Metadata   PodMetadata `yaml:"metadata"`
	Spec       PodSpec     `yaml:"spec"`
}

type User struct {
	Username string `bson:"username"`
	Tenant   string `bson:"tenant"`
}

type Tenant struct {
	Username string `bson:"username"`
}
