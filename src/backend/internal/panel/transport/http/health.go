package httptransport

import "net/http"

// health проверяет что хттпчасть панели отвечает
// если что доступность до базы данных и хостов тут проверять не будем
func health(w http.ResponseWriter, r *http.Request){
	// w-запись ответа  r-входящий запрос
	w.WriteHeader(http.StatusOK) // ретурним 200 без тела ответа
}