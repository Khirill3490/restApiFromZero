INSERT INTO tasks (title, description, completed)
VALUES
 ('Настроить проект', 'Инициализация Go проекта', true),
 ('Подключить Postgres', 'docker-compose + postgres', true),
 ('Добавить миграции', 'golang-migrate', true),
 ('Создать таблицу tasks', 'Основная таблица задач', true),
 ('Seed данные', 'Начальные данные в БД', true),
 ('Добавить API', 'CRUD для задач', false),
 ('Валидация', 'Проверка входящих данных', false),
 ('JWT авторизация', 'Access + Refresh tokens', false),
 ('Логирование', 'zap / slog', false),
 ('Docker build', 'Сборка Go сервиса', false)
ON CONFLICT (title) DO NOTHING;
