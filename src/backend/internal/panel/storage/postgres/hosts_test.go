package postgres

import (
	"os"
	"testing"
	"context"
	"time"
	"errors"

	"github.com/bearbine/EasyGround/src/backend/internal/panel/hosts"
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
//проверяем создание хоста в PostgreSQL
func TestHostStoreCreate(t *testing.T){
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == ""{t.Skip("нету TEST_DATABASE_URL")}

	ctx,cancel := context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()

	pool,err := Open(ctx,databaseURL)
	if err != nil {t.Fatalf("подключение к тестовой базе:%v",err)}
	defer pool.Close()

	//проверка ли мы реально в тестовой базе
	//я знаю что параноидально, но бывало что я проебывался на этом моменте и подключался к основной
	var databaseName string
	if err := pool.QueryRow(ctx, "SELECT current_database()").Scan(&databaseName); err != nil {t.Fatalf("проверка нейма базы: %v", err)}
	if databaseName != "easyground_test"{t.Fatalf("для теста нужна easyground_test, а у нас %q",databaseName)}

	//создаем хост
	store := NewHostStore(pool)
	host,err := store.Create(ctx,"  integration-test-host")
	if err != nil {t.Fatalf("создание тестового хоста: %v",err)}
	//удаляем только запись созданную ТОЛЬКО ТУТ
	defer func(){
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cleanupCancel()
		if _,err := pool.Exec(cleanupCtx, "DELETE FROM public.hosts WHERE id = $1::uuid", host.ID); err != nil{t.Errorf("удаление тестового хоста:%v",err)}
	}()

	//проверка что база заполнила ид и время создания + пробелы
	if host.ID == ""{t.Fatal("ID созданного хоста не заполнен")}
	if host.CreatedAt.IsZero(){t.Fatal("время создания хоста не заполнено")}
	if host.Name != "integration-test-host"{t.Fatalf("Name:получил %q, хотя %q",host.Name,"integration-test-host")}

	//проверка что нейм засейвился в базе без пробела
	var savedName string
	err = pool.QueryRow(ctx,`
		SELECT name FROM public.hosts WHERE id = $1::uuid
	`,host.ID).Scan(&savedName)
	if err != nil {t.Fatalf("чтение созданного хоста: %v",err)}
	if savedName != "integration-test-host"{t.Fatalf("Name в базе:получил %q,хотя %q",savedName,"integration-test-host")}
}
//проверяем получение одного хоста из PostgreSQL
func TestHostStoreGet(t *testing.T){
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == ""{t.Skip("нету TEST_DATABASE_URL")}

	ctx,cancel := context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()

	pool,err := Open(ctx,databaseURL)
	if err != nil {t.Fatalf("подключение к тестовой базе:%v",err)}
	defer pool.Close()

	//проверка ли мы реально в тестовой базе
	//я знаю что параноидально, но бывало что я проебывался на этом моменте и подключался к основной
	var databaseName string
	if err := pool.QueryRow(ctx, "SELECT current_database()").Scan(&databaseName); err != nil {t.Fatalf("проверка нейма базы: %v", err)}
	if databaseName != "easyground_test"{t.Fatalf("для теста нужна easyground_test, а у нас %q",databaseName)}

	//создаем хост
	store := NewHostStore(pool)
	host,err := store.Create(ctx,"integration-test-host")
	if err != nil {t.Fatalf("создание тестового хоста: %v",err)}

	//удаляем только запись созданную ТОЛЬКО ТУТ
	defer func(){
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cleanupCancel()
		if _,err := pool.Exec(cleanupCtx, "DELETE FROM public.hosts WHERE id = $1::uuid", host.ID); err != nil{t.Errorf("удаление тестового хоста:%v",err)}
	}()

	//получаем созданный хост по его ID
	got,err := store.Get(ctx,host.ID)
	if err != nil {t.Fatalf("получение хоста: %v",err)}

	if got.ID != host.ID{t.Fatalf("ID:получил %q, хотя %q",got.ID,host.ID)}
	if got.Name != host.Name{t.Fatalf("Name:получил %q, хотя %q",got.Name,host.Name)}
	if !got.CreatedAt.Equal(host.CreatedAt){t.Fatalf("CreatedAt:получил %v, хотя %v",got.CreatedAt,host.CreatedAt)}

	t.Run("not_found",func(t *testing.T){
		//удаляем нашу запись чтобы проверить отсутствие хоста
		_,err := pool.Exec(ctx,"DELETE FROM public.hosts WHERE id = $1::uuid",host.ID)
		if err != nil {t.Fatalf("удаление тестового хоста:%v",err)}

		_,err = store.Get(ctx,host.ID)
		if !errors.Is(err,hosts.ErrNotFound){t.Fatalf("ошибка:получил %v, хотя %v",err,hosts.ErrNotFound)}
	})
	t.Run("invalid_id",func(t *testing.T){
		_,err := store.Get(ctx,"not-a-uuid")
		if !errors.Is(err,hosts.ErrInvalidID){t.Fatalf("ошибка:получил %v, хотя %v",err,hosts.ErrInvalidID)}
	})
}
//проверяем переименование хоста в PostgreSQL
func TestHostStoreRename(t *testing.T){
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == ""{t.Skip("нету TEST_DATABASE_URL")}

	ctx,cancel := context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()

	pool,err := Open(ctx,databaseURL)
	if err != nil {t.Fatalf("подключение к тестовой базе:%v",err)}
	defer pool.Close()

	//проверка ли мы реально в тестовой базе
	//я знаю что параноидально, но бывало что я проебывался на этом моменте и подключался к основной
	var databaseName string
	if err := pool.QueryRow(ctx, "SELECT current_database()").Scan(&databaseName); err != nil {t.Fatalf("проверка нейма базы: %v", err)}
	if databaseName != "easyground_test"{t.Fatalf("для теста нужна easyground_test, а у нас %q",databaseName)}

	//создаем хост
	store := NewHostStore(pool)
	host,err := store.Create(ctx,"integration-test-host")
	if err != nil {t.Fatalf("создание тестового хоста: %v",err)}

	//удаляем только запись созданную ТОЛЬКО ТУТ
	defer func(){
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cleanupCancel()
		if _,err := pool.Exec(cleanupCtx, "DELETE FROM public.hosts WHERE id = $1::uuid", host.ID); err != nil{t.Errorf("удаление тестового хоста:%v",err)}
	}()

	//меняем название с удалением пробелов по краям
	got,err := store.Rename(ctx,host.ID,"  renamed-host  ")
	if err != nil {t.Fatalf("переименование хоста: %v",err)}

	if got.ID != host.ID{t.Fatalf("ID:получил %q, хотя %q",got.ID,host.ID)}
	if got.Name != "renamed-host"{t.Fatalf("Name:получил %q, хотя %q",got.Name,"renamed-host")}
	if !got.CreatedAt.Equal(host.CreatedAt){t.Fatalf("CreatedAt:получил %v, хотя %v",got.CreatedAt,host.CreatedAt)}

	//проверяем что новое название сохранилось в базе
	var savedName string
	err = pool.QueryRow(ctx,`
		SELECT name FROM public.hosts WHERE id = $1::uuid
	`,host.ID).Scan(&savedName)
	if err != nil {t.Fatalf("чтение переименованного хоста: %v",err)}
	if savedName != "renamed-host"{t.Fatalf("Name в базе:получил %q, хотя %q",savedName,"renamed-host")}

	t.Run("invalid_id",func(t *testing.T){
		_,err := store.Rename(ctx,"not-a-uuid","other-name")
		if !errors.Is(err,hosts.ErrInvalidID){t.Fatalf("ошибка: получил %v, хотя %v",err,hosts.ErrInvalidID)}
	})

	t.Run("invalid_name",func(t *testing.T){
		_,err := store.Rename(ctx,host.ID,"   ")
		if !errors.Is(err,hosts.ErrInvalidName){t.Fatalf("ошибка:получил %v, хотя %v",err,hosts.ErrInvalidName)}
	})
	//после неправильных запросов название должно сохраниться
	saved,err := store.Get(ctx,host.ID)
	if err != nil {t.Fatalf("получение хоста: %v",err)}
	if saved.Name != "renamed-host"{t.Fatalf("Name:получил %q, хотя %q",saved.Name,"renamed-host")}
	t.Run("not_found",func(t *testing.T){
		//удаляем нашу запись чтобы проверить отсутствие хоста
		_,err := pool.Exec(ctx,"DELETE FROM public.hosts WHERE id = $1::uuid",host.ID)
		if err != nil {t.Fatalf("удаление тестового хоста:%v",err)}

		_,err = store.Rename(ctx,host.ID,"other-name")
		if !errors.Is(err,hosts.ErrNotFound){t.Fatalf("ошибка:получил %v, хотя %v",err,hosts.ErrNotFound)}
	})
}
//проверяем удаление хоста из PostgreSQL
func TestHostStoreDelete(t *testing.T){
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == ""{t.Skip("нету TEST_DATABASE_URL")}

	ctx,cancel := context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()

	pool,err := Open(ctx,databaseURL)
	if err != nil {t.Fatalf("подключение к тестовой базе:%v",err)}
	defer pool.Close()

	//проверка ли мы реально в тестовой базе
	//я знаю что параноидально, но бывало что я проебывался на этом моменте и подключался к основной
	var databaseName string
	if err := pool.QueryRow(ctx, "SELECT current_database()").Scan(&databaseName); err != nil {t.Fatalf("проверка нейма базы: %v", err)}
	if databaseName != "easyground_test"{t.Fatalf("для теста нужна easyground_test, а у нас %q",databaseName)}

	//создаем хост
	store := NewHostStore(pool)
	host,err := store.Create(ctx,"integration-test-host")
	if err != nil {t.Fatalf("создание тестового хоста: %v",err)}

	//удаляем только запись созданную ТОЛЬКО ТУТ
	defer func(){
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cleanupCancel()
		if _,err := pool.Exec(cleanupCtx, "DELETE FROM public.hosts WHERE id = $1::uuid", host.ID); err != nil{t.Errorf("удаление тестового хоста:%v",err)}
	}()

	t.Run("invalid_id",func(t *testing.T){
		err := store.Delete(ctx,"not-a-uuid")
		if !errors.Is(err,hosts.ErrInvalidID){t.Fatalf("ошибка:получил %v, хотя %v",err,hosts.ErrInvalidID)}
	})

	//после неправильного запроса хост должен остаться
	if _,err := store.Get(ctx,host.ID); err != nil{t.Fatalf("получение хоста: %v",err)}

	//удаляем хост и проверяем что записи больше нет
	if err := store.Delete(ctx,host.ID); err != nil{t.Fatalf("удаление хоста: %v",err)}
	_,err = store.Get(ctx,host.ID)
	if !errors.Is(err,hosts.ErrNotFound){t.Fatalf("ошибка:получил %v, хотя %v",err,hosts.ErrNotFound)}

	t.Run("not_found",func(t *testing.T){
		err := store.Delete(ctx,host.ID)
		if !errors.Is(err,hosts.ErrNotFound){t.Fatalf("ошибка:получил %v, хотя %v",err,hosts.ErrNotFound)}
	})
}