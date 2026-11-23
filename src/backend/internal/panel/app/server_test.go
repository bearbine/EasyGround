package app

import (
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

//============================
//проверяем прекращение чтения незавершенного тела запроса
func TestServerReadTimeout(t *testing.T){
	readErr := make(chan error,1)
	handler := http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
		_,err := io.Copy(io.Discard,r.Body)
		readErr <- err
	})
	//берем настройки настоящего сервера
	server := httptest.NewUnstartedServer(handler)
	server.Config = newServer("",handler)
	
	//сокращаем настроенній срок только для быстрого теста
	server.Config.ReadTimeout /= 10
	server.Start()
	defer server.Close()

	//подключаемся напрямую что бы управлять отправкой тела
	conn,err := net.DialTimeout("tcp",server.Listener.Addr().String(),3*time.Second)
	if err != nil{t.Fatalf("подключение к серверу %v",err)}
	defer conn.Close()
	if err := conn.SetWriteDeadline(time.Now().Add(3*time.Second)); err != nil{t.Fatalf("установка срока записи: %v",err)}

	//обещаем 10 байт тела но отправляем только один
	request := "POST / HTTP/1.1\r\nHost: localhost\r\nContent-Length: 10\r\n\r\n{"
	if _,err := io.WriteString(conn,request); err != nil{t.Fatalf("отправка запроса: %v",err)}
	//ждем именно ошибку таймаута чтения на сервере
	select{
		case err := <-readErr:
			var timeoutErr net.Error
			if !errors.As(err,&timeoutErr) || !timeoutErr.Timeout(){t.Fatalf("ожидали таймаут чтения, получили: %v",err)}
		case <-time.After(5*time.Second):
			t.Fatal("сервер не прекратил чтение тела запроса")
	}
}
//проверяем прекращение записи клиенту который не читает ответ
func TestServerWriteTimeout(t *testing.T){
	writeErr := make(chan error,1)

	handler := http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
		chunk := make([]byte,32*1024)
		for{
			if _,err := w.Write(chunk); err != nil{
				writeErr <- err
				return
			}
		}
	})
	//берем настройки настоящего сервера
	server := httptest.NewUnstartedServer(handler)
	server.Config = newServer("",handler)

	//сокращаем настроенный срок только для быстрого теста
	server.Config.WriteTimeout /= 10
	server.Start()
	defer server.Close()

	//подключаемся но читать ответ не будем
	conn,err := net.DialTimeout("tcp",server.Listener.Addr().String(),3*time.Second)
	if err != nil{t.Fatalf("подключение к серверу %v",err)}
	defer conn.Close()
	if err := conn.SetWriteDeadline(time.Now().Add(3*time.Second)); err != nil{t.Fatalf("установка срока записи: %v",err)}
	//отправляем запрос но не читаем ответ
	request := "GET / HTTP/1.1\r\nHost: localhost\r\n\r\n"
	if _,err := io.WriteString(conn,request); err != nil{t.Fatalf("отправка запроса: %v",err)}

	//ждем именно ошибку таймаута записи на сервере
	select{
		case err := <-writeErr:
			var timeoutErr net.Error
			if !errors.As(err,&timeoutErr) || !timeoutErr.Timeout(){
				t.Fatalf("ожидали таймаут записи, получили: %v",err)
			}
		case <-time.After(5*time.Second):
			t.Fatal("сервер не прекратил запись ответа")
	}
}