-- +goose Up
--  goose StatementBegin
update workouts
set user_id=1
where user_id is null;
-- goose StatementEnd
-- +goose Down
