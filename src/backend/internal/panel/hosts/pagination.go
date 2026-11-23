package hosts

//размер страницы по умолчанию
const (
	DefaultListLimit = 50
	MaxListLimit = 100
)
//==================
//проверяем размер страницы и количество пропускаемых записей
func ValidatePagination(limit int,offset int)error{
	if limit < 1 || limit > MaxListLimit || offset < 0{return ErrInvalidPagination}
	return nil
}