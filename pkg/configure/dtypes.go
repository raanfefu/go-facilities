package configure

import "crypto/tls"

const (
	SEPARATOR = "-"
)

type Configuration interface {
	// Registrar un componente que impletamenta las ConfigurationService
	// para leer parametros desde la linea de comandos
	RegistryService(key string, service interface{})
	// Carga la configuraciones de los diferentes componentes registrados con RegistryService
	// incluyendo archivos y cerficados y ejecuta las validaciones implementadas con PostConfiguration
	LoadConfiguration()
}

type ConfigurationService interface {
	PreConfiguration()
	PostConfiguration() error

	Key() *string
	SetKey(key string)
	StringVar(p *string, name string, value string, usage string)
	UintVar(p *uint, name string, value uint, usage string)
	GetNames() []string
	GetValue(name string) *CertificateArgsType
	GetJsonNames() []string
	GetJsonValue(name string) *JsonFileArgsType
}

type CertificateArgsType struct {
	Name            string
	ModuleName      string
	CertificatePath string
	KeyPath         string
	Certificate     *tls.Certificate
	Errors          []string
}

type JsonFileArgsType struct {
	Name       string
	ModuleName string
	Content    *interface{}
	Errors     []string
}
