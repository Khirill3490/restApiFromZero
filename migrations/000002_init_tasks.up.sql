CREATE TABLE tasks (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Уникальный title В РАМКАХ одного пользователя
CREATE UNIQUE INDEX tasks_user_title_unique ON tasks(user_id, title);

-- Для быстрых выборок задач пользователя
CREATE INDEX tasks_user_id_idx ON tasks(user_id);
