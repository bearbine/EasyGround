package hosts

import (
	"errors"
	"testing"
)

//====================
//проверяем допустимые границы параметров страницы
func TestValidatePagination(t *testing.T){
	tests := []struct{
		name string
		limit int
		offset int
		err error
	}{
		{"default",DefaultListLimit,0,nil},
		{"minimum",1,0,nil},
		{"maximum",MaxListLimit,10,nil},
		{"zero_limit",0,0,ErrInvalidPagination},
		{"negative_limit",-1,0,ErrInvalidPagination},
		{"too_large_limit",MaxListLimit+1,0,ErrInvalidPagination},
		{"negative_offset",DefaultListLimit,-1,ErrInvalidPagination},
	}

	for _,tt := range tests{
		t.Run(tt.name,func(t *testing.T){
			err := ValidatePagination(tt.limit,tt.offset)
			if !errors.Is(err,tt.err){t.Fatalf("ошибка:получил %v, хотя %v",err,tt.err)}
		})
	}
}