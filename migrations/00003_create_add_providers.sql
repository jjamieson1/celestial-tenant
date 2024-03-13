-- +goose Up
-- SQL in this section is executed when the migration is applied.
create table IF NOT EXISTS service_provider (
    `id`                  char(36)        not null,
    `name`                varchar(255)    not null,
    `created_by`          varchar(50)     default "system" not null,
    `created_timestamp`   timestamp       default current_timestamp() not null,
    `updated_by`          char(36)        default "system" not null,
    `updated_timestamp`   timestamp       null on update current_timestamp(),
    primary key (id)
);

INSERT INTO service_provider (`id`, `name`) VALUES 
    ("0c922fb5-5e3a-4304-be76-0a37f2bfcde2", "content"),
    ("c9b1ca96-38a2-4410-81c1-421e1c76a9ad", "auth"),
    ("2f8f5e3a-991a-4665-8e01-9054303f376c", "catalog"),
    ("e41f90bc-8045-4c83-a2c2-816cacb3105c", "audit"),
    ("8a05417a-17ef-43ec-83d4-668378f12800", "messaging"),
    ("6b301e55-ff14-40c2-9776-0f5c0210e6b7", "payment");

