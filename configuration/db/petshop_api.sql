--- New Schema
create schema petshop_api

-- auto-generated definition

    create table address
    (
        id           SERIAL       NOT NULL
            CONSTRAINT petshop_api_address_pkey PRIMARY KEY,
        street       VARCHAR(255) NOT NULL,
        number       VARCHAR(255) NOT NULL,
        complement    VARCHAR(10),
        block        VARCHAR(10),
        neighborhood VARCHAR(255) NOT NULL,
        zip_code      VARCHAR(20) NOT NULL,
        city         VARCHAR(255) NOT NULL,
        state        VARCHAR(2) NOT NULL,
        country      VARCHAR(100) NOT NULL
    )

    create
        unique index petshop_api_address_id_uindex
        on address (id)

    create table contract
    (
        id            serial       not null
            constraint petshop_api_contract_pkey primary key,
        name          varchar(255) not null,
        email         varchar(255) not null unique,
        document      varchar(255) not null unique,
        person_type   varchar(255) not null,
        date_created  timestamp default timezone('BRT'::text, now()),
        fk_id_address int          not null unique,
        FOREIGN KEY (fk_id_address) references address (id),
        CONSTRAINT chk_person_type_value
            CHECK (person_type IN ('individual', 'legal'))
    )

    create
        unique index petshop_api_contract_id_uindex
        on contract (id)

    create
        unique index petshop_api_contract_document_uindex
        on contract (document)

    create table customer
    (
        id             serial       not null unique ,
        name           varchar(255) not null,
        email          varchar(255) not null,
        document       varchar(255) not null,
        person_type    varchar(255) not null,
        fk_id_address  int          not null unique,
        fk_id_contract int          not null,
        FOREIGN KEY (fk_id_contract) references contract (id),
        FOREIGN KEY (fk_id_address) references address (id),
        constraint petshop_api_customer_pkey PRIMARY KEY (id, email, document, fk_id_contract),
        CONSTRAINT chk_customer_person_type_value
            CHECK (person_type IN ('individual', 'legal'))
    )

    create
        unique index petshop_api_customer_id_uindex
        on customer (id)

    create
        index petshop_api_customer_email_uindex
        on customer (email)

    create
        index petshop_api_customer_document_uindex
        on customer (document)

    create table customer_history
    (
        id serial not null
            constraint petshop_api_customer_history_pkey primary key,
        fk_id_customer int          not null,
        date   timestamp default timezone('BRT'::text, now()),
        description varchar(255) not null,
        FOREIGN KEY (fk_id_customer) references customer (id)
    )

    create
        index petshop_api_customer_history_fk_id_customer_uindex
        on customer_history (fk_id_customer)

    create table phone
    (
        id           serial       not null
            constraint petshop_api_phone_pkey primary key,
        number       varchar(255) not null,
        code_area    varchar(255) not null,
        phone_type   varchar(255) not null
    )

    create
        unique index petshop_api_phone_id_uindex
        on phone (id)

    create table phone_user
    (
        id          serial       not null
            constraint petshop_api_phone_user_pkey primary key,
        fk_id_phone int          not null unique,
        fk_id_user  int          not null,
        user_type   varchar(255) not null,
        FOREIGN KEY (fk_id_phone) references phone (id),
        CONSTRAINT chk_phone_user_type_value
            CHECK (user_type IN ('contract', 'customer', 'employee'))
    )

    create
        index petshop_api_phone_user_uindex
        on phone_user (fk_id_user)

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
        id             serial       not null
            constraint petshop_api_pet_pkey primary key,
        name           varchar(255) not null,
        date_created   timestamp default timezone('BRT'::text, now()),
        date_birthday  date      default timezone('BRT'::text, null),
        date_deleted   date      default timezone('BRT'::text, null),
        fk_id_customer int          not null,
        fk_id_breed    int          not null,
        fk_id_contract int          not null,
        FOREIGN KEY (fk_id_contract) references contract (id),
        FOREIGN KEY (fk_id_customer) references customer (id),
        FOREIGN KEY (fk_id_breed) references breed (id)
    )

    create
        unique index petshop_api_pet_id_uindex
        on pet (id)

    create
        index petshop_api_pet_fk_id_customer_uindex
        on pet (fk_id_customer)

    create table service
    (
        id             serial       not null
            constraint petshop_api_service_pkey primary key,
        name           varchar(255) not null,
        price          decimal      not null default 0,
        active         bool         not null default true,
        description    varchar(255) not null,
        fk_id_contract int          not null,
        FOREIGN KEY (fk_id_contract) references contract (id)
    )

    create
        unique index petshop_api_service_id_uindex
        on service (id)

    create table employee
    (
        id             serial       not null unique,
        name           varchar(255) not null,
        register       varchar(255) not null unique,
        document       varchar(255) not null unique,
        date_created   timestamp             default timezone('BRT'::text, now()),
        active         bool         not null default true,
        fk_id_contract int          not null,
        FOREIGN KEY (fk_id_contract) references contract (id),
        constraint petshop_api_employee_pkey PRIMARY KEY (id, document, fk_id_contract)
    )

    create
        unique index petshop_api_employee_id_uindex
        on employee (id)

    create
        unique index petshop_api_employee_register_uindex
        on employee (register)

    create
        unique index petshop_api_employee_document_uindex
        on employee (document)

    create table service_employee_attention_time
    (
        id             serial       not null
            constraint petshop_api_service_employee_attention_time_pkey primary key,
        initial_time   varchar(255) not null,
--         final_time     varchar(255) not null,
        active         bool         not null default false,
        fk_id_service  int          not null,
        fk_id_contract int          not null,
        fk_id_employee int          not null,
        FOREIGN KEY (fk_id_contract) references contract (id),
        FOREIGN KEY (fk_id_service) references service (id),
        FOREIGN KEY (fk_id_employee) references employee (id)
    )

    create
        unique index petshop_api_service_employee_attention_time_id_uindex
        on service_employee_attention_time (id)

    create table schedule
    (
        id                                    serial       not null
            constraint petshop_api_schedule_pkey primary key,
        date_created                          timestamp             default timezone('BRT'::text, now()),
        date_declined                         timestamp,
        number                                varchar(255) not null, -- 2023dez10.000001
        booked_at                             date         not null,
        price                                 decimal      not null default 0,
        fk_id_pet                             int          not null,
        fk_id_service_employee_attention_time int          not null,
--         fk_id_contract                        int          not null,
--         FOREIGN KEY (fk_id_contract) references contract (id),
        FOREIGN KEY (fk_id_pet) references pet (id),
        FOREIGN KEY (fk_id_service_employee_attention_time) references service_employee_attention_time (id)
    )

    create
        unique index petshop_api_schedule_id_uindex
        on schedule (id);


