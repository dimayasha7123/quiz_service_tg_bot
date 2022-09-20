-- noinspection SqlNoDataSourceInspectionForFile

-- +goose Up
-- +goose StatementBegin
create table user_account
(
    id bigserial primary key,
    username varchar not null,
    tg_id bigint unique not null,
    qs_id bigint unique not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists user_account;
-- +goose StatementEnd
