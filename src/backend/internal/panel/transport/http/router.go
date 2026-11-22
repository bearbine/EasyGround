package httptransport

import (
	"net/http"
	"log/slog"
)

//=======================
// NewHandler создает маршрутизатор хттп запросов панели
func NewHandler(logger *slog.Logger, store HostStore) http.Handler{
	mux := http.NewServeMux() 
	// направляем запрос в обработчик
	mux.HandleFunc("GET /healthz", health)
	mux.HandleFunc("GET /api/hosts", listHosts(store,logger))
	//создание хоста
	mux.HandleFunc("POST /api/hosts",createHost(store,logger))
	//получение хоста
	mux.HandleFunc("GET /api/hosts/{id}",getHost(store,logger))
	//переименование хоста
	mux.HandleFunc("PATCH /api/hosts/{id}",renameHost(store,logger))
	//удаление записи хоста
	mux.HandleFunc("DELETE /api/hosts/{id}",deleteHost(store,logger))
	//оборачивание всех маршрутов
	return requestLogging(logger,mux)
}