package config

var AppConfig = &Config{}

type Config struct {
	BaseUrl     string `yaml:"base_url"`
	RedisConfig struct {
		DBHost   string `yaml:"db_host"`
		DBPort   int    `yaml:"db_port"`
		DBPasswd string `yaml:"db_passwd"`
		DB       int    `yaml:"db"`
	} `yaml:"redis_config"`
}
