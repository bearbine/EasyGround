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
	items,err := store.List(ctx,hosts.DefaultListLimit,0)
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
//проверяем постраничное получение хостов из PostgreSQL
func TestHostStoreListPagination(t *testing.T){
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
	//учитываем существующие записи чтобы не менять чужие данные
	var existingCount int
	var latestCreatedAt time.Time
	err = pool.QueryRow(ctx,`
		SELECT count(*), COALESCE(max(created_at),now())
		FROM public.hosts
	`).Scan(&existingCount,&latestCreatedAt)
	if err != nil {t.Fatalf("получение исходного состояния: %v",err)}

	//сохраняем ID только созданных этим тестом записей
	var hostIDs []string

	//удаляем только записи созданные ТОЛЬКО ТУТ
	defer func(){
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(),5*time.Second)
		defer cleanupCancel()

		for _,id := range hostIDs{
			if _,err := pool.Exec(cleanupCtx,"DELETE FROM public.hosts WHERE id = $1::uuid",id); err != nil{t.Errorf("удаление тестового хоста:%v",err)}
		}
	}()
	//размещаем три тестовых хоста после существующих записей
	for i := 0; i < 3; i++{
		createdAt := latestCreatedAt.Add(time.Duration(i+1)*time.Second)
		var id string

		err = pool.QueryRow(ctx,`
			INSERT INTO public.hosts (name,created_at)
			VALUES ($1,$2)
			RETURNING id::text
		`,"pagination-test-host",createdAt).Scan(&id)
		if err != nil {t.Fatalf("создание тестового хоста: %v",err)}

		hostIDs = append(hostIDs,id)
	}
	store := NewHostStore(pool)

	//пропускаем старые записи и получаем первые два тестовых хоста
	firstPage,err := store.List(ctx,2,existingCount)
	if err != nil {t.Fatalf("получение первой страницы: %v",err)}
	if len(firstPage) != 2{t.Fatalf("размер первой страницы:получил %d, хотя 2",len(firstPage))}
	for i,host := range firstPage{
		if host.ID != hostIDs[i]{t.Fatalf("ID на позиции %d:получил %q, хотя %q",i,host.ID,hostIDs[i])}
	}
	//на второй странице остается только третий тестовый хост
	secondPage,err := store.List(ctx,2,existingCount+2)
	if err != nil {t.Fatalf("получение второй страницы: %v",err)}
	if len(secondPage) != 1{t.Fatalf("размер второй страницы:получил %d, хотя 1",len(secondPage))}
	if secondPage[0].ID != hostIDs[2]{t.Fatalf("ID:получил %q, хотя %q",secondPage[0].ID,hostIDs[2])}
	//после последней записи должна вернуться пустая страница
	emptyPage,err := store.List(ctx,2,existingCount+3)
	if err != nil {t.Fatalf("получение пустой страницы: %v",err)}
	if len(emptyPage) != 0{t.Fatalf("размер пустой страницы:получил %d, хотя 0",len(emptyPage))}
}