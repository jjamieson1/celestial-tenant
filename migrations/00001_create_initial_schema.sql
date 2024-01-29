-- +goose Up
-- SQL in this section is executed when the migration is applied.
create table IF NOT EXISTS tenant (
    tenant_id           char(36)        not null,
    parent_tenant_id    char(36)        not null,
    url                 varchar(255)    not null,
    common_name         varchar(255)    not null,
    logo_primary_url            varchar(255)    null,
    is_available        int             DEFAULT 1 not null,
    status_id           int             DEFAULT 2 not null,
    created_by          varchar(50)     not null,
    created_timestamp   timestamp       default current_timestamp() not null,
    updated_by          char(36)        null,
    updated_timestamp   timestamp       null on update current_timestamp(),
    primary key (tenant_id)
);

create table IF NOT EXISTS tenant_type (
    tenant_type_id          char(36)          not null,
    tenant_type_name        varchar(100) not null,
    primary key (tenant_type_id)
);

create table IF NOT EXISTS tenant_type_to_tenant (
    id int not null AUTO_INCREMENT PRIMARY KEY,
    tenant_type_id char(36) NOT NULL,
    tenant_id char(36) NOT NULL
);

INSERT INTO tenant_type (tenant_type_id, tenant_type_name) VALUES ('65c5a9b2-4c23-47df-8a10-1b563e9cf4c8','Basic Tenant');
INSERT INTO tenant_type (tenant_type_id, tenant_type_name) VALUES ('65c5a9b2-4c23-47df-8a10-1b563f9cf4c9','Extended Tenant');


create table  IF NOT EXISTS status
(
    status_id         int                                   not null,
    description       varchar(50)                           null,
    created_by        varchar(50)                           not null,
    created_timestamp timestamp default current_timestamp() not null,
    updated_by        varchar(50)                           null,
    updated_timestamp timestamp                             null on update current_timestamp(),
    primary key (status_id)
);

INSERT INTO status (status_id, description, created_by, created_timestamp, updated_by, updated_timestamp) VALUES (1, 'DISABLED', 'SYSTEM', now(), null, null);
INSERT INTO status (status_id, description, created_by, created_timestamp, updated_by, updated_timestamp) VALUES (2, 'PENDING', 'SYSTEM', now(), null, null);
INSERT INTO status (status_id, description, created_by, created_timestamp, updated_by, updated_timestamp) VALUES (3, 'ACTIVE', 'SYSTEM', now(), null, null);
INSERT INTO status (status_id, description, created_by, created_timestamp, updated_by, updated_timestamp) VALUES (4, 'ARCHIVED', 'SYSTEM', now(), null, null);
INSERT INTO status (status_id, description, created_by) VALUES (5, 'VERIFIED', 'SYSTEM');

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

