package configure

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
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

	for module, service := range c.services {
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
						log.Printf("Loading certificate from (%s) (%s)... Done ✓", module, name)
					}
				}
			}
		}

		// Load Jsons

		namejsons := service.(ConfigurationService).GetJsonNames()
		for _, name := range namejsons {
			jsonv := service.(ConfigurationService).GetJsonValue(name)
			if jsonv != nil {
				if jsonv.Name != "" {
					jsonv.Errors = make([]string, 0)
					fname := jsonv.Name
					jsonFile, err := os.Open(fname)
					if err != nil {
						jsonv.Errors = append(jsonv.Errors, fmt.Sprintf("failed to load file json -  %s", fname))
						log.Printf("Failed to load file: %v\n", err)
					} else {
						bytesValue, _ := io.ReadAll(jsonFile)
						err = json.Unmarshal(bytesValue, jsonv.Content)
						if err != nil {
							jsonv.Errors = append(jsonv.Errors, fmt.Sprintf("failed to parse file json - %s", fname))
							log.Printf("Failed to parse json file: %v\n", err)
						}
					}

				}

			}
			log.Printf("Loading json from (%s) (%s)... Done ✓", module, name)
		}

		// Post Validations
		err := service.(ConfigurationService).PostConfiguration()
		if err != nil {
			fmt.Printf("%s\n\n", err)
			flag.PrintDefaults()
			os.Exit(0)
		}
		log.Printf("Loading parameters (%s)... Done ✓", module)
	}

}