-- Create default inserts

-- contract
INSERT INTO petshop_api.address (street, number, complement, block, neighborhood, zip_code, city, state, country)
VALUES ('Rua Jose Bonifácio', 1432, 403, 'B', 'Centro', '36025-200', 'Juiz de Fora', 'MG', 'Brasil');

-- INSERT INTO petshop_api.contract (name, email, date_created, fk_id_address, document, person_type)
-- VALUES ('petshop juiz de fora', 'pet_jf@gmail.com', now(), 1, '38988657000181', 'legal');

INSERT INTO petshop_api.phone (number, code_area, phone_type)
VALUES ('912345674', '72', 'celular');

INSERT INTO petshop_api.phone_user(fk_id_phone, fk_id_user, user_type)
VALUES (1, 1, 'contract');

-- first customer
INSERT INTO petshop_api.address (street, number, complement, block, neighborhood, zip_code, city, state, country)
VALUES ('Rua Lechitz', 11, 201, 'A', 'São Mateus', '36025-290', 'Juiz de Fora', 'MG', 'Brasil');

INSERT INTO petshop_api.customer (name, fk_id_address, email,  fk_id_contract, document, person_type)
VALUES ('siclano', 2, 'siclano@gmail.com', 1, '22233344409', 'individual');

INSERT INTO petshop_api.phone (number, code_area, phone_type)
VALUES ('912345000', '72', 'celular');

INSERT INTO petshop_api.phone_user(fk_id_phone, fk_id_user, user_type)
VALUES (2, 1, 'customer');

-- second customer

INSERT INTO petshop_api.address (street, number, complement, block, neighborhood, zip_code, city, state, country)
VALUES ('Av. Juiz de Fora', 1001, null, null, 'Centro', '36025-100', 'Juiz de Fora', 'MG', 'Brasil')

INSERT INTO petshop_api.customer (name, fk_id_address, email, fk_id_contract, document, person_type)
VALUES ('testando cnpj', 3, 'company@gmail.com', 1, '38988657000182', 'legal');

INSERT INTO petshop_api.phone (number, code_area, phone_type)
VALUES ('900045678', '72', 'celular');

INSERT INTO petshop_api.phone_user(fk_id_phone, fk_id_user, user_type)
VALUES (3, 2, 'customer');

