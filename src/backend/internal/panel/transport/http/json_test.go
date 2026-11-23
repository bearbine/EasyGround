package httptransport

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"io"
	"log/slog"
	"strings"
	"context"
	
	"github.com/bearbine/EasyGround/src/backend/internal/panel/hosts"
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
//проверяем отклонение повторяющихся полей и поврежденных строк
func TestHostRoutesStrictJSON(t *testing.T){
	tests := []struct{
		name string
		body string
	}{
		{"duplicate_name",`{"name":"first","name":"second"}`},
		{"invalid_utf8","{\"name\":\"\xff\"}"},
		{"invalid_surrogate",`{"name":"\ud800"}`},
	}

	for _,method := range []string{http.MethodPost,http.MethodPatch}{
		for _,tt := range tests{
			t.Run(method+"/"+tt.name,func(t *testing.T){
				calls := 0
				store := hostStoreStub{
					createFn:func(ctx context.Context,name string)(hosts.Host,error){
						calls++
						return hosts.Host{},nil
					},
					renameFn:func(ctx context.Context,id string,name string)(hosts.Host,error){
						calls++
						return hosts.Host{},nil
					},
				}
				logger := slog.New(slog.NewTextHandler(io.Discard,nil))
				handler := NewHandler(logger,store)

				path := "/api/hosts"
				if method == http.MethodPatch{path += "/550e8400-e29b-41d4-a716-446655440000"}

				r := httptest.NewRequest(method,path,strings.NewReader(tt.body))
				r.Header.Set("Content-Type","application/json")
				w := httptest.NewRecorder()
				handler.ServeHTTP(w,r)

				if w.Code != http.StatusBadRequest{t.Fatalf("статус:получил %d, хотя %d",w.Code,http.StatusBadRequest)}
				if calls != 0{t.Fatalf("вызовы хранилища:получил %d, хотя 0",calls)}
			})
		}
	}
}
