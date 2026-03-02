package config

type ScyllaDB struct {
	Endpoint string `env:"ENDPOINT" envDefault:"localhost:9042"`
	Keyspace string `env:"KEYSPACE" envDefault:"prf_example"`
}
