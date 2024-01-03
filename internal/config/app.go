package config

type App struct {
	DB    Database
	Redis RedisConfig
	Email EmailConfig
}
