package httptransport

import "net/http"

// NewHandler создает маршрутизатор хттп запросов панели
func NewHandler() http.Handler{
	mux := http.NewServeMux() 
	// направляем запрос в обработчик
	mux.HandleFunc("GET /healthz", health)

	return mux
}