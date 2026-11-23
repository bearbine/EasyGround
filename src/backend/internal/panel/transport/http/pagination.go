package httptransport

import (
	"net/url"
	"strconv"

	"github.com/bearbine/EasyGround/src/backend/internal/panel/hosts"
)

//==================================
//разбираем параметры страницы из строки запроса
func parsePagination(rawQuery string)(int,int,error){
	query,err := url.ParseQuery(rawQuery)
	if err != nil {return 0,0,hosts.ErrInvalidPagination}

	limit := hosts.DefaultListLimit
	offset := 0

	//отсутствующий лимит оставляет значение по умолчанию
	if values,exists := query["limit"]; exists{
		if len(values) != 1{return 0,0,hosts.ErrInvalidPagination}
		
		limit,err =strconv.Atoi(values[0])
		if err != nil{return 0,0,hosts.ErrInvalidPagination}
	}
	//отсутствующий оффсет оставляет нулевое смещение
	if values,exists := query["offset"]; exists{
		if len(values) != 1{return 0,0,hosts.ErrInvalidPagination}

		offset,err = strconv.Atoi(values[0])
		if err != nil {return 0,0,hosts.ErrInvalidPagination}
	}

	if err := hosts.ValidatePagination(limit,offset); err != nil{return 0,0,err}
	return limit,offset,nil
}