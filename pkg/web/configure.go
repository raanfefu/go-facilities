package web

import (
	"errors"

	"github.com/raanfefu/go-facilities/pkg/common"
)

func (s *impl) PreConfiguration() {
	s.StringVar(&s.Params.ModeValue, "mode", "", "mode value is https / http")
	s.UintVar(&s.Params.Port, "port", 0, "port using listing service")
	s.X509KeyPairVar(&s.Params.Certs, "tls", "certificado server")
}

func (s *impl) PostConfiguration() error {
	if common.StringEmptyOrNil(&s.Params.ModeValue) {
		return errors.New("mode is requeried")
	}
	vmode, err := common.ParseMode(&s.Params.ModeValue)
	if err != nil {
		return errors.New("mode value must be http / https")
	}
	s.Params.Mode = vmode

	if *vmode == common.TLS {
		if s.Params.Certs.PrivateKey == nil {
			return errors.New("no se puedo cargar el certificado 2")
		}
	}
	return nil
}
