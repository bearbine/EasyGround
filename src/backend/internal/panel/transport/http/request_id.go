package httptransport

import (
	"context"
)

type requestIDKey struct{} //собственный тип ключа
// исключаем совпадение с ключами других пакетов

//=============================
//ретурним ид из контекста или пустую строку если его нет
func RequestID(ctx context.Context) string {
	id,_ := ctx.Value(requestIDKey{}).(string) //обез строка
	return id
}