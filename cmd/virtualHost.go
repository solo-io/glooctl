package cmd

import (
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/glooctl/platform"
)

type VHost struct {
	vhost v1.VirtualHost
}

func NewVHost() *VHost {
	return &VHost{vhost: v1.VirtualHost{}}
}

func (vh *VHost) Get() *v1.VirtualHost {
	return &vh.vhost
}

func GetVhostParams() *platform.VhostParams {
	return &platform.VhostParams{}
}
