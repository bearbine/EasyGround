package config

import (
	"errors"
	"os"
	"strings"
)
//настройки панели
type Config struct{
	DatabaseURL string
}
//=============================
//читает настройки из окружения процесса
func Load() (Config,error) {
	cfg := Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}

	// не запускаем приложение с отсутствующей настройкой базы
	if strings.TrimSpace(cfg.DatabaseURL) ==""{return Config{},errors.New("не задана переменная DATABASE_URL")}
	
	return cfg,nil
}