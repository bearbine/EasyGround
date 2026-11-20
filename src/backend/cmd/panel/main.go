package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/bearbine/EasyGround/src/backend/internal/panel/app"
	"github.com/bearbine/EasyGround/src/backend/internal/panel/config"
)

//========================
func main(){
	os.Exit(run())
}

//ретурним код завершения и позваляем делать дефер перед выходом
func run() int {
	//адрес можно менять параметром -addr при запуске(не забыть)
	addr := flag.String("addr", "127.0.0.1:8080", "HTTP")
	flag.Parse()

	//структурированные логи
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	//читаем и чекаем настройки до запуска
	cfg,err := config.Load()
	if err != nil{
		logger.Error("ошибка настроек","error",err)
		return 1
	}

	//CTRL+C иЛИ SIGTERM ОТМЕНЯЕТ КОНТЕКСТ
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	logger.Info("запуск","addr", *addr)
	if err := app.Run(ctx, *addr,cfg,logger); err != nil {
		logger.Error("завершение с ошибкой","error",err)
		return 1
	}
	logger.Info("остановлено")
	return 0
}