package httptransport

import (
	"testing"
	"errors"
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"

	"github.com/bearbine/EasyGround/src/backend/internal/panel/hosts"
)
//============================
//проверяем разбор допустимых параметров страницы
func TestParsePagination(t *testing.T){
	tests := []struct{
		name string
		query string
		limit int
		offset int
	}{
		{"defaults","",50,0},
		{"limit_only","limit=10",10,0},
		{"offset_only","offset=20",50,20},
		{"both","limit=2&offset=3",2,3},
		{"maximum","limit=100",100,0},
		{"encoded","limit=%31%30&offset=%32",10,2},
	}

	for _,tt := range tests{
		t.Run(tt.name,func(t *testing.T){
			limit,offset,err := parsePagination(tt.query)
			if err != nil {t.Fatalf("неожиданная ошибка: %v",err)}
			if limit != tt.limit{t.Fatalf("limit:получил %d, хотя %d",limit,tt.limit)}
			if offset != tt.offset{t.Fatalf("offset:получил %d, хотя %d",offset,tt.offset)}
		})
	}
}
//проверяем отклонение неправильных параметров страницы
func TestParsePaginationRejectsInvalid(t *testing.T){
	tests := []struct{
		name string
		query string
	}{
		{"empty_limit","limit="},
		{"empty_offset","offset="},
		{"zero_limit","limit=0"},
		{"negative_limit","limit=-1"},
		{"large_limit","limit=101"},
		{"negative_offset","offset=-1"},
		{"text_limit","limit=abc"},
		{"text_offset","offset=abc"},
		{"duplicate_limit","limit=10&limit=20"},
		{"duplicate_offset","offset=0&offset=1"},
		{"overflow","offset=999999999999999999999999999999"},
		{"invalid_encoding","limit=%zz"},
	}

	for _,tt := range tests{
		t.Run(tt.name,func(t *testing.T){
			_,_,err := parsePagination(tt.query)
			if !errors.Is(err,hosts.ErrInvalidPagination){t.Fatalf("ошибка:получил %v, хотя %v",err,hosts.ErrInvalidPagination)}
		})
	}
}
//проверяем передачу параметров страницы из API в хранилище
func TestListHostsPagination(t *testing.T){
	tests := []struct{
		name string
		query string
		limit int
		offset int
	}{
		{"defaults","",50,0},
		{"explicit","?limit=2&offset=3",2,3},
		{"limit_only","?limit=10",10,0},
		{"offset_only","?offset=20",50,20},
	}

	for _,tt := range tests{
		t.Run(tt.name,func(t *testing.T){
			calls := 0
			var receivedLimit int
			var receivedOffset int
			store := hostStoreStub{
				listFn:func(ctx context.Context,limit int,offset int)([]hosts.Host,error){
					calls++
					receivedLimit = limit
					receivedOffset = offset
					return []hosts.Host{},nil
				},
			}
			logger := slog.New(slog.NewTextHandler(io.Discard,nil))
			handler := NewHandler(logger,store)

			r := httptest.NewRequest(http.MethodGet,"/api/hosts"+tt.query,nil)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w,r)

			if w.Code != http.StatusOK{t.Fatalf("статус:получил %d, хотя %d",w.Code,http.StatusOK)}
			if calls != 1{t.Fatalf("вызовы List:получил %d, хотя 1",calls)}
			if receivedLimit != tt.limit{t.Fatalf("limit:получил %d, хотя %d",receivedLimit,tt.limit)}
			if receivedOffset != tt.offset{t.Fatalf("offset:получил %d, хотя %d",receivedOffset,tt.offset)}
		})
	}
}
//проверяем отказ до обращения к хранилищу
func TestListHostsRejectsInvalidPagination(t *testing.T){
	tests := []struct{
		name string
		query string
	}{
		{"zero_limit","limit=0"},
		{"large_limit","limit=101"},
		{"negative_offset","offset=-1"},
		{"empty_limit","limit="},
		{"text_limit","limit=abc"},
		{"duplicate_limit","limit=10&limit=20"},
		{"invalid_encoding","limit=%zz"},
	}

	for _,tt := range tests{
		t.Run(tt.name,func(t *testing.T){
			calls := 0
			store := hostStoreStub{
				listFn:func(ctx context.Context,limit int,offset int)([]hosts.Host,error){
					calls++
					return []hosts.Host{},nil
				},
			}
			logger := slog.New(slog.NewTextHandler(io.Discard,nil))
			handler := NewHandler(logger,store)

			r := httptest.NewRequest(http.MethodGet,"/api/hosts?"+tt.query,nil)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w,r)

			if w.Code != http.StatusBadRequest{t.Fatalf("статус:получил %d, хотя %d",w.Code,http.StatusBadRequest)}
			if calls != 0{t.Fatalf("вызовы List:получил %d, хотя 0",calls)}
		})
	}
}