-- +goose Up
--реестр хостов которыми управляет панель
CREATE TABLE public.hosts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    name TEXT NOT NULL
        CHECK (
            char_length(name) BETWEEN 1 AND 100
            AND btrim(name) <> ''
        ),
        
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
--бекенд пока что получает только право читать хосты
--ключевое слово ПОКА
GRANT SELECT ON TABLE public.hosts TO easyground_app;

-- +goose Down
--откат удаляет таблицу вместе со всеми данными
DROP TABLE public.hosts;
