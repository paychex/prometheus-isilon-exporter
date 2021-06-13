package isiconfig

import (
	"flag"

	"github.com/jamiealquiza/envy"
)

type isiConfig struct {
	UserName string
	Password string
	MgmtPort int
	IsiURL   string
}

type exporterConfig struct {
	BindAddress string
	BindPort    int
	LogLevel    string
	MultiQuery  bool
}

// Config is a container for settings modifiable by the user
type Config struct {
	ISI      isiConfig
	Exporter exporterConfig
}

var (
	isiMgmtPort   = flag.Int("mgmtport", 8080, "The port which isilon listens to for administration")
	isiUserName   = flag.String("username", "defaultUser", "Username")
	isiPassword   = flag.String("password", "defaultPass", "Password")
	listenAddress = flag.String("bindaddress", "localhost", "Exporter bind address")
	listenPort    = flag.Int("bindport", 9437, "Exporter bind port")
	isiURL        = flag.String("url", "", "Base URL of the Isilon management interface.  Normally something like https://my-isilon.something.x")
	multiQuery    = flag.Bool("multi", false, "Enable query endpoint")
)

func init() {

	envy.Parse("ISIENV") // looks for ISIENV_USERNAME, ISIENV_PASSWORD, ISIENV_BINDPORT etc
	flag.Parse()

}

// GetConfig returns an instance of Config containing the resulting parameters
// to the program
func GetConfig() *Config {
	return &Config{
		ISI: isiConfig{
			UserName: *isiUserName,
			Password: *isiPassword,
			MgmtPort: *isiMgmtPort,
			IsiURL:   *isiURL,
		},
		Exporter: exporterConfig{
			BindAddress: *listenAddress,
			BindPort:    *listenPort,
			MultiQuery:  *multiQuery,
		},
	}
}
