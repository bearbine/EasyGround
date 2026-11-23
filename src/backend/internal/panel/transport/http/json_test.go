package httptransport

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"io"
	"log/slog"
	"strings"
)

//=========================
//проверяем поддерживаемый тип содержимого запроса
func TestRequireJSON(t *testing.T){
	tests := []struct{
		name string
		contentType string
		status int
		calls int
	}{
		{"json","application/json",http.StatusNoContent,1},
		{"json_with_charset","application/json; charset=utf-8",http.StatusNoContent,1},
		{"uppercase","Application/JSON",http.StatusNoContent,1},
		{"missing","",http.StatusUnsupportedMediaType,0},
		{"text","text/plain",http.StatusUnsupportedMediaType,0},
		{"malformed","application/json; charset",http.StatusUnsupportedMediaType,0},
	}

	for _,tt := range tests{
		t.Run(tt.name,func(t *testing.T){
			calls := 0
			handler := requireJSON(func(w http.ResponseWriter,r *http.Request){
				calls++
				w.WriteHeader(http.StatusNoContent)
			})

			r := httptest.NewRequest(http.MethodPost,"/",nil)
			if tt.contentType != ""{r.Header.Set("Content-Type",tt.contentType)}
			w := httptest.NewRecorder()
			handler.ServeHTTP(w,r)

			if w.Code != tt.status{t.Fatalf("статус:получил %d, хотя %d",w.Code,tt.status)}
			if calls != tt.calls{t.Fatalf("вызовы обработчика:получил %d, хотя %d",calls,tt.calls)}
		})
	}
}
//проверяем защиту POST и PATCH через маршрутизатор
func TestHostRoutesRequireJSON(t *testing.T){
	tests := []struct{
		name string
		method string
		path string
		contentType string
	}{
		{"post_missing",http.MethodPost,"/api/hosts",""},
		{"post_text",http.MethodPost,"/api/hosts","text/plain"},
		{"patch_missing",http.MethodPatch,"/api/hosts/550e8400-e29b-41d4-a716-446655440000",""},
		{"patch_text",http.MethodPatch,"/api/hosts/550e8400-e29b-41d4-a716-446655440000","text/plain"},
	}

	for _,tt := range tests{
		t.Run(tt.name,func(t *testing.T){
			logger := slog.New(slog.NewTextHandler(io.Discard,nil))
			handler := NewHandler(logger,emptyHostStore{})

			//JSON правильный, отклонить запрос должен именно заголовок
			r := httptest.NewRequest(tt.method,tt.path,strings.NewReader(`{"name":"test-host"}`))
			if tt.contentType != ""{r.Header.Set("Content-Type",tt.contentType)}
			w := httptest.NewRecorder()
			handler.ServeHTTP(w,r)

			if w.Code != http.StatusUnsupportedMediaType{t.Fatalf("статус:получил %d, хотя %d",w.Code,http.StatusUnsupportedMediaType)}
		})
	}
}