package httptransport

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)
//======================
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