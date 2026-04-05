package processSubmission

// Workspace abstracts an isolated directory for a specific submission
type Workspace interface {
	BaseDir() string
	WriteFile(relativePath string, content []byte) (string, error)
	GetAbsolutePath(relativePath string) string
	Clean() error
}

// FileStorage abstracts the creation and management of Workspaces
type FileStorage interface {
	CreateWorkspace(submissionID string) (Workspace, error)
}
