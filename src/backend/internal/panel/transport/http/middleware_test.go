package httptransport

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"log/slog"
	"io"
)

//======================
//проверка запись тела без WriteHeader дает 200
func TestStatusWriterImplicitOK(t *testing.T){
	recorder := httptest.NewRecorder()
	w := &statusWriter{ResponseWriter:recorder}
	
	_, err := w.Write([]byte("hello"))
	if err != nil {
		t.Fatalf("ошибка записи ответа: %v",err)
	}

	//обертка должна запоминать отправленные маршруты
	if w.status != http.StatusOK {
		t.Fatalf("статус обертки: получил %d, хотя 200", w.status)
	}
	//клиенту должен уйти тот же статус и ИСХОДНОЕ тело
	if recorder.Code != http.StatusOK{
		t.Fatalf("статус ответа: получил %d, хотя 200", recorder.Code)
	}
	if recorder.Body.String() != "hello" {
		t.Fatalf("тело ответа: получил %q, хотя %q", recorder.Body.String(),"hello")
	}
}

//первый окончательный статус сохраняется при повторном WriteHeader
func TestStatusWriterKeepsFirstStatus(t *testing.T){
	recorder := httptest.NewRecorder()
	w := &statusWriter{ResponseWriter:recorder}

	w.WriteHeader(http.StatusCreated)
	w.WriteHeader(http.StatusInternalServerError)

	//проверяем статус, который будет записан в логах
	if w.status != http.StatusCreated{
		t.Fatalf("статус обертки: получил %d, хотя 201",w.status)
	}
	//проверяем статус, отправленный клиенту
	if recorder.Code != http.StatusCreated{
		t.Fatalf("статус ответа: получил %d, хотя 201", recorder.Code)
	}
}

//обработчик и клиент должны получить один и тот же ид
func TestRequestIDReachesHandler(t *testing.T){
	logger := slog.New(slog.NewTextHandler(io.Discard,nil))
	var handlerID string

	//тестовый обработчик читает ид из контекста
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerID = RequestID(r.Context())
		w.WriteHeader(http.StatusNoContent)
	})
	handler := requestLogging(logger,next)
	r := httptest.NewRequest(http.MethodGet, "/test", nil)

	//панель должна создать свой ид НО НИКАК НЕ ДОВЕРЯЯ ВХОДЯЩЕМУ ЗАПРОСУ
	r.Header.Set("X-Request-ID", "client-supplied-id")

	w := httptest.NewRecorder()
	handler.ServeHTTP(w,r)

	response := w.Result()
	defer response.Body.Close()

	if handlerID == "" {
		t.Fatal("обработчик не получил request_id")
	}
	if response.Header.Get("X-Request-ID") != handlerID{
		t.Fatal("ид в ответе не совпадает с ид в контексте")
	}
	if handlerID == "client-supplied-id"{
		t.Fatal("панель использования ид клиента вместо своего")
	}
}