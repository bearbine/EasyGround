package httptransport

import (
	"log/slog"
	"net/http"
	"time"
	"crypto/rand"
	"context"
)
//сохраняем статус когда передаем ответ исходному w
type statusWriter struct{
	http.ResponseWriter
	status int
}
//====================
//передаем статус клиенту и сохраняем для лога
func (w *statusWriter) WriteHeader(code int){
	//окончательный статус отправляем ТОЛЬКО ОДИН раз
	if w.status != 0 {return}
	if code >= 100 && code < 200 && code != http.StatusSwitchingProtocols{
		w.ResponseWriter.WriteHeader(code)
		return
	}

	w.ResponseWriter.WriteHeader(code)
	w.status = code
}
//передаем тело ответа с учетом 200
func(w *statusWriter) Write(body []byte) (int,error){
	//запись тела без явного статуса будет 200
	if w.status == 0 {w.WriteHeader(http.StatusOK)}
	return w.ResponseWriter.Write(body)
}
//позволяем получить исходный w для ResponseController
func(w *statusWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}
//оборачиваем обработчик логирование запросов
func requestLogging(logger *slog.Logger,next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
								// w-запись ответа  r-входящий запрос
		start := time.Now()

		//создаем собственный ид для каджого входящего запроса
		requestID := rand.Text()
		w.Header().Set("X-Request-ID",requestID) //передаем ид клиенту до отправки статуса и тела 
		ctx := context.WithValue(r.Context(),requestIDKey{},requestID)//сохраняем ид в контексте запроса
		r = r.WithContext(ctx)

		

		//передаем обработчику обертку для статуса
		// если ничего не записано то ответ 200
		writer := &statusWriter{ResponseWriter:w}
		
		next.ServeHTTP(writer, r)
		status := writer.status
		if status == 0 {status = http.StatusOK}

		logger.Info("HTTP-запрос обработан",
		"method", r.Method,
		"path", r.URL.Path,
		"duration",time.Since(start),
		"request_id",requestID,
		"status", status,
		)
	})
}

