package httptransport

import (
	"testing"
	"errors"

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