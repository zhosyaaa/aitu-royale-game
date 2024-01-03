package config

type Database struct {
	Host     string `env:"DB_Host" envDefault:"localhost"`
	Port     string `env:"DB_Port" envDefault:"5342"`
	Sslmode  string `json:"DB_Sslmode" envDefault:"disable"`
	Name     string `json:"DB_Name" envDefault:"postgres"`
	User     string `json:"DB_User" envDefault:"postgres"`
	Password string `json:"DB_Password" envDefault:""`
}
