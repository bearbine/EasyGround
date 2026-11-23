package httptransport

import (
	"mime"
	"net/http"
	"errors"
	"io"
	json "encoding/json/v2"

	
)

//=============================
//пропускаем запрос только с поддерживаемый типом внутренности
func requireJSON(next http.HandlerFunc) http.HandlerFunc{
	return func(w http.ResponseWriter,r *http.Request){
		mediaType,_,err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil || mediaType != "application/json"{
			http.Error(w,"ожидается Content-Type application/json",http.StatusUnsupportedMediaType)
			return
		}

		next(w,r)
	}
}
//читаем ограниченное тело и строго разбираем один JSON 
func readJSON(w http.ResponseWriter,r *http.Request,dst any)bool{
	body,err := io.ReadAll(http.MaxBytesReader(w,r.Body,1024))
	if err != nil{
		var sizeErr *http.MaxBytesError
		if errors.As(err,&sizeErr){http.Error(w,"слишком большой запрос",http.StatusRequestEntityTooLarge);return false}
		http.Error(w,"не удалось прочитать тело запроса",http.StatusBadRequest)
		return false
	}

	if err := json.Unmarshal(body,dst,json.RejectUnknownMembers(true));
	err != nil{
		http.Error(w,"некоррект JSON",http.StatusBadRequest)
		return false
	}

	return true
}