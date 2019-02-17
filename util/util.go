package util

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GetProjectPath returns absolute path of current project.
func GetProjectPath() string {
	b, err := exec.Command("go", "env", "GOPATH").CombinedOutput()
	if err != nil {
		panic(err)
	}

	projectPath := ""
	for _, p := range filepath.SplitList(strings.TrimSpace(string(b))) {
		p = filepath.Join(p, filepath.FromSlash("/src/github.com/YuheiNakasaka/sandbox-golang"))
		if _, err = os.Stat(p); err == nil {
			projectPath = p
			break
		}
	}
	return projectPath
}
