-- +goose Up
--разрешаем удалять записи хостов
GRANT DELETE ON TABLE public.hosts TO easyground_app;

-- +goose Down
--при откате убираем право удаления хостов
REVOKE DELETE ON TABLE public.hosts FROM easyground_app;