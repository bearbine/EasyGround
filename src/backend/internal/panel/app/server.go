package app 

import (
	"net/http"
	"time"
	
)
//===============================
// newServer создаст настройки без запуска
func newServer(addr string, handler http.Handler) *http.Server{
	return &http.Server{
		Addr:addr,
		Handler: handler,

		//макс получение
		ReadHeaderTimeout:5*time.Second,
		//макс простой
		IdleTimeout:60*time.Second,

	}
}