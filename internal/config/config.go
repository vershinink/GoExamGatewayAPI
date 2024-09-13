// Пакет для работы с файлом конфига.
package config

import (
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

// Структура конфига
type Config struct {
	News       string `yaml:"news_service"`
	Comments   string `yaml:"comments_service"`
	Censor     string `yaml:"censor_service"`
	HTTPServer `yaml:"http_server"`
}
type HTTPServer struct {
	Address      string        `yaml:"address"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

// MustLoad - инициализирует данные из конфиг файла. Путь к файлу берет из
// переменной окружения GATEWAY_CONFIG_PATH. Если не удается, то завершает
// приложение с ошибкой.
func MustLoad() *Config {
	configPath := os.Getenv("GATEWAY_CONFIG_PATH")
	if configPath == "" {
		log.Fatal("GATEWAY_CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("cannot read config file: %s, %s", configPath, err)
	}

	var cfg Config
	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		log.Fatalf("cannot decode config file: %s, %s", configPath, err)
	}

	return &cfg
}
