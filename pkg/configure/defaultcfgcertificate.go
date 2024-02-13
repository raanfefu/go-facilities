package configure

import (
	"crypto/tls"
	"flag"
	"fmt"
)

func (d *DefaultConfiguraionService) X509KeyPairVar(p *tls.Certificate, name string, usage string) {
	if d.certificates == nil {
		d.certificates = make(map[string]*CertificateArgsType)
	}
	crt := CertificateArgsType{
		Certificate: p,
	}
	flag.StringVar(&crt.CertificatePath, fmt.Sprintf("%s%s%s%scrt", *d.key, SEPARATOR, name, SEPARATOR), "", usage)
	flag.StringVar(&crt.KeyPath, fmt.Sprintf("%s%s%s%skey", *d.key, SEPARATOR, name, SEPARATOR), "", usage)

	d.certificates[fmt.Sprintf("%s:%s", *d.key, name)] = &crt
}

func (d *DefaultConfiguraionService) GetNames() []string {
	if d.certificates != nil {
		keys := make([]string, 0)
		for k := range d.certificates {
			keys = append(keys, k)
		}
		return keys
	}
	return nil
}

func (d *DefaultConfiguraionService) GetValue(name string) *CertificateArgsType {
	if d.certificates != nil {
		return d.certificates[fmt.Sprintf("%s:%s", *d.key, name)]
	}
	return nil
}
