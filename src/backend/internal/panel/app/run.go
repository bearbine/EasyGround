package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)
//===============================
//запускаем сервер. работаем до ошибки или отмены
func Run(ctx context.Context, addr string) error {
	server := newServer(addr)

	//канал дает нам результат
	//в то время буфер дает отправить результат без простоя
	serverErr := make(chan error,1)

	//запустим сервак что бы следить за сигналом
	go func() {
		serverErr <- server.ListenAndServe()
	}()
	
	//смотрим ошибку или запрос 
	select{
	case err := <-serverErr:
		if errors.Is(err, http.ErrServerClosed) {return nil}
		return fmt.Errorf("запуск HTTP: %w",err)

	case <-ctx.Done():
		//переходим к норм запросу
	}

	// отдельный срок запроса где исходный ктх отменен
	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	//стопим соединение и ждем запросы
	if err := server.Shutdown(shutdownCtx); err != nil {
		//если время вышло, то закрываем
		_ = server.Close()
		return fmt.Errorf("остановка HTTP: %w", err)
	}

	return nil 
}