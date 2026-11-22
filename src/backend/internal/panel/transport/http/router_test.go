package httptransport

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"log/slog"
	"io"
	"context"

	"github.com/bearbine/EasyGround/src/backend/internal/panel/hosts"
)
//подменя хранилища
type emptyHostStore struct{}
//============================
//тест маршрутов не требует постгре
func (emptyHostStore) List(ctx context.Context)([]hosts.Host,error){return []hosts.Host{},nil}
//эта подмена исполдьзуется только для проверки чтения
func (emptyHostStore) Create(ctx context.Context,name string)(hosts.Host,error){panic("неожиданный вызов Create у emptyHostStore")}
//получение одного хоста в этой подмене не предусмотрено
func (emptyHostStore) Get(ctx context.Context,id string)(hosts.Host,error){panic("неожиданный вызов Get у emptyHostStore")}
//переименование в этой подмене не предусмотрено
func (emptyHostStore) Rename(ctx context.Context,id string,name string)(hosts.Host,error){panic("вызов Rename у emptyHostStore")}
//проверяет маршруты через тот же обработчик
//котоырй потом подключи к хттпсерверу
func TestRoutes(t *testing.T){
	//в тестах маршрутов отбрасываем логи что бы не засорять вывод
	logger := slog.New(slog.NewTextHandler(io.Discard,nil))
	handler := NewHandler(logger, emptyHostStore{})

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