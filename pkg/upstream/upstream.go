package upstream

import (
	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/storage/file"
	"github.com/solo-io/glooctl/pkg/util"
)

// Options - represents CLI parameters related to upstream
type Options struct {
	Filename    string
	Output      string
	Template    string
	Interactive bool
}

// ParseFile parse a given file into an upstream
// If the filename is '-' it reads the standard input
func ParseFile(filename string) (*v1.Upstream, error) {
	var u v1.Upstream
	// special case: reading from stdin
	if filename == "-" {
		return &u, util.ReadStdinInto(&u)
	}
	err := file.ReadFileInto(filename, &u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
