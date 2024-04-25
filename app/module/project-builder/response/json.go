package response

type ProjectCreateInfo struct {
	Status  bool   `json:"status"`
	Phase   int    `json:"phase"`
	Message string `json:"message"`
}
