package domain

// ObjectCode type
type ObjectCode string

// Object describes object configuration
type Object struct {
	Code        ObjectCode         `json:"code"`
	Owner       string             `json:"owner"`
	ProjectCode ProjectCode        `json:"project_code" bson:"project_code"`
	EnvCode     EnvironmentCode    `json:"env_code" bson:"env_code"`
	Description string             `json:"description"`
	Inherits    *ObjectInheritance `json:"inherits"`
	Parameters  []*Parameter       `json:"parameters"`
}

// ObjectInheritance type
type ObjectInheritance struct {
	ProjectCode ProjectCode     `json:"project_code" bson:"project_code"`
	EnvCode     EnvironmentCode `json:"env_code" bson:"env_code"`
	ObjectCode  ObjectCode      `json:"object_code" bson:"object_code"`
}
