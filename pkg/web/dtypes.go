package web

import (
	"crypto/tls"
	"net/http"

	"github.com/raanfefu/go-facilities/pkg/common"
)

type Server interface {
	Init()
	AddEndpoint(path string, handler func(http.ResponseWriter, *http.Request), methods ...string)
	Start()
}

type ServerParameters struct {
	Mode      *common.ModeType
	ModeValue string
	Port      uint
	Certs     tls.Certificate
}
