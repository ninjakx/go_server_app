package config

type Config struct {
	DB *DBConfig
}

type DBConfig struct {
	Dialect  string
	Host     string
	Port     int
	User     string
	Password string
	DBname   string
}

func GetConfig() *Config {
	return &Config{
		DB: &DBConfig{
			Dialect:  "postgres",
			Host:     "localhost",
			Port:     5432,
			User:     "kriti",
			Password: "nkx01",
			DBname:   "go_dummy",
		},
	}
}
