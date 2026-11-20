package httptransport

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

//TestRoutes проверяет маршруты через тот же обработчик
//котоырй потом подключи к хттпсерверу
func TestRoutes(t *testing.T){
	handler := NewHandler()

	//каждый елемент описывает отдельную проверку
	tests := []struct {
		name string
		method string
		path string
		status int
	}{
		{"health", http.MethodGet,"/healthz",http.StatusOK},
		{"wrong_method", http.MethodPost, "/healthz",http.StatusMethodNotAllowed},
		{"not_found",http.MethodGet,"/unknown",http.StatusNotFound},
	}

	for _, tt := range tests{
		t.Run(tt.name, func(t *testing.T) {
			//делаем запрос без тема и записываем в память
			r := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w,r)

			if w.Code != tt.status {
				t.Fatalf("статус: получил %d, хотя %d", w.Code, tt.status)
			}
		})
	}
}