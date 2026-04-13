// Модуль конфигурации
package config

import "os"

// Config содержит базовые настройки приложения
type Config struct {
	Port        string
	StoragePath string
}

// Load считывает конфигурацию
// Сначала ставим дефолтные значения, затем переписываем их,
// если в ENV что-то есть
func Load() *Config {
	cfg := &Config{
		Port:        "8080",       // Порт по умолчанию
		StoragePath: "tasks.json", // БД по умолчанию
	}

	if port := os.Getenv("HTTP_PORT"); port != "" {
		cfg.Port = port
	}

	if path := os.Getenv("STORAGE_PATH"); path != "" {
		cfg.StoragePath = path
	}
}
