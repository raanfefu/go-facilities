package configure

import (
	"flag"
	"fmt"
)

type DefaultConfiguraionService struct {
	key          *string
	Certificates map[string]*CertificateArgsType
}

func (d *DefaultConfiguraionService) Key() *string {
	return d.key
}

func (d *DefaultConfiguraionService) SetKey(pkey string) {
	if d.key == nil {
		d.key = &pkey
	}
}

func (d *DefaultConfiguraionService) StringVar(p *string, name string, value string, usage string) {
	flag.StringVar(p, fmt.Sprintf("%s%s%s", *d.key, SEPARATOR, name), value, usage)
}

func (d *DefaultConfiguraionService) UintVar(p *uint, name string, value uint, usage string) {
	flag.UintVar(p, fmt.Sprintf("%s%s%s", *d.key, SEPARATOR, name), value, usage)
}
