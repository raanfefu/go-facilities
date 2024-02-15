package samplemod

import (
	"errors"

	"github.com/raanfefu/go-facilities/pkg/configure"
)

type Sample interface {
	Get() string
}

type impl struct {
	configure.DefaultConfiguraionService
	value string
}

func NewSampleMod(cfg configure.Configuration) Sample {
	a := &impl{}
	if cfg != nil {
		cfg.RegistryService("s", a)
	}

	return a
}

func (i *impl) Get() string {
	return i.value
}

func (s *impl) PreConfiguration() {
	s.StringVar(&s.value, "value", "t", "se requiere password")

}

func (s *impl) PostConfiguration() error {
	if s.value == "" {
		return errors.New("vacio")
	}

	return nil
}
