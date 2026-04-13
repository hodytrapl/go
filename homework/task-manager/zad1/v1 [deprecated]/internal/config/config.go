package config

import (
	"os"
)

type Config struct {
	Port string
	StoragePath string
}

func Load() *Config{
	cfg :=&Config{
		Port :"8080"
		StoragePath:"tasks.json"
	}

	if port:=os.Getenv("HTTP_PORT");port!=""{
		cfg.Port=port
	}

	if path:=os.Getenv("STORAGE_PATH");path!=""{
		cfg.StoragePath=path
	}
}