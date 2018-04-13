package upstream

import (
	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/storage/file"
	"github.com/solo-io/glooctl/pkg/util"
)

func parseFile(filename string) (*v1.Upstream, error) {
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
