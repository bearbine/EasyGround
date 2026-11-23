package httptransport

import (
	"context"
	"time"
	"encoding/json"
	"log/slog"
	"net/http"
	"errors"

	"github.com/bearbine/EasyGround/src/backend/internal/panel/hosts"
)
//описываем данные запроса на создание хоста
type createHostRequest struct {Name string `json:"name"`}
//описываем возможность получить список хостов
type HostLister interface {List(ctx context.Context,limit int,offset int)([]hosts.Host,error)}
//определяем поля хоста в HTTPответе
type hostResponse struct{
	ID string `json:"id"`
	Name string `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
//описываем возможность создать хост
type HostCreator interface {Create(ctx context.Context,name string)(hosts.Host,error)}
//описываем возможность получить один хост
type HostGetter interface {Get(ctx context.Context,id string)(hosts.Host,error)}
//описываем возможность переименовать хост
type HostRenamer interface {Rename(ctx context.Context,id string,name string)(hosts.Host,error)}
//описываем данные запроса на переименование хоста
type renameHostRequest struct {Name string `json:"name"`}
//описываем возхможность удалить запись хоста
type HostDeleter interface {Delete(ctx context.Context,id string)error}
//ообьединяем возможности хранилища для маршрутизатора
type HostStore interface{
	HostLister
	HostCreator
	HostGetter
	HostRenamer
	HostDeleter
}
//=====================================
// создаем обработчик с нужным хранилищем и логгером
func listHosts(store HostLister, logger *slog.Logger) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		//проверяем параметры страницы до обращения к хранилищу
		limit,offset,err := parsePagination(r.URL.RawQuery)
		if err != nil{
			http.Error(w,"некорректные параметры страницы",http.StatusBadRequest)
			return
		}
		//ограничение ожидание базы в секах
		ctx,cancel := context.WithTimeout(r.Context(),3*time.Second)//тут поставил 3 сек пока что
		defer cancel()

		items,err := store.List(ctx,limit,offset)
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
		//строго разбираем JSON через общую функцию
		var request createHostRequest
		if !readJSON(w,r,&request){return}

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
//получаем один хост по ид из адреса запроса
func getHost(store HostGetter,logger *slog.Logger) http.HandlerFunc{
	return func(w http.ResponseWriter,r *http.Request){
		ctx,cancel := context.WithTimeout(r.Context(),3*time.Second)
		defer cancel()

		host,err := store.Get(ctx,r.PathValue("id"))
		if errors.Is(err,hosts.ErrInvalidID){http.Error(w,"некорректный ID хоста",http.StatusBadRequest);return}
		if errors.Is(err,hosts.ErrNotFound){http.Error(w,"хост не найден",http.StatusNotFound);return}
		if err != nil{
			logger.Error("получение хоста","error",err,"request_id",RequestID(r.Context()))
			http.Error(w,"не удалось получить хост",http.StatusInternalServerError)
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
			logger.Error("подготовка ответа с хостом","error",err,"request_id",RequestID(r.Context()))
			http.Error(w,"не удалось подготовить ответ",http.StatusInternalServerError)
			return
		}

		//ретурним хост в json
		w.Header().Set("Content-Type","application/json")
		if _,err := w.Write(body); err != nil{logger.Error("отправка хоста","error",err,"request_id",RequestID(r.Context()))}
	}
}
//переименовываем хост по ИД из адреса запроса
func renameHost(store HostRenamer,logger *slog.Logger) http.HandlerFunc{
	return func(w http.ResponseWriter,r *http.Request){
		//строго разбираем JSON через общую функцию
		var request renameHostRequest
		if !readJSON(w,r,&request){return}

		ctx,cancel := context.WithTimeout(r.Context(),3*time.Second)
		defer cancel()

		host,err := store.Rename(ctx,r.PathValue("id"),request.Name)
		if errors.Is(err,hosts.ErrInvalidID){http.Error(w,"некорректный ID хоста",http.StatusBadRequest); return}
		if errors.Is(err,hosts.ErrInvalidName){http.Error(w,"некорректное название хоста",http.StatusBadRequest);return}
		if errors.Is(err,hosts.ErrNotFound){http.Error(w,"хост не найден",http.StatusNotFound);return}
		if err != nil{
			logger.Error("переименование хоста","error",err,"request_id",RequestID(r.Context()))
			http.Error(w,"не удалось переименовать хост",http.StatusInternalServerError)
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
			logger.Error("подготовка ответа с переименованным хостом","error",err,"request_id",RequestID(r.Context()))
			http.Error(w,"не удалось подготовить ответ",http.StatusInternalServerError)
			return
		}

		//ретурним переименованный хост в json
		w.Header().Set("Content-Type","application/json")
		if _,err := w.Write(body); err != nil{logger.Error("отправка переименованного хоста","error",err,"request_id",RequestID(r.Context()))}

	}
}
//удаляем хост по ID 
func deleteHost(store HostDeleter,logger *slog.Logger) http.HandlerFunc{
	return func(w http.ResponseWriter,r *http.Request){
		ctx,cancel := context.WithTimeout(r.Context(),3*time.Second)
		defer cancel()

		err := store.Delete(ctx,r.PathValue("id"))
		if errors.Is(err,hosts.ErrInvalidID){http.Error(w,"некорректный ID хоста",http.StatusBadRequest);return}
		if errors.Is(err,hosts.ErrNotFound){http.Error(w,"хост не найден",http.StatusNotFound); return}
		if err != nil{
			logger.Error("удаление хоста","error",err,"request_id",RequestID(r.Context()))
			http.Error(w,"не удалось удалить хост", http.StatusInternalServerError)
			return
		}

		//запись удалена, тело ответа не нуждаетса
		w.WriteHeader(http.StatusNoContent)
	}
}