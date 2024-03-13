-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE IF NOT EXISTS service_provider_tenant_data (
    `id`                    varchar(36) not null, 
    `tenant_id`             varchar(36) not null,
    `service_provider_id`   varchar(36) not null,
    `URL`                   text,
    `created_by`            varchar(50)     default "system" not null,
    `created_timestamp`     timestamp       default current_timestamp() not null,
    `updated_by`            char(36)        default "system" not null,
    `updated_timestamp`     timestamp       null on update current_timestamp(),
    primary key (id)
);

INSERT INTO service_provider_tenant_data (`id`, `tenant_id`, `service_provider_id`, `URL`) VALUES 
("4d08113d-35aa-47ce-aba8-e13aeb3a02d8", "95086244-7674-11ed-ab0f-ee4d31332ff3", "0c922fb5-5e3a-4304-be76-0a37f2bfcde2", "http://localhost:3001"),
("4f393c79-b25c-4782-a701-e0830d806dee", "95086244-7674-11ed-ab0f-ee4d31332ff3", "c9b1ca96-38a2-4410-81c1-421e1c76a9ad", "http://localhost:3000"),
("fac2e2db-1d6c-4601-9710-5030716eb5a7", "95086244-7674-11ed-ab0f-ee4d31332ff3", "2f8f5e3a-991a-4665-8e01-9054303f376c", "http://localhost:3004"),
("09d5b876-576b-411b-a25f-b1acda629e76", "95086244-7674-11ed-ab0f-ee4d31332ff3", "e41f90bc-8045-4c83-a2c2-816cacb3105c", "http://localhost:3003"),
("e163b6c3-67e1-432a-b3e3-854c8b6f8ea8", "95086244-7674-11ed-ab0f-ee4d31332ff3", "8a05417a-17ef-43ec-83d4-668378f12800", "http://localhost:3005"),
("ce173bb0-6b0e-4e45-879f-4dfc543a1522", "95086244-7674-11ed-ab0f-ee4d31332ff3", "6b301e55-ff14-40c2-9776-0f5c0210e6b7", "http://localhost:3006"),
("b3806ee3-5d68-48f5-97a5-a204c47e8b34", "95086244-7674-11ed-ab0f-ee4d31332ff3", "711fc83a-54f9-4335-9ae6-f097e0b724fe", "https://files.celestialtech.ca");


INSERT INTO service_provider (`id`, `name`) VALUES 
    ("711fc83a-54f9-4335-9ae6-f097e0b724fe", "file");