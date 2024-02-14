package web

import (
	"crypto/tls"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/raanfefu/go-facilities/pkg/common"
	"github.com/raanfefu/go-facilities/pkg/configure"
)

type Server interface {
	Init()
	Start()
	GetRouter() *mux.Router
}

type impl struct {
	configure.DefaultConfiguraionService
	Params *ServerParameters
	Server *http.Server
	Router *mux.Router
}

type ServerParameters struct {
	Mode      *common.ModeType
	ModeValue string
	Port      uint
	Certs     tls.Certificate
}

type Handler struct {
	Func  http.HandlerFunc
	Route func(r *mux.Route)
}

func (h *Handler) Registry(s Server) {
	h.AddRoute(s.GetRouter())
}

func (h *Handler) AddRoute(r *mux.Router) {
	h.Route(r.NewRoute().HandlerFunc(h.Func))
}
