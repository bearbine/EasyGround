package httptransport

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"context"
	"time"
	"encoding/json"
	"errors"

	"github.com/bearbine/EasyGround/src/backend/internal/panel/hosts"
)

//храним список хостов, который нужно ретурнить в тесте
type hostStoreStub struct{
	items []hosts.Host
	err error
	createFn func(ctx context.Context,name string)(hosts.Host,error)
}
//======================
//ретурним подготовленные данные без обращение к постре
func (s hostStoreStub) List(ctx context.Context)([]hosts.Host,error){return s.items,s.err}
//поведение создания задается отдельно в каждом месте
func (s hostStoreStub) Create(ctx context.Context,name string)(hosts.Host,error){
	if s.createFn == nil {panic("неожиданный вызов Create у hostStoreStub")}
	return s.createFn(ctx,name)
}
//проверяем ответ API, когда список хостов пуст
func TestListHostsEmpty(t *testing.T){
	logger := slog.New(slog.NewTextHandler(io.Discard,nil))
	handler := NewHandler(logger,emptyHostStore{})

	//создаем запрос и записываем ответ в память
	r := httptest.NewRequest(http.MethodGet,"/api/hosts",nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w,r)
	if w.Code != http.StatusOK{t.Fatalf("статус: получил %d,хотя %d", w.Code,http.StatusOK)}

	//тело ответа содержит джсон
	if got := w.Header().Get("Content-Type"); got != "application/json"{t.Fatalf("Content-Type: получил %q, хотя application/json",got)}
	//пустой список должен быть массивом
	if got := strings.TrimSpace(w.Body.String()); got != "[]"{t.Fatalf("тело ответа: получил %q, хотя []",got)}
}
//проверка выдачу списка с одним хостом
func TestListHostsWithData(t *testing.T){
	createdAt := time.Date(2026, time.November,22,12,0,0,0,time.UTC)
	store := hostStoreStub{
		items:[]hosts.Host{
			{
				ID: "8493234-432-423",
				Name:"testhost",
				CreatedAt: createdAt,
			},
		},
	}
	logger := slog.New(slog.NewTextHandler(io.Discard,nil))
	handler := NewHandler(logger,store)
	r := httptest.NewRequest(http.MethodGet, "/api/hosts",nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w,r)

	//проверка на ли успешно
	if w.Code != http.StatusOK{
		t.Fatalf("статус: получил %d, хотя %d", w.Code,http.StatusOK)
	}
	//разбор джсона по публичнім названиям полей апишки
	var response []struct{
		ID string `json:"id"`
		Name string `json:"name"`
		CreatedAt time.Time `json:"created_at"`
	}
	if err := json.Unmarshal(w.Body.Bytes(),&response); err != nil {t.Fatalf("не удалось разобрать JSON: %v",err)}
	//учтем что при ретурне должен быть один готовый хост
	if len(response) != 1 {t.Fatalf("количество хостов: получил %d,ожидал 1",len(response))}
	
	got := response[0]
	want := store.items[0]
	if got.ID != want.ID{t.Fatalf("ID:получил %q, хотя %q", got.ID, want.ID)}
	if got.Name != want.Name{t.Fatalf("Name:получил %q, хотя %q", got.Name, want.Name)}
	if !got.CreatedAt.Equal(want.CreatedAt){t.Fatalf("CreatedAt:получил %v, хотя %v", got.CreatedAt, want.CreatedAt)}
}
//проверка ответа апишки при ошибке в хранилище
func TestListHostsStorageError(t *testing.T){
	storeErr := errors.New("внутрення ошибка подключение к PostgreSQL")
	store := hostStoreStub{err:storeErr,}
	logger := slog.New(slog.NewTextHandler(io.Discard,nil))
	handler := NewHandler(logger,store)

	r := httptest.NewRequest(http.MethodGet,"/api/hosts",nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w,r)

	//ошибка хранилища должна быть 500
	if w.Code != http.StatusInternalServerError{
		t.Fatalf("статус: получил %d, хотя %d", w.Code,http.StatusInternalServerError)
	}
	//внутренности ошибки не должны быть у клиента. важно!
	if strings.Contains(w.Body.String(),storeErr.Error()){t.Fatalf("ответ содержит внутренююю ошибку: %q", w.Body.String())}

	if got := strings.TrimSpace(w.Body.String()); got != "не удалось получить хосты"{t.Fatalf("тело ответа: получил %q, ожидал %q",got,"не удалось получить хосты")}
}
//проверка на успех создания хоста через апишку
func TestCreateHost(t *testing.T){
	want := hosts.Host{
		ID:"88989989-43434343-43",
		Name:"test-host",
		CreatedAt:time.Date(2026,time.November,22,12,0,0,0,time.UTC),
	}
	//запоминаем вызов хранилища и переданное название
	
	calls := 0
	var receivedName string
	
	store := hostStoreStub{
		createFn:func(ctx context.Context,name string)(hosts.Host,error){
			calls++
			receivedName = name
			return want,nil
		},
	}
	logger := slog.New(slog.NewTextHandler(io.Discard,nil))
	handler := NewHandler(logger,store)

	//передаем название в теле пост запроса
	r := httptest.NewRequest(http.MethodPost,"/api/hosts",strings.NewReader(`{"name":"test-host"}`))
	r.Header.Set("Content-Type","application/json")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w,r)

	//создание должно ретурнить 201
	if w.Code != http.StatusCreated{t.Fatalf("статус:получил %d,хотя %d", w.Code,http.StatusCreated)}
	if calls != 1{t.Fatalf("вызовы create:получил %d, хотя 1",calls)}
	if receivedName != "test-host"{t.Fatalf("Name:получил %q, хотя %q",receivedName,"test-host")}
	if got := w.Header().Get("Content-Type"); got != "application/json"{t.Fatalf("Content-Type:получил %q, хотя application/json",got)}
	
	//разбиарем созданный хост из ответа
	var response struct{
		ID string `json:"id"`
		Name string `json:"name"`
		CreatedAt time.Time `json:"created_at"`
	}
	if err := json.Unmarshal(w.Body.Bytes(),&response); err != nil{t.Fatalf("не удалось разобрать JSON: %v",err)}
	if response.ID != want.ID{t.Fatalf("ID:получил %q, хотя %q",response.ID,want.ID)}
	if response.Name != want.Name{t.Fatalf("Name:получил %q, хотя %q",response.Name,want.Name)}
	if !response.CreatedAt.Equal(want.CreatedAt){t.Fatalf("CreatedAt:получил %v, хотя %v",response.CreatedAt,want.CreatedAt)}
}
//проверяем что неправильный JSON не доходит до хранилища
func TestCreateHostInvalidJSON(t *testing.T){
	tests := []struct{
		name string
		body string
	}{
		{"broken_json",`{"name":}`},
		{"unknown_field",`{"name":"test-host","id":"123"}`},
		{"two_values",`{"name":"test-host"} {"name":"other-host"}`},
	}
	for _,tt := range tests{
			t.Run(tt.name,func(t *testing.T){
			calls := 0
			store := hostStoreStub{
				createFn:func(ctx context.Context,name string)(hosts.Host,error){
					calls++
					return hosts.Host{},nil
				},
			}
			logger := slog.New(slog.NewTextHandler(io.Discard,nil))
			handler := NewHandler(logger,store)

			//передаем тело запроса из текущего случая 
			r := httptest.NewRequest(http.MethodPost,"/api/hosts", strings.NewReader(tt.body))
			r.Header.Set("Content-Type","application/json")
			w := httptest.NewRecorder()
			handler.ServeHTTP(w,r)

			if w.Code != http.StatusBadRequest{t.Fatalf("статус: получил %d, хотя %d", w.Code,http.StatusBadRequest)}
			if calls != 0{t.Fatalf("вызовы Create:получил %d, хотя 0",calls)}
		})
	}
} 