create schema petshop_api

-- auto-generated definition
    create table cliente
    (
        id            serial       not null
            constraint petshop_api_cliente_pkey primary key,
        nome          varchar(255) not null,
        telefone      json         NOT NULL,
        endereco      varchar(255) not null,
        email         varchar(255) not null,
        data_cadastro timestamp default timezone('BRT'::text, now())
    )

    create
        unique index petshop_api_cliente_id_uindex
        on cliente (id)

    create table pet
    (
        id              serial       not null
            constraint petshop_api_pet_pkey primary key,
        nome            varchar(255) not null,
        especie         varchar(255) not null,
        raca            varchar(255) not null,
        data_cadastro   timestamp default timezone('BRT'::text, now()),
        data_nascimento date default timezone('BRT'::text, null),
        fk_id_cliente   int          not null,
        FOREIGN KEY (fk_id_cliente) references cliente (id)
    )

    create
        unique index petshop_api_pet_id_uindex
        on pet (id);

-- Create default clientes inserts
INSERT INTO petshop_api.cliente (nome, telefone, endereco, email, data_cadastro)
VALUES ('siclano', '{
  "celular": "5521933235455"
}', 'Av. Rio Branco, 156', 'siclano@gmail.com', now());

INSERT INTO petshop_api.pet (nome, especie, raca, data_cadastro, data_nascimento, fk_id_cliente)
VALUES ('Bingo', 'Cachorro', 'Bordercole', now(), to_date('12/12/2016','dd/MM/yyyy'), 1);

