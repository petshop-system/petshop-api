create schema petshop_api

-- auto-generated definition
    create table address
    (
        id         serial       not null
            constraint petshop_api_address_pkey primary key,
        street varchar(255) not null,
        number     varchar(255) not null
    )

    create
        unique index petshop_api_address_id_uindex
        on address (id)

    create table person
    (
        id serial not null
            constraint petshop_api_person_pkey primary key,
        document varchar(255) not null unique ,
        person_type varchar(255) not null
    )

    create
        unique index petshop_api_person_id_uindex
        on person (id)

    create table customer
    (
        id             serial       not null
            constraint petshop_api_customer_pkey primary key,
        name           varchar(255) not null,
        email          varchar(255) not null unique ,
        date_created  timestamp default timezone('BRT'::text, now()),
        fk_id_address int          not null unique ,
        fk_id_person   int          not null unique ,
        FOREIGN KEY (fk_id_address) references address (id),
        FOREIGN KEY (fk_id_person) references person (id)
    )

    create
        unique index petshop_api_customer_id_uindex
        on customer (id)

    create table phone
    (
        id            serial       not null
            constraint petshop_api_phone_pkey primary key,
        number        varchar(255) not null,
        location           varchar(255) not null,
        phone_type varchar(255) not null,
        fk_id_person  int          not null,
        FOREIGN KEY (fk_id_person) references person (id)
    )

    create
        unique index petshop_api_phone_id_uindex
        on phone (id)

    create table species
    (
        id   serial       not null
            constraint petshop_api_species_pkey primary key,
        name varchar(255) not null unique
    )

    create
        unique index petshop_api_species_id_uindex
        on species (id)

    create table breed
    (
        id            serial       not null
            constraint petshop_api_breed_pkey primary key,
        name          varchar(255) not null,
        fk_id_species int          not null,
        FOREIGN KEY (fk_id_species) references species (id)
    )

    create
        unique index petshop_api_breed_id_uindex
        on breed (id)

    create table pet
    (
        id              serial       not null
            constraint petshop_api_pet_pkey primary key,
        name            varchar(255) not null,
        date_created   timestamp default timezone('BRT'::text, now()),
        date_birthday date      default timezone('BRT'::text, null),
        fk_id_customer   int          not null,
        fk_id_breed      int          not null,
        FOREIGN KEY (fk_id_customer) references customer (id),
        FOREIGN KEY (fk_id_breed) references breed (id)
    )

    create
        unique index petshop_api_pet_id_uindex
        on pet (id);


-- Create default clientes inserts

INSERT INTO petshop_api.address (street, number)
VALUES ('Rua Jose Bonifácio', 1432);

INSERT INTO petshop_api.address (street, number)
VALUES ('Rua Lechitz', 11);

INSERT INTO petshop_api.person (document, person_type) VALUES ('22233344409', 'individual');
INSERT INTO petshop_api.person (document, person_type) VALUES ('38988657000181', 'legal');

INSERT INTO petshop_api.customer (name, fk_id_address, email, date_created, fk_id_person)
VALUES ('siclano', 1, 'siclano@gmail.com', now(), 1);

INSERT INTO petshop_api.customer (name, fk_id_address, email, date_created, fk_id_person)
VALUES ('testando cnpj', 2, 'company@gmail.com', now(), 2);

INSERT INTO petshop_api.phone (number, location, phone_type, fk_id_person)
VALUES ('912345678', '72', 'celular', 1);

INSERT INTO petshop_api.species (name)
VALUES ('Canino');
INSERT INTO petshop_api.species (name)
VALUES ('Felino');

INSERT INTO petshop_api.breed (name, fk_id_species)
VALUES ('Pastor Alemao', 1);

INSERT INTO petshop_api.pet (name, date_created, date_birthday, fk_id_customer, fk_id_breed)
VALUES ('Rex', now(), to_date('12/12/2016', 'dd/MM/yyyy'), 1, 1);

