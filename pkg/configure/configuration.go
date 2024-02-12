package configure

import (
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
)

type impl struct {
	services map[string]interface{}
}

func NewConfiguration() Configuration {
	return &impl{
		services: make(map[string]interface{}),
	}
}

func (c *impl) RegistryService(key string, service interface{}) {

	_, isConfigurationService := service.(ConfigurationService)
	if isConfigurationService {
		_, ok := c.services[key]
		if !ok {
			c.services[key] = service
			service.(ConfigurationService).SetKey(key)
			service.(ConfigurationService).PreConfiguration()
		} else {
			panic(errors.New("exist key services"))
		}
	} else {
		panic(errors.New("must be setup ConfigurationService"))
	}
}

func (c *impl) LoadConfiguration() {
	flag.Parse()

	for _, service := range c.services {
		// Load Certs

		names := service.(ConfigurationService).GetNames()
		for _, name := range names {
			crt := service.(ConfigurationService).GetValue(name)
			if crt != nil {
				if crt.CertificatePath != "" && crt.KeyPath != "" {
					vcrt, err := tls.LoadX509KeyPair(crt.CertificatePath, crt.KeyPath)
					if err != nil {
						log.Printf("Failed to load key pair: %v\n", err)
					} else {
						*crt.Certificate = vcrt
						log.Printf("Loading certificate TLS ... Done ✓")
					}
				}

			}

		}

		// Post Validations
		err := service.(ConfigurationService).PostConfiguration()
		if err != nil {
			fmt.Printf("%s\n\n", err)
			flag.PrintDefaults()
			os.Exit(0)
		}
	}
	log.Println("Loading parameters... Done ✓")
}
