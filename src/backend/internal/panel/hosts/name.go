package hosts

import (
	"errors"
	"strings"
	"unicode/utf8"
)

// название хоста не прошло проверку
var ErrInvalidName = errors.New("некореектное название хоста")

//============================
//убираем пробелы по краям и проверяем название хоста
func NormalizeName(name string)(string,error){
	//отклоянем ненорм UTF-8 и нулевой символ если будет
	if !utf8.ValidString(name) || strings.ContainsRune(name, '\x00'){return "", ErrInvalidName}
	name = strings.TrimSpace(name)
	length := utf8.RuneCountInString(name)
	if length <1 || length >100 {return "", ErrInvalidName}

	return name,nil
}