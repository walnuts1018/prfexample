package config

type ScyllaDB struct {
	Host     string `env:"HOST" envDefault:"127.0.0.1"`
	Port     int    `env:"PORT" envDefault:"9042"`
	Keyspace string `env:"KEYSPACE" envDefault:"prf_example"`
}
