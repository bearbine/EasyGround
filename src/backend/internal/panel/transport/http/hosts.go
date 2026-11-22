package httptransport

import (
	"context"
	"time"
	"encoding/json"
	"log/slog"
	"net/http"
	"errors"
	"io"

	"github.com/bearbine/EasyGround/src/backend/internal/panel/hosts"
)
//описываем данные запроса на создание хоста
type createHostRequest struct {Name string `json:"name"`}
//описываем возможность получить список хостов
type HostLister interface {List(ctx context.Context)([]hosts.Host, error)}
//определяем поля хоста в HTTPответе
type hostResponse struct{
	ID string `json:"id"`
	Name string `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
//описываем возможность создать хост
type HostCreator interface {Create(ctx context.Context,name string)(hosts.Host,error)}
//ообьединяем возможности хранилища для маршрутизатора
type HostStore interface{
	HostLister
	HostCreator
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
//создаем обработчик запроса на добавление хоста
func createHost(store HostCreator,logger *slog.Logger) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		r.Body = http.MaxBytesReader(w,r.Body,1024) //ограничение тело одним КБ
	

		//разбираеим джсон и запрещаем анкнойвн поля
		var request createHostRequest
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		if err := decoder.Decode(&request); err != nil{
			var sizeErr *http.MaxBytesError
			if errors.As(err,&sizeErr){http.Error(w,"слишком большой запрос",http.StatusRequestEntityTooLarge) 
			return
		}

		http.Error(w,"некоррект JSON",http.StatusBadRequest)
		return
		}

		//после первого JSON допускаем только конец тела запроса
		var extra any
		if err := decoder.Decode(&extra);
		!errors.Is(err,io.EOF){
			var sizeErr *http.MaxBytesError
			if errors.As(err,&sizeErr){
				http.Error(w,"слишком большой запрос",http.StatusRequestEntityTooLarge)
				return
			}
			http.Error(w,"запрос должен содержать один JSON",http.StatusBadRequest)
			return
		}

		//ограничение время создания хоста
		ctx,cancel := context.WithTimeout(r.Context(),3*time.Second)
		defer cancel()

		host,err := store.Create(ctx,request.Name)
		if errors.Is(err,hosts.ErrInvalidName){
			http.Error(w,"некоррект название хоста",http.StatusBadRequest)
			return
		}
		if err != nil {
			logger.Error("создание хоста","error",err,"request_id",RequestID(r.Context()))
			http.Error(w,"не удалось создать хост",http.StatusInternalServerError)
			return
		}

		//ретурним только паблик поля хоста
		response := hostResponse{
			ID:host.ID,
			Name:host.Name,
			CreatedAt:host.CreatedAt,
		}
		body,err := json.Marshal(response)
		if err != nil{
			logger.Error("подготовка ответа с созданным хостом","error",err,"request_id",RequestID(r.Context()))
			http.Error(w,"не удалось подготовить ответ",http.StatusInternalServerError)
			return
		}
		//говорим что хост создан и ретурним json
		w.Header().Set("Content-Type","application/json")
		w.WriteHeader(http.StatusCreated)
		if _,err := w.Write(body); err != nil{logger.Error("отправка созданного хоста","error",err,"request_id",RequestID(r.Context()))}
	}
}