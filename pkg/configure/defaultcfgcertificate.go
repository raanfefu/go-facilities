package configure

import (
	"crypto/tls"
	"flag"
	"fmt"
)

func (d *DefaultConfiguraionService) X509KeyPairVar(p *tls.Certificate, name string, usage string) {
	crt := CertificateArgsType{
		Certificate: p,
	}
	flag.StringVar(&crt.CertificatePath, fmt.Sprintf("%s%s%s%scrt", *d.key, SEPARATOR, name, SEPARATOR), "", usage)
	flag.StringVar(&crt.KeyPath, fmt.Sprintf("%s%s%s%skey", *d.key, SEPARATOR, name, SEPARATOR), "", usage)
	if d.Certificates == nil {
		d.Certificates = make(map[string]*CertificateArgsType)
	}
	d.Certificates[name] = &crt
}

func (d *DefaultConfiguraionService) GetNames() []string {
	if d.Certificates != nil {
		keys := make([]string, 0)
		for k := range d.Certificates {
			keys = append(keys, k)
		}
		return keys
	}
	return nil
}

func (d *DefaultConfiguraionService) GetValue(name string) *CertificateArgsType {
	if d.Certificates != nil {
		return d.Certificates[name]
	}
	return nil
}
