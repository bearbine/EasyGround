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

	"github.com/bearbine/EasyGround/src/backend/internal/panel/hosts"
)

//храним список хостов, который нужно ретурнить в тесте
type hostStoreStub struct{items []hosts.Host}
//======================
//ретурним подготовленные данные без обращение к постре
func (s hostStoreStub) List(ctx context.Context)([]hosts.Host,error){return s.items,nil}
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