package redis

import (
	"errors"

	"github.com/raanfefu/go-facilities/pkg/common"
)

func (s *impl) PreConfiguration() {
	s.BoolVar(&s.Parameters.RequiredPass, "required-pass", false, "se requiere password")
	s.StringVar(&s.Parameters.Host, "host", "", "servidor redis en formato <host:port>")
	s.UintVar(&s.Parameters.Port, "port", 6379, "password ")
	s.StringVar(&s.Parameters.Username, "user", "", "username ")
	s.StringVar(&s.Parameters.Password, "pass", "", "password ")

}

func (s *impl) PostConfiguration() error {
	if common.StringEmptyOrNil(&s.Parameters.Host) {
		return errors.New("hostname es requerido")
	}
	if s.Parameters.RequiredPass {
		if common.StringEmptyOrNil(&s.Parameters.Username) {
			return errors.New("username es requerido")
		}
		if common.StringEmptyOrNil(&s.Parameters.Password) {
			return errors.New("password es requerido")
		}
	}

	return nil
}
