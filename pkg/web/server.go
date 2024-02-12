package web

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/raanfefu/go-facilities/pkg/configure"
)

type impl struct {
	configure.DefaultConfiguraionService
	Params *ServerParameters
	Server *http.Server
	Router *mux.Router
}

func NewServer(conf configure.Configuration) Server {
	obj := &impl{
		Params: &ServerParameters{},
	}
	conf.RegistryService("w", obj)
	return obj
}
