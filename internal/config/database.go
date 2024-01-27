package config

type Database struct {
	Host     string `env:"DB_Host" envDefault:"localhost"`
	Port     string `env:"DB_Port" envDefault:"5342"`
	Sslmode  string `env:"DB_Sslmode" envDefault:"disable"`
	Name     string `env:"DB_Name" envDefault:"postgres"`
	User     string `env:"DB_User" envDefault:"postgres"`
	Password string `env:"DB_Password" envDefault:""`
}
