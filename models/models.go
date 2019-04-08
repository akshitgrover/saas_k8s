package models

type ResourceQuota struct {
	ApiVersion string                       `yaml:"apiVersion"`
	Kind       string                       `yaml:"kind"`
	Metadata   map[string]string            `yaml:"metadata"`
	Spec       map[string]map[string]string `yaml:"spec"`
}

type User struct {
	Username string `bson:"username"`
	Tenant   string `bson:"tenant"`
}

type Tenant struct {
	Username string `bson:"username"`
}
