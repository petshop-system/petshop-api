create schema petshop_api

-- auto-generated definition
    create table cliente
    (
        id            serial       not null
            constraint petshop_api_cliente_pkey primary key,
        nome          varchar(255) not null,
        telefone      json         NOT NULL,
        endereco      varchar(255) not null,
        data_cadastro timestamp default timezone('BRT'::text, now())
    )

    create
        unique index petshop_api_cliente_id_uindex
        on cliente (id);

-- Create default accounts inserts
INSERT INTO petshop_api.cliente (nome, telefone, endereco, data_cadastro)
VALUES ('siclano', '{"celular": "5521933235455"}', 'Av. Rio Branco, 156', now());