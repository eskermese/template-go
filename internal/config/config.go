package config

import (
	"fmt"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
)

const (
	EnvLocal = "local"
	Prod     = "prod"
)

type Config struct {
	*Database
	*Servers
}

type Database struct {
	DSURL string `short:"u" long:"ds-url" env:"DATASTORE_URL" description:"DataStore URL (format: postgresql://postgres:qwerty123@product-db/postgres)" required:"false" default:"postgresql://postgres:qwerty123@localhost:5432/postgres"`
}

type Servers struct {
	ClientHTTP *ClientServiceConfigs

	Environment string `long:"env" env:"ENVIRONMENT" description:"app environment" default:"local"`
	JWTKey      string `long:"jwt-key" env:"JWT_KEY" description:"JWT secret key" required:"false" default:"some-secret"`
}

type ClientServiceConfigs struct {
	Schema             string        `long:"schema" env:"HTTP_SCHEMA" description:"HTTP url schema" required:"false" default:"http"`
	Host               string        `long:"host" env:"HTTP_HOST" description:"HTTP host" required:"false" default:"localhost"`
	Port               int           `long:"port" env:"HTTP_PORT" description:"HTTP port" required:"false" default:"8000"`
	ReadTimeout        time.Duration `long:"read-timeout" env:"HTTP_READ_TIMEOUT" description:"HTTP read timeout" required:"false" default:"10s"`
	WriteTimeout       time.Duration `long:"write-timeout" env:"HTTP_WRITE_TIMEOUT" description:"HTTP write timeout" required:"false" default:"10s"`
	MaxHeaderMegabytes int           `long:"max-header-mg" env:"HTTP_MAX_HEADER_MEGABYTES" description:"HTTP max header mg" required:"false" default:"1"`
	ListenAddr         string        `long:"listen" env:"LISTEN" description:"Listen Address (format: :8080|127.0.0.1:8080)" required:"false" default:":8000"`
}

func New() (*Config, error) {
	defer os.Clearenv()

	c := &Config{}
	p := flags.NewParser(c, flags.Default|flags.IgnoreUnknown)

	if _, err := p.Parse(); err != nil {
		return nil, fmt.Errorf("error parsing config options: %w", err)
	}

	return c, nil
}
