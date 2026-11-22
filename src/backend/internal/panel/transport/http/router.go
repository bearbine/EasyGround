package httptransport

import (
	"net/http"
	"log/slog"
)

//=======================
// NewHandler создает маршрутизатор хттп запросов панели
func NewHandler(logger *slog.Logger, store HostLister) http.Handler{
	mux := http.NewServeMux() 
	// направляем запрос в обработчик
	mux.HandleFunc("GET /healthz", health)
	mux.HandleFunc("GET /api/hosts", listHosts(store,logger))

	//оборачивание всех маршрутов
	return requestLogging(logger,mux)
}