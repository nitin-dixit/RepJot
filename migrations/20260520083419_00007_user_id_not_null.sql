-- +goose Up
-- goose StatementBegin
alter table workouts
alter column user_id set not null;
-- goose StatementEnd
