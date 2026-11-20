package app 

import (
	"net/http"
	"time"

	httptransport "github.com/bearbine/EasyGround/src/backend/internal/panel/transport/http"
)
//===============================
// newServer создаст настройки без запуска
func newServer(addr string) *http.Server{
	return &http.Server{
		Addr:addr,
		Handler:httptransport.NewHandler(),

		//макс получение
		ReadHeaderTimeout:5*time.Second,
		//макс простой
		IdleTimeout:60*time.Second,

	}
}