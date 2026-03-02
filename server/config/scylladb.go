package config

type ScyllaDB struct {
	Endpoint string `env:"ENDPOINT" envDefault:"localhost:9042"`
	Keyspace string `env:"KEYSPACE" envDefault:"prfexample"`

	CACertPath     string `env:"CA_CERT_PATH"`
	ClientCertPath string `env:"CLIENT_CERT_PATH"`
	ClientKeyPath  string `env:"CLIENT_KEY_PATH"`

	User     string `env:"USER" envDefault:"cassandra"`
	Password string `env:"PASSWORD" envDefault:"cassandra"`
}
