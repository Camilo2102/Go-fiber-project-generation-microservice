package request

type ProjectInfo struct {
	UserId      string      `json:"userId"`
	ProjectName string      `json:"projectName"`
	ProjectType string      `json:"projectType"`
	Models      []ModelInfo `json:"models,omitempty"`
}

type ModelInfo struct {
	ModelName     string      `json:"modelName"`
	AttributeList []Attribute `json:"attributeList"`
}

type Attribute struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	DefaultValue string `json:"defaultValue"`
}
