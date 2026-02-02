package config

type ConfigDatabase struct {
	Host string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port string `yaml:"port" env:"DB_PORT" env-default:"5432"`
	Name string `yaml:"user" env:"DB_USER" env-required:"true"`
	User string `yaml:"user" env:"DB_USER" env-required:"true"`
	Password string `yaml:"password" env:"DB_PASSWORD" env-required:"true"`
}
