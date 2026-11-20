package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)
//===================================================
// создаем пул и проверяем доступность постгре
// если все пройшло норм то код обязан закрыть пулу
func Open(ctx context.Context, connectionString string)(*pgxpool.Pool, error){
	config, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		// ошибка разбора может иметь в себе строку подключение с паролем
		return nil, errors.New("некорректные настройки подключения к PostgreSQL")
	}

	//начальное ограничение, потом буду менять но пока что остановлюсь на этом
	config.MaxConns = 4

	//ограничаем время для первоначального подключения
	connectCtx, cancel := context.WithTimeout(ctx,5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(connectCtx,config)
	if err != nil {return nil, fmt.Errorf("создание пула: %w", err)}
	if err := pool.Ping(connectCtx);err !=nil {
		//создание пулы не гарантирует соединение.. при ошибке освобождаем ресурсы
		pool.Close()
		return nil,fmt.Errorf("подключение: %w", err)
	}
	return pool,nil
}