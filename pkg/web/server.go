package web

import (
	"github.com/raanfefu/go-facilities/pkg/configure"
)

func NewServer(conf configure.Configuration) Server {
	obj := &impl{
		Params: &ServerParameters{},
	}
	if conf != nil {
		conf.RegistryService("w", obj)
	}
	return obj
}
