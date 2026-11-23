package httptransport

import (
	"mime"
	"net/http"
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