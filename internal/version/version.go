package version

import (
	"fmt"
	"runtime"
)

var (
	Name    = ""
	Version = ""
	Commit  = ""
)

// Info defines the application version information.
type Info struct {
	Name      string `json:"name" yaml:"name"`
	Version   string `json:"version" yaml:"version"`
	GitCommit string `json:"commit" yaml:"commit"`
	GoVersion string `json:"go" yaml:"go"`
}

func NewInfo() Info {
	return Info{
		Name:      Name,
		Version:   Version,
		GitCommit: Commit,
		GoVersion: fmt.Sprintf("go version %s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH),
	}
}

func (vi Info) String() string {
	return fmt.Sprintf(`%s: %s
git commit: %s
%s`,
		vi.Name, vi.Version, vi.GitCommit, vi.GoVersion,
	)
}
