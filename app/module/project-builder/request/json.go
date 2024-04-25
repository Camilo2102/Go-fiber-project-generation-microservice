package request

type ProjectInfo struct {
	UserId      string   `json:"userId"`
	ProjectName string   `json:"projectName"`
	Modules     []Module `json:"modules,omitempty"`
}

type Module struct {
	ModuleName string  `json:"moduleName"`
	Config     string  `json:"config"`
	Models     []Model `json:"models,omitempty"`
}

type Model struct {
	ModelName     string      `json:"modelName"`
	AttributeList []Attribute `json:"attributeList"`
}

type Attribute struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	DefaultValue string `json:"defaultValue"`
}
