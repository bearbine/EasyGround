-- +goose Up
--разрешаем менять только название хоста
GRANT UPDATE (name) ON TABLE public.hosts TO easyground_app;

-- +goose Down
--при откате убираем право менять название
REVOKE UPDATE (name) ON TABLE public.hosts FROM easyground_app;