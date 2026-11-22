package httptransport

import (
	"context"
	"time"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/bearbine/EasyGround/src/backend/internal/panel/hosts"
)
//описываем возможность получить список хостов
type HostLister interface {List(ctx context.Context)([]hosts.Host, error)}
//определяем поля хоста в HTTPответе
type hostResponse struct{
	ID string `json:"id"`
	Name string `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
//=====================================
// создаем обработчик с нужным хранилищем и логгером
func listHosts(store HostLister, logger *slog.Logger) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		//ограничение ожидание базы в секах
		ctx,cancel := context.WithTimeout(r.Context(),3*time.Second)//тут поставил 3 сек пока что
		defer cancel()

		items,err := store.List(ctx)
		if err != nil{
			logger.Error("получение хостов","error",err,"request_id",RequestID(r.Context()),)
			http.Error(w, "не удалось получить хосты", http.StatusInternalServerError)
			return
		}
		// собираем поля которые РАЗРЕШЕНЫ В API
		response := make([]hostResponse,0,len(items))
		for _,host := range items{
			response = append(response, hostResponse{
				ID:host.ID,
				Name:host.Name,
				CreatedAt:host.CreatedAt,
			})
		}

		//готовим джсон до отправки заголовка
		// если при каком то будет ошбикиа то ретурним статус 500
		body, err := json.Marshal(response)
		if err != nil{
			logger.Error("подготовка ответа хостами","error",err,"request_id",RequestID(r.Context()),)
			http.Error(w,"не удалось подготовить ответ", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type","application/json")
		if _,err := w.Write(body); err != nil {logger.Error("отправка списка хостов","error",err,"request_id",RequestID(r.Context()),)}
	}
}