-- pet control

INSERT INTO petshop_api.species (name)
VALUES ('Canino'), ('Felino');

INSERT INTO petshop_api.breed (name, fk_id_species)
VALUES ('Pastor Alemao', 1), ('Siames', 2);

INSERT INTO petshop_api.pet (name, date_created, date_birthday, fk_id_customer, fk_id_breed, fk_id_contract)
VALUES ('Rex', now(), to_date('12/12/2016', 'dd/MM/yyyy'), 1, 1, 1),
       ('Rex', now(), to_date('12/09/2023', 'dd/MM/yyyy'), 1, 2, 1);

INSERT INTO petshop_api.service (name, price, active, fk_id_contract, description)
VALUES ('TOSA', 50.65, true, 1, 'Tosa com tesoura.');

INSERT INTO petshop_api.service (name, price, active, fk_id_contract, description)
VALUES ('BANHO', 55.99, true, 1, 'Banho com sais minerais e água morna.');

INSERT INTO petshop_api.service (name, price, active, fk_id_contract, description)
VALUES ('VACINA ANTIRRABICA', 112.70, true, 1, 'Vacina antirrabica para cachorros.');

-- Employee

INSERT INTO petshop_api.employee(name, register, fk_id_contract, document)
VALUES ('Fulana da Silva Sauro', 'FUNC-0001', 1, '63609931043');

INSERT INTO petshop_api.employee(name, register, fk_id_contract, document)
VALUES ('Ciclano da Silva Sauro', 'FUNC-0002', 1, '56689159051');

INSERT INTO petshop_api.employee(name, register, fk_id_contract, document)
VALUES ('Brave Vacinador', 'FUNC-0003', 1, '05740847036');

-- service employee attention time

INSERT INTO petshop_api.service_employee_attention_time(active, initial_time,  fk_id_service, fk_id_contract, fk_id_employee)
VALUES (true, '9:00', 1, 1, 1),
       (true, '9:00', 2, 1, 1),
       (true, '10:00', 1, 1, 1),
       (true, '11:00', 2, 1, 1),
       (true, '10:00', 2, 1, 2),
       (true, '13:00', 2, 1, 2),
       (false, '8:00', 3, 1, 3);

INSERT INTO petshop_api.schedule(date_created, number, booked_at, price, fk_id_pet, fk_id_service_employee_attention_time)
VALUES (now(), '2024020001', now() + interval '1 day', 10.50, 1, 1),
       (now(), '2024020002', now() + interval '1 day', 100.50, 2, 4);


-- FUNCTIONS

CREATE OR REPLACE FUNCTION petshop_api.GET_SERVICE_ATTENTION_AVAILABLE(P_DATE_SCHEDULE VARCHAR(12), P_SERVICE_ID INTEGER)
    RETURNS TABLE
            (
                service_attention_id int,
                service_attention_active bool,
                service_attention_initial_time varchar(255),
                service_attention_fk_id_service int,
                service_attention_fk_id_contract int,
                service_attention_fk_id_employee int
            )
    LANGUAGE plpgsql
AS
$$
BEGIN
    RETURN QUERY
        select
            service_attention.id::integer,
            service_attention.active,
            service_attention.initial_time,
            service_attention.fk_id_service,
            service_attention.fk_id_contract,
            service_attention.fk_id_employee
        from (select service_attention.*
              from petshop_api.service_employee_attention_time service_attention
              where 1 = 1
                and service_attention.active = true
                and not exists (select 1
                                from petshop_api.schedule schedule
                                where service_attention.id = schedule.fk_id_service_employee_attention_time
                                  and schedule.booked_at = TO_DATE(P_DATE_SCHEDULE, 'YYYY-MM-DD')))
                 service_attention -- getting all available services attention in order to schedules
        where 1 = 1
          and not exists( -- denying select
            -- select to get employees with possibles appointments in the same hour
            select 1 from petshop_api.schedule schedule
                              inner join petshop_api.service_employee_attention_time seat
                                         on schedule.fk_id_service_employee_attention_time = seat.id
            where 1 = 1
              and schedule.booked_at = TO_DATE(P_DATE_SCHEDULE, 'YYYY-MM-DD')
              and service_attention.initial_time = seat.initial_time
              and seat.fk_id_employee = service_attention.fk_id_employee

        )
          and service_attention.fk_id_service = P_SERVICE_ID -- getting only service attention for specific service
        order by cast(SPLIT_PART(initial_time, ':', 1) as INTEGER);
end;
$$