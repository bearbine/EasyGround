package hosts

import (
	"errors"
)
//хост с указанным ид не найден
var ErrNotFound = errors.New("хост не найден")
//ID хоста имеет неправильный формат
var ErrInvalidID = errors.New("некорректный ид хоста")

//=====================================
