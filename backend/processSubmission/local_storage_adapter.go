package processSubmission

import (
	"fmt"
	"os"
	"path/filepath"
)

type LocalWorkspace struct {
	baseDir string
}

func (w *LocalWorkspace) BaseDir() string {
	return w.baseDir
}

func (w *LocalWorkspace) WriteFile(relativePath string, content []byte) (string, error) {
	absPath := w.GetAbsolutePath(relativePath)
	err := os.MkdirAll(filepath.Dir(absPath), 0755)
	if err != nil {
		return "", err
	}
	err = os.WriteFile(absPath, content, 0644)
	return absPath, err
}

func (w *LocalWorkspace) GetAbsolutePath(relativePath string) string {
	return filepath.Join(w.baseDir, relativePath)
}

func (w *LocalWorkspace) Clean() error {
	return os.RemoveAll(w.baseDir)
}

type LocalTempStorageAdapter struct{}

func (a *LocalTempStorageAdapter) CreateWorkspace(submissionID string) (Workspace, error) {
	dir := filepath.Join(os.TempDir(), fmt.Sprintf("velox_sub_%s", submissionID))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	return &LocalWorkspace{baseDir: dir}, nil
}
