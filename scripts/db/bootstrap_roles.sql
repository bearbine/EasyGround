-- подготовка роли перед миграциями
-- права на таблицы остаются в миграциях

-- при ошибке не забываем остановить выполнение скрипта
\set ON_ERROR_STOP on

--пароль получаем из окружение процесса psql
\getenv app_password EASYGROUND_APP_PASSWORD
\if :{?app_password}
\else
DO $$
BEGIN
    RAISE EXCEPTION 'переменная окружения EASYGROUND_APP_PASSWORD не задана';
END
$$;
\endif

BEGIN;
SET LOCAL log_statement = 'none';
SET LOCAL log_min_error_statement = 'panic';
SET LOCAL log_min_duration_statement = -1;
SET LOCAL log_min_duration_sample = -1;
SET LOCAL log_transaction_sample_rate = 0;

SELECT length(:'app_password') >= 16 as app_password_ok --если что
-- то тут пустой или слишком короткий не будет идти
\gset

\if :app_password_ok
\else
DO $$
BEGIN 
    RAISE EXCEPTION 'EASYGROUND_APP_PASSWORD должен быть 16 символов или больше';
END
$$;
\endif

SELECT 'CREATE ROLE easyground_app'-- если ее нет
WHERE NOT EXISTS (
    SELECT 1
    FROM pg_catalog.pg_roles
    WHERE rolname = 'easyground_app'
)
\gexec

ALTER ROLE easyground_app WITH
    LOGIN
    NOSUPERUSER
    NOCREATEDB
    NOCREATEROLE
    NOREPLICATION
    NOBYPASSRLS
    PASSWORD :'app_password';

--подключение к базе разрешаем явно
--бтв нужен доступ к схеме
REVOKE CONNECT ON DATABASE :"DBNAME" FROM PUBLIC;
GRANT CONNECT ON DATABASE :"DBNAME" TO easyground_app;
REVOKE CREATE ON SCHEMA public FROM PUBLIC;
GRANT USAGE ON SCHEMA public TO easyground_app;

COMMIT;
\unset app_password