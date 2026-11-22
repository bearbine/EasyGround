package hosts

import (
	"testing"
	"strings"
	"errors"
)

//=========================
//проверяем удаление пробелов по краям
func TestNormalizeNameTrimsSpace(t *testing.T){
	got,err := NormalizeName("  test test ")
	if err != nil {t.Fatalf("неожиданная ошибка: %v",err)}
	if got != "test test"{t.Fatalf("Name:получил %q, хотя %q",got,"test test")}
}
//проверяем нейм которые должны быть отклонены
func TestNormalizeNameRejectsInvalid(t *testing.T){
	tests := []struct{
		name string
		input string
	}{
		{"empty", ""},
		{"whitespace", " \t\n "},
		{"too_long", strings.Repeat("я", 101)},
		{"nul", "сервер\x00"},
		{"invalid_utf8", string([]byte{0xff})},
	}

	for _,tt := range tests {
		t.Run(tt.name, func(t *testing.T){
			_, err := NormalizeName(tt.input)
			if !errors.Is(err, ErrInvalidName) {
			t.Fatalf("ошибка: получил %v,хотя %v",err, ErrInvalidName)}
		})
	}
}
//проверяем нейм граница длины
func TestNormalizeNameAcceptsBoundaryLengths(t *testing.T){
	tests := []struct{
		name string
		input string
	}{
		{"one_character", "я"},
		{"hundred_characters", strings.Repeat("я", 100)},
	}

	for _,tt := range tests{
		t.Run(tt.name, func(t *testing.T){
			got,err := NormalizeName(tt.input)
			if err != nil {t.Fatalf("неожиданная ошибка: %v",err)}
			if got != tt.input{t.Fatalf("Name:получил %q, хотя %q",got,tt.input)}
		})
	}
}