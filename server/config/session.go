package config

type Session struct {
	Issuer string `env:"ISSUER" envDefault:"https://prfexample.walnuts.dev"`
}
