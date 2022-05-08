package shared

import (
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type AppConfig struct {
	WeatherApiKey                  string  `env:"WEATHERAPI_KEY"`
	WeatherApiMaxRequestsPerSecond int     `env:"WEATHERAPI_MAX_REQUESTS_PER_SECOND" env-default:"1000"`
	WeatherApiClientTimeoutSec     int     `env:"WEATHERAPI_CLIENT_TIMEOUT_SEC" env-default:"7"`
	CacheMaxLifeWindowSec          int     `env:"CACHE_MAX_LIFE_WINDOW_SEC" env-default:"1600"`
	AppDebug                       bool    `env:"APP_DEBUG" env-default:"false"`
	AppNumDecimalPlaces            int     `env:"APP_NUM_DECIMAL_PLACES" env-default:"2"`
	AppMinDegreeStep               float32 `env:"APP_MIN_DEGREE_STEP" env-default:"0.1"`
}

var lock = &sync.Mutex{}
var Config AppConfig

func LoadConfig() *AppConfig {
	lock.Lock()
	defer lock.Unlock()

	if Config == (AppConfig{}) {

		err := cleanenv.ReadEnv(&Config)
		if err != nil {

		}
	}

	return &Config
}
