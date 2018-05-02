package bootstrap

import (
	"github.com/solo-io/gloo/pkg/bootstrap"
)

type Options struct {
	bootstrap.Options
	IngressOptions IngressOptions
}

type IngressOptions struct {
	BindAddress string
	Port        int
	SecurePort  int
}
