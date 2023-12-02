create schema petshop_api

-- auto-generated definition
    create table address
    (
        id         serial       not null
            constraint petshop_api_address_pkey primary key,
        logradouro varchar(255) not null,
        numero     varchar(255) not null
    )

    create
        unique index petshop_api_address_id_uindex
        on address (id)

    create table pessoa
    (
        id serial not null
            constraint petshop_api_pessoa_pkey primary key,
        cpf_cnpj varchar(255) not null,
        tipo_pessoa varchar(255) not null
    )

    create
        unique index petshop_api_pessoa_id_uindex
        on pessoa (id)

    create table cliente
    (
        id             serial       not null
            constraint petshop_api_cliente_pkey primary key,
        nome           varchar(255) not null,
        email          varchar(255) not null,
        data_cadastro  timestamp default timezone('BRT'::text, now()),
        fk_id_address int          not null,
        fk_id_pessoa   int          not null,
        FOREIGN KEY (fk_id_address) references address (id),
        FOREIGN KEY (fk_id_pessoa) references pessoa (id)
    )

    create
        unique index petshop_api_cliente_id_uindex
        on cliente (id)

    create table telefone
    (
        id            serial       not null
            constraint petshop_api_telefone_pkey primary key,
        numero        varchar(255) not null,
        ddd           varchar(255) not null,
        tipo_telefone varchar(255) not null,
        fk_id_pessoa  int          not null,
        FOREIGN KEY (fk_id_pessoa) references pessoa (id)
    )

    create
        unique index petshop_api_telefone_id_uindex
        on telefone (id)

    create table especie
    (
        id   serial       not null
            constraint petshop_api_especie_pkey primary key,
        nome varchar(255) not null
    )

    create
        unique index petshop_api_especie_id_uindex
        on especie (id)

    create table raca
    (
        id            serial       not null
            constraint petshop_api_raca_pkey primary key,
        nome          varchar(255) not null,
        fk_id_especie int          not null,
        FOREIGN KEY (fk_id_especie) references especie (id)
    )

    create
        unique index petshop_api_raca_id_uindex
        on raca (id)

    create table pet
    (
        id              serial       not null
            constraint petshop_api_pet_pkey primary key,
        nome            varchar(255) not null,
        data_cadastro   timestamp default timezone('BRT'::text, now()),
        data_nascimento date      default timezone('BRT'::text, null),
        fk_id_cliente   int          not null,
        fk_id_raca      int          not null,
        FOREIGN KEY (fk_id_cliente) references cliente (id),
        FOREIGN KEY (fk_id_raca) references raca (id)
    )

    create
        unique index petshop_api_pet_id_uindex
        on pet (id);


-- Create default clientes inserts

INSERT INTO petshop_api.address (logradouro, numero)
VALUES ('Rua Jose Bonifácio', 1432);

INSERT INTO petshop_api.pessoa (id)
VALUES ('Rua Jose Bonifácio', 1432);

INSERT INTO petshop_api.cliente (nome, fk_id_address, email, data_cadastro, fk_id_pessoa)
VALUES ('siclano', 1, 'siclano@gmail.com', now(), 1);

INSERT INTO petshop_api.telefone (numero, ddd, tipo_telefone, fk_id_pessoa)
VALUES ('912345678', '72', 'celular', 1);

INSERT INTO petshop_api.especie (nome)
VALUES ('Canino');
INSERT INTO petshop_api.especie (nome)
VALUES ('Felino');

INSERT INTO petshop_api.raca (nome, fk_id_especie)
VALUES ('Pastor Alemao', 1);

INSERT INTO petshop_api.pet (nome, data_cadastro, data_nascimento, fk_id_cliente, fk_id_raca)
VALUES ('Rex', now(), to_date('12/12/2016', 'dd/MM/yyyy'), 1, 1);

