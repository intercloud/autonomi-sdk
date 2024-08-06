package models

type CreateWorkspace struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateWorkspace struct {
	CreateWorkspace
}

type Workspace struct {
	BaseModel
	Name        string `json:"name"`
	Description string `json:"description"`
	AccountID   string `json:"accountId"`
}

type WorkspaceResponse struct {
	Data Workspace `json:"data"`
}
