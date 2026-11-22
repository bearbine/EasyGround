package postgres

import (
	"os"
	"testing"
	"context"
	"time"
)
//=========================
func TestHostStoreList(t *testing.T){
	//подключение к тестовой базе
	//отдельно задается к какой именно
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == ""{t.Skip("нету TEST_DATABASE_URL")}

	//ограничиваем время подключение и запросов
	ctx,cancel :=
	context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool,err := Open(ctx,databaseURL)
	if err != nil {t.Fatalf("подключение к тестовойц базе: %v",err)}
	defer pool.Close()
	//проверка ли мы реально в тестовой базе
	//я знаю что параноидально, но бывало что я проебывался на этом моменте и подключался к основной
	var databaseName string
	if err := pool.QueryRow(ctx, "SELECT current_database()").Scan(&databaseName); err != nil {t.Fatalf("проверка нейма базы: %v", err)}
	if databaseName != "easyground_test"{t.Fatalf("для теста нужна easyground_test, а у нас %q",databaseName)}

	var hostID string
	err = pool.QueryRow(ctx, `
		INSERT INTO public.hosts (name)
		VALUES ($1)
		RETURNING id::text
	`, "integration-test-host").Scan(&hostID)
	if err != nil{t.Fatalf("создание тестового хоста: %v",err)}

	//удаляем только запись созданную ТОЛЬКО ТУТ
	defer func(){
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cleanupCancel()
		if _,err := pool.Exec(cleanupCtx, "DELETE FROM public.hosts WHERE id = $1::uuid", hostID); err != nil{t.Errorf("удаление тестового хоста:%v",err)}
	}()

	store := NewHostStore(pool)
	items, err := store.List(ctx)
	if err != nil {t.Fatalf("получение списка хостов: %v",err)}

	//ищем нашу запись
	// потом в базе могут быть другие тестовые данные
	hostIndex := -1
	for i, host := range items{
		if host.ID == hostID{
			hostIndex = i
			break
		}
	}
	if hostIndex == -1 {t.Fatalf("созданных хост %q отстуетвует в списке", hostID)}

	got := items[hostIndex]
	if got.Name !="integration-test-host"{t.Fatalf("название: получил %q, ожидал integration-test-host",got.Name)}
	if got.CreatedAt.IsZero(){t.Fatal("время создания хоста не заполнено")}
}