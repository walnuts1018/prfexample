package config

import "net/url"

type Server struct {
	Port              int      `env:"PORT" envDefault:"8080" validate:"gte=1,lte=65535"`
	Origin            *url.URL `env:"ORIGIN" envDefault:"https://prfexample.walnuts.dev"`
	AdditionalOrigins []string `env:"ADDITIONAL_ORIGINS" envSeparator:","`
}
