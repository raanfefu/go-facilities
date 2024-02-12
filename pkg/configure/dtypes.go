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
	// La implementaciond de este metodo contiene los paramtros que se desean capturar como flags
	// en la linea de comando asi como referencia de la variables que almacen el valor capturado
	PreConfiguration()
	// la implementacion de este metodo permite validar si se cargaron correctamente los valores
	// capturados durante le analisis de la paramtros recibidos por la linea de comandos como flags.
	PostConfiguration() error

	Key() *string
	SetKey(key string)
	StringVar(p *string, name string, value string, usage string)
	UintVar(p *uint, name string, value uint, usage string)
	GetNames() []string
	GetValue(name string) *CertificateArgsType
}

type CertificateArgsType struct {
	Name            string
	ModuleName      string
	CertificatePath string
	KeyPath         string
	Certificate     *tls.Certificate
}
