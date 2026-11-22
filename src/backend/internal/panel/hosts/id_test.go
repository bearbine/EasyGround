package hosts

import (
	"strings"
	"testing"
	"errors"
)

//===============================
//проверка на допустимый uuid и приведение к нижнему регистру
func TestNormalizeID(t *testing.T){
	want := "550e8400-e29b-41d4-a716-446655440000"
	tests := []struct{
		name string
		input string
	}{
		{"lowercase",want},
		{"uppercase",strings.ToUpper(want)},
	}

	for _,tt := range tests{
		t.Run(tt.name,func(t *testing.T){
			got,err := NormalizeID(tt.input)
			if err != nil {t.Fatalf("неожиданная ошибка: %v",err)}
			if got != want{t.Fatalf("ID:получил %q, хотя %q",got,want)}
		})
	}
}
//проверяем UUID которые должны быть отклонены
func TestNormalizeIDRejectsInvalid(t *testing.T){
	tests := []struct{
		name string
		input string
	}{
		{"empty",""},
		{"too_short","550e8400-e29b-41d4-a716"},
		{"too_long","550e8400-e29b-41d4-a716-4466554400000"},
		{"without_hyphens","550e8400e29b41d4a716446655440000"},
		{"wrong_separator","550e8400_e29b-41d4-a716-446655440000"},
		{"invalid_hex","550e8400-e29b-41d4-a716-44665544000g"},
		{"extra_hyphens","550e8400-e29b-41d4-a716-4466554400--"},
		{"surrounding_spaces"," 550e8400-e29b-41d4-a716-446655440000 "},
	}

	for _,tt := range tests{
		t.Run(tt.name,func(t *testing.T){
			got,err := NormalizeID(tt.input)
			if !errors.Is(err,ErrInvalidID){t.Fatalf("ошибка:получил %v, хотя %v",err,ErrInvalidID)}
			if got != ""{t.Fatalf("ID:получил %q, хотя пустую строку",got)}
		})
	}
}