package models

type CreateWorkspace struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateWorkspace struct {
	Name        string `json:"name"`
	Description string `json:"description"`
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

type WorkspacesResponse struct {
	Data []Workspace `json:"data"`
}
