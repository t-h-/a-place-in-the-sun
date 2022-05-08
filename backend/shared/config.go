package shared

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type AppConfig struct {
	WeatherApiKey                  string  `yaml:"WeatherApiKey" env:"WEATHERAPI_KEY"`
	WeatherApiMaxRequestsPerSecond int     `yaml:"WeatherApiMaxRequestsPerSecond" env:"WEATHERAPI_MAX_REQUESTS_PER_SECOND" env-default:"1000"`
	WeatherApiClientTimeoutSec     int     `yaml:"WeatherApiClientTimeoutSec" env:"WEATHERAPI_CLIENT_TIMEOUT_SEC" env-default:"7"`
	CacheMaxLifeWindowSec          int     `yaml:"CacheMaxLifeWindowSec" env:"CACHE_MAX_LIFE_WINDOW_SEC" env-default:"1600"`
	AppDebug                       bool    `yaml:"AppDebug" env:"APP_DEBUG" env-default:"false"`
	AppNumDecimalPlaces            int     `yaml:"AppNumDecimalPlaces" env:"APP_NUM_DECIMAL_PLACES" env-default:"2"`
	AppMinDegreeStep               float32 `yaml:"AppMinDegreeStep" env:"APP_MIN_DEGREE_STEP" env-default:"0.1"`
}

var Config AppConfig

// LoadConfig has to be run on startup. Thereafter, the var Config wil contain the values loaded from the environment.
// TODO Consider making passing a config instance to the components explicitly in the main().
func LoadConfigFromEnv() {
	if Config == (AppConfig{}) {

		err := cleanenv.ReadEnv(&Config)
		if err != nil {
			log.Fatalf("Could not load config from environment: %v", err)
		}
	}
}

func LoadConfigFromYaml(path string) {
	if Config == (AppConfig{}) {

		err := cleanenv.ReadConfig(path, &Config)
		if err != nil {
			log.Fatalf("Could not load config from config file: %v", err)
		}
	}
}
