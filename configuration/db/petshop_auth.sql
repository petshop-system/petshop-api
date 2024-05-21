create schema petshop_auth

    create table profile
    (
        name        varchar(100) not null unique primary key,
        description varchar(200)
    )

    create
        unique index petshop_auth_api_profile_name_uindex
        on profile (name)

    create table access
    (
        action      varchar(100) not null unique primary key,
        description varchar(100)
    )

    create
        unique index petshop_auth_api_access_action_uindex
        on access (action)

    create table profile_access
    (
        fk_profile varchar not null,
        fk_access  varchar not null,
        FOREIGN KEY (fk_profile) references profile (name),
        FOREIGN KEY (fk_access) references access (action),
        constraint petshop_auth_api_profile_access_pkey PRIMARY KEY (fk_profile, fk_access)
    )

    create
        unique index petshop_auth_api_profile_access_uindex
        on profile_access (fk_profile, fk_access)

    create table authentication
    (
        id         serial       not null
            constraint petshop_auth_api_authentication_pkey primary key,
        login      varchar(100) not null unique,
        password   varchar(255),
        id_user    int,
        active     bool         not null default false,
        fk_profile varchar      not null default '',
        FOREIGN KEY (fk_profile) references profile (name)
    )

    create
        unique index petshop_auth_api_authentication_id_uindex
        on authentication (id)

    create
        unique index petshop_auth_api_id_authentication_user_uindex
        on authentication (id_user)

    create
        unique index petshop_auth_api_authentication_login_uindex
        on authentication (login)

    create table access_token
    (
        token                varchar(200) not null unique primary key,
        fk_id_authentication int          not null,
        date_created         timestamp default timezone('BRT'::text, now()),
        FOREIGN KEY (fk_id_authentication) REFERENCES authentication (id)
    )

    create
        unique index petshop_auth_api_access_token_token_uindex
        on access_token (token);

INSERT INTO petshop_auth.profile(name, description)
VALUES ('ADMINISTRATOR', 'system administrator'),
       ('API', 'api request'),
       ('CUSTOMER', 'costumer user'),
       ('EMPLOYEE', 'employee user'),
       ('MANAGER', 'manager user');

INSERT INTO petshop_auth.access(action, description)
VALUES ('CUSTOMER_CREATE', 'access to create a new customer'),
       ('CUSTOMER_UPDATE', 'access to update a known customer'),
       ('EMPLOYEE_CREATE', 'access to create a ner employee'),
       ('EMPLOYEE_UPDATE', 'access to update a known employee');

INSERT INTO petshop_auth.profile_access(fk_profile, fk_access)
VALUES ('ADMINISTRATOR', 'CUSTOMER_CREATE'),
       ('ADMINISTRATOR', 'CUSTOMER_UPDATE'),
       ('ADMINISTRATOR', 'EMPLOYEE_CREATE'),
       ('ADMINISTRATOR', 'EMPLOYEE_UPDATE'),
       ('API', 'CUSTOMER_CREATE'),
       ('API', 'CUSTOMER_UPDATE'),
       ('CUSTOMER', 'CUSTOMER_CREATE'),
       ('CUSTOMER', 'CUSTOMER_UPDATE'),
       ('EMPLOYEE', 'CUSTOMER_CREATE'),
       ('EMPLOYEE', 'CUSTOMER_UPDATE'),
       ('MANAGER', 'CUSTOMER_CREATE'),
       ('MANAGER', 'CUSTOMER_UPDATE'),
       ('MANAGER', 'EMPLOYEE_CREATE'),
       ('MANAGER', 'EMPLOYEE_UPDATE');

INSERT INTO petshop_auth.authentication (login, password, id_user, fk_profile)
VALUES ('admin@petshopsystem.com', '$2a$10$InNNO0QNeLe3Yc2mWPYmvOta421VA64e1Hq1mpPGvgymdm6w0uBvq', 1,
        'ADMINISTRATOR'); -- senha 1234
