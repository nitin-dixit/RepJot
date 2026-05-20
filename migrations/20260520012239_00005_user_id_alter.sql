-- +goose Up
-- goose StatementBegin
alter table workouts
add column user_id bigint references users(id) on delete cascade;
-- goose StatementEnd

-- +goose Down
alter table workouts drop column user_id;

