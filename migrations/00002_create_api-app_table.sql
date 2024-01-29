-- +goose Up
-- SQL in this section is executed when the migration is applied.
create table IF NOT EXISTS secret_keys (
    tenant_id           char(36)        not null,
    api_key             char(36)        not null
);