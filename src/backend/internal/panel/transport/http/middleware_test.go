package httptransport

import (
	"net/http"
	"net/http/httptest"
	"testing"
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