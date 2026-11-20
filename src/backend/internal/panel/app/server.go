package app 

import (
	"net/http"
	"time"
	"log/slog"

	httptransport "github.com/bearbine/EasyGround/src/backend/internal/panel/transport/http"
)
//===============================
// newServer создаст настройки без запуска
func newServer(addr string, logger *slog.Logger) *http.Server{
	return &http.Server{
		Addr:addr,
		Handler:httptransport.NewHandler(logger),

		//макс получение
		ReadHeaderTimeout:5*time.Second,
		//макс простой
		IdleTimeout:60*time.Second,

	}
}