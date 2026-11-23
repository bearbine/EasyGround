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
	getFn func(ctx context.Context,id string)(hosts.Host,error)
	renameFn func(ctx context.Context,id string,name string)(hosts.Host,error)
	deleteFn func(ctx context.Context,id string)error
	listFn func(ctx context.Context,limit int,offset int)([]hosts.Host,error)
}
//======================
//ретурним подготовленные данные без обращение к постре
func (s hostStoreStub) List(ctx context.Context,limit int,offset int)([]hosts.Host,error){
	if s.listFn != nil{return s.listFn(ctx,limit,offset)}
	return s.items,s.err
}
//поведение создания задается отдельно в каждом месте
func (s hostStoreStub) Create(ctx context.Context,name string)(hosts.Host,error){
	if s.createFn == nil {panic("неожиданный вызов Create у hostStoreStub")}
	return s.createFn(ctx,name)
}
//поведение получения задается отдельно в каждом месте
func (s hostStoreStub) Get(ctx context.Context,id string)(hosts.Host,error){
	if s.getFn == nil {panic("неожиданный вызов Get у hostStoreStub")}
	return s.getFn(ctx,id)
}
//поведение переименования задается отдельно в каждом месте
func (s hostStoreStub) Rename(ctx context.Context,id string,name string)(hosts.Host,error){
	if s.renameFn == nil {panic("неожиданный вызов Rename у hostStoreStub")}
	return s.renameFn(ctx,id,name)
}
//поведение удаления задается отдельно в каждом месте
func (s hostStoreStub) Delete(ctx context.Context,id string)error{
	if s.deleteFn == nil {panic("неожиданный вызов Delete у hostStoreStub")}
	return s.deleteFn(ctx,id)
	
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
//проверка первышения размера тела запроса
func TestCreateHostTooLarge(t *testing.T){
	tests := []struct{
		name string
		body string
	}{
		{"large_json",`{"name":"`+strings.Repeat("a",1024)+`"}`},
		{"large_trailing_space",`{"name":"test-host"}`+strings.Repeat(" ",1024)},
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

			if w.Code != http.StatusRequestEntityTooLarge{t.Fatalf("статус:получил %d, хотя %d",w.Code,http.StatusRequestEntityTooLarge)}
			if calls != 0{t.Fatalf("вызовы Create:получил %d, хотя 0",calls)}
		})
	}
}
//проверяем ответ API при неправильном названии хоста
func TestCreateHostInvalidName(t *testing.T){
	calls := 0
	store := hostStoreStub{
		createFn:func(ctx context.Context,name string)(hosts.Host,error){
			calls++
			return hosts.Host{},hosts.ErrInvalidName
		},
	}
	logger := slog.New(slog.NewTextHandler(io.Discard,nil))
	handler := NewHandler(logger,store)

	//передаем пустое название в теле пост запроса
	r := httptest.NewRequest(http.MethodPost,"/api/hosts", strings.NewReader(`{"name":""}`))
	r.Header.Set("Content-Type","application/json")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w,r)

	if w.Code != http.StatusBadRequest{t.Fatalf("статус: получил %d, хотя %d",w.Code,http.StatusBadRequest)}
	if calls != 1{t.Fatalf("вызовы Create:получил %d, хотя 1",calls)}
}
//проверяем ответ API при ошибке создания в хранилище
func TestCreateHostStorageError(t *testing.T){
	storeErr := errors.New("внутренняя ошибка подключения к PostgreSQL")
	calls := 0
	store := hostStoreStub{
		createFn:func(ctx context.Context,name string)(hosts.Host,error){
			calls++
			return hosts.Host{},storeErr
		},
	}
	logger := slog.New(slog.NewTextHandler(io.Discard,nil))
	handler := NewHandler(logger,store)

	//передаем название в теле пост запроса
	r := httptest.NewRequest(http.MethodPost,"/api/hosts", strings.NewReader(`{"name":"test-host"}`))
	r.Header.Set("Content-Type","application/json")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w,r)

	if w.Code != http.StatusInternalServerError{t.Fatalf("статус: получил %d, хотя %d",w.Code,http.StatusInternalServerError)}
	if calls != 1{t.Fatalf("вызовы Create:получил %d, хотя 1",calls)}

	//внутренние подробности не должны попадать клиенту
	if strings.Contains(w.Body.String(),storeErr.Error()){t.Fatalf("ответ содержит внутреннюю ошибку: %q",w.Body.String())}
	if got := strings.TrimSpace(w.Body.String()); got != "не удалось создать хост"{t.Fatalf("тело ответа:получил %q, хотя %q",got,"не удалось создать хост")}
}
//проверяем получение одного хоста через API
func TestGetHost(t *testing.T){
	want := hosts.Host{
		ID:"550e8400-e29b-41d4-a716-446655440000",
		Name:"test-host",
		CreatedAt:time.Date(2026,time.November,22,12,0,0,0,time.UTC),
	}
	calls := 0
	var receivedID string
	store := hostStoreStub{
		getFn:func(ctx context.Context,id string)(hosts.Host,error){
			calls++
			receivedID = id
			return want,nil
		},
	}
	logger := slog.New(slog.NewTextHandler(io.Discard,nil))
	handler := NewHandler(logger,store)

	//передаем ID хоста в адресе запроса
	r := httptest.NewRequest(http.MethodGet,"/api/hosts/"+want.ID,nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w,r)

	if w.Code != http.StatusOK{t.Fatalf("статус:получил %d, хотя %d",w.Code,http.StatusOK)}
	if calls != 1{t.Fatalf("вызовы Get:получил %d, хотя 1",calls)}
	if receivedID != want.ID{t.Fatalf("ID:получил %q, хотя %q",receivedID,want.ID)}
	if got := w.Header().Get("Content-Type"); got != "application/json"{t.Fatalf("Content-Type:получил %q, хотя application/json",got)}

	//разбираем хост из ответа
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
//проверяем ответы API при ошибках получения хоста
func TestGetHostErrors(t *testing.T){
	tests := []struct{
		name string
		id string
		err error
		status int
		body string
	}{
		{"invalid_id","not-a-uuid",hosts.ErrInvalidID,http.StatusBadRequest,"некорректный ID хоста"},
		{"not_found","550e8400-e29b-41d4-a716-446655440000",hosts.ErrNotFound,http.StatusNotFound,"хост не найден"},
		{"storage_error","550e8400-e29b-41d4-a716-446655440000",errors.New("внутренняя ошибка PostgreSQL"),http.StatusInternalServerError,"не удалось получить хост"},
	}

	for _,tt := range tests{
		t.Run(tt.name,func(t *testing.T){
			calls := 0
			store := hostStoreStub{
				getFn:func(ctx context.Context,id string)(hosts.Host,error){
					calls++
					return hosts.Host{},tt.err
				},
			}
			logger := slog.New(slog.NewTextHandler(io.Discard,nil))
			handler := NewHandler(logger,store)

			//передаем ID из текущего случая
			r := httptest.NewRequest(http.MethodGet,"/api/hosts/"+tt.id,nil)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w,r)

			if w.Code != tt.status{t.Fatalf("статус:получил %d, хотя %d",w.Code,tt.status)}
			if calls != 1{t.Fatalf("вызовы Get:получил %d, хотя 1",calls)}
			if got := strings.TrimSpace(w.Body.String()); got != tt.body{t.Fatalf("тело ответа:получил %q, хотя %q",got,tt.body)}
		})
	}
}
//проверяем переименование хоста через API
func TestRenameHost(t *testing.T){
	want := hosts.Host{
		ID:"550e8400-e29b-41d4-a716-446655440000",
		Name:"renamed-host",
		CreatedAt:time.Date(2026,time.November,22,12,0,0,0,time.UTC),
	}
	calls := 0
	var receivedID string
	var receivedName string
	store := hostStoreStub{
		renameFn:func(ctx context.Context,id string,name string)(hosts.Host,error){
			calls++
			receivedID = id
			receivedName = name
			return want,nil
		},
	}
	logger := slog.New(slog.NewTextHandler(io.Discard,nil))
	handler := NewHandler(logger,store)

	//передаем ID в адресе и новое название в теле запроса
	r := httptest.NewRequest(http.MethodPatch,"/api/hosts/"+want.ID,strings.NewReader(`{"name":"renamed-host"}`))
	r.Header.Set("Content-Type","application/json")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w,r)

	if w.Code != http.StatusOK{t.Fatalf("статус:получил %d, хотя %d",w.Code,http.StatusOK)}
	if calls != 1{t.Fatalf("вызовы Rename:получил %d, хотя 1",calls)}
	if receivedID != want.ID{t.Fatalf("ID:получил %q, хотя %q",receivedID,want.ID)}
	if receivedName != want.Name{t.Fatalf("Name:получил %q, хотя %q",receivedName,want.Name)}
	if got := w.Header().Get("Content-Type"); got != "application/json"{t.Fatalf("Content-Type:получил %q, хотя application/json",got)}

	//разбираем переименованный хост из ответа
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
//проверяем ответы API при ошибках переименования хоста
func TestRenameHostErrors(t *testing.T){
	tests := []struct{
		name string
		err error
		status int
		body string
	}{
		{"invalid_id",hosts.ErrInvalidID,http.StatusBadRequest,"некорректный ID хоста"},
		{"invalid_name",hosts.ErrInvalidName,http.StatusBadRequest,"некорректное название хоста"},
		{"not_found",hosts.ErrNotFound,http.StatusNotFound,"хост не найден"},
		{"storage_error",errors.New("внутренняя ошибка PostgreSQL"),http.StatusInternalServerError,"не удалось переименовать хост"},
	}

	for _,tt := range tests{
		t.Run(tt.name,func(t *testing.T){
			calls := 0
			store := hostStoreStub{
				renameFn:func(ctx context.Context,id string,name string)(hosts.Host,error){
					calls++
					return hosts.Host{},tt.err
				},
			}
			logger := slog.New(slog.NewTextHandler(io.Discard,nil))
			handler := NewHandler(logger,store)

			//хранилище возвращает ошибку из текущего случая
			r := httptest.NewRequest(http.MethodPatch,"/api/hosts/550e8400-e29b-41d4-a716-446655440000",strings.NewReader(`{"name":"renamed-host"}`))
			r.Header.Set("Content-Type","application/json")
			w := httptest.NewRecorder()
			handler.ServeHTTP(w,r)

			if w.Code != tt.status{t.Fatalf("статус:получил %d, хотя %d",w.Code,tt.status)}
			if calls != 1{t.Fatalf("вызовы Rename:получил %d, хотя 1",calls)}
			if got := strings.TrimSpace(w.Body.String()); got != tt.body{t.Fatalf("тело ответа:получил %q, хотя %q",got,tt.body)}
		})
	}
}
//проверяем что неправильное тело запроса не доходит до хранилища
func TestRenameHostInvalidBody(t *testing.T){
	tests := []struct{
		name string
		body string
		status int
	}{
		{"broken_json",`{"name":}`,http.StatusBadRequest},
		{"unknown_field",`{"name":"renamed-host","id":"123"}`,http.StatusBadRequest},
		{"two_values",`{"name":"renamed-host"} {"name":"other-host"}`,http.StatusBadRequest},
		{"large_json",`{"name":"`+strings.Repeat("a",1024)+`"}`,http.StatusRequestEntityTooLarge},
		{"large_trailing_space",`{"name":"renamed-host"}`+strings.Repeat(" ",1024),http.StatusRequestEntityTooLarge},
	}

	for _,tt := range tests{
		t.Run(tt.name,func(t *testing.T){
			calls := 0
			store := hostStoreStub{
				renameFn:func(ctx context.Context,id string,name string)(hosts.Host,error){
					calls++
					return hosts.Host{},nil
				},
			}
			logger := slog.New(slog.NewTextHandler(io.Discard,nil))
			handler := NewHandler(logger,store)

			//передаем тело запроса из текущего случая
			r := httptest.NewRequest(http.MethodPatch,"/api/hosts/550e8400-e29b-41d4-a716-446655440000",strings.NewReader(tt.body))
			r.Header.Set("Content-Type","application/json")
			w := httptest.NewRecorder()
			handler.ServeHTTP(w,r)

			if w.Code != tt.status{t.Fatalf("статус:получил %d, хотя %d",w.Code,tt.status)}
			if calls != 0{t.Fatalf("вызовы Rename:получил %d, хотя 0",calls)}
		})
	}
}
//проверяем удаление хоста через API
func TestDeleteHost(t *testing.T){
	id := "550e8400-e29b-41d4-a716-446655440000"
	calls := 0
	var receivedID string
	store := hostStoreStub{
		deleteFn:func(ctx context.Context,id string)error{
			calls++
			receivedID = id
			return nil
		},
	}
	logger := slog.New(slog.NewTextHandler(io.Discard,nil))
	handler := NewHandler(logger,store)

	//передаем ID хоста в адресе запроса
	r := httptest.NewRequest(http.MethodDelete,"/api/hosts/"+id,nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w,r)

	if w.Code != http.StatusNoContent{t.Fatalf("статус:получил %d, хотя %d",w.Code,http.StatusNoContent)}
	if calls != 1{t.Fatalf("вызовы Delete:получил %d, хотя 1",calls)}
	if receivedID != id{t.Fatalf("ID:получил %q, хотя %q",receivedID,id)}
	if w.Body.Len() != 0{t.Fatalf("тело ответа должно быть пустым, получил %q",w.Body.String())}
}
//проверяем ответы API при ошибках удаления хоста
func TestDeleteHostErrors(t *testing.T){
	tests := []struct{
		name string
		err error
		status int
		body string
	}{
		{"invalid_id",hosts.ErrInvalidID,http.StatusBadRequest,"некорректный ID хоста"},
		{"not_found",hosts.ErrNotFound,http.StatusNotFound,"хост не найден"},
		{"storage_error",errors.New("внутренняя ошибка PostgreSQL"),http.StatusInternalServerError,"не удалось удалить хост"},
	}

	for _,tt := range tests{
		t.Run(tt.name,func(t *testing.T){
			calls := 0
			store := hostStoreStub{
				deleteFn:func(ctx context.Context,id string)error{
					calls++
					return tt.err
				},
			}
			logger := slog.New(slog.NewTextHandler(io.Discard,nil))
			handler := NewHandler(logger,store)

			//хранилище возвращает ошибку из текущего случая
			r := httptest.NewRequest(http.MethodDelete,"/api/hosts/550e8400-e29b-41d4-a716-446655440000",nil)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w,r)

			if w.Code != tt.status{t.Fatalf("статус:получил %d, хотя %d",w.Code,tt.status)}
			if calls != 1{t.Fatalf("вызовы Delete:получил %d, хотя 1",calls)}
			if got := strings.TrimSpace(w.Body.String()); got != tt.body{t.Fatalf("тело ответа:получил %q, хотя %q",got,tt.body)}
		})
	}
}