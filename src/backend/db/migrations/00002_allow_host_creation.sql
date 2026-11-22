-- +goose Up
--разрешаем задавать название при создании хоста
GRANT INSERT (name) ON TABLE public.hosts TO easyground_app;

-- +goose Down
-- при откате убираем право добавлять хосты
REVOKE INSERT (name) ON TABLE public.hosts FROM easyground_app;