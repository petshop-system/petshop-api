create schema petshop_gateway

    create table gateway
    (
        id serial not null
            constraint petshop_api_gateway_pkey primary key,
        router  varchar(255) not null,
        configuration jsonb default '{"":""}'
    )

    create
        unique index petshop_api_gateway_id_uindex
        on gateway (id);

INSERT INTO petshop_gateway.gateway (router, configuration)
VALUES
    ('address', '{"host": "http://petshop-api:5001", "app-context": "petshop-api"}'),
    ('customer', '{"host": "http://petshop-api:5001", "app-context": "petshop-api"}'),
    ('employee', '{"host": "http://petshop-admin-api:5002", "app-context": "petshop-admin-api"}'),
    ('schedule', '{"host": "https://demo2908199.mockable.io", "app-context": "petshop-api"}'),
    ('schedule-request', '{"host": "http://petshop-message-api:5003","app-context": "petshop-message-api"}'),
    ('service', '{"host": "http://petshop-admin-api:5002", "app-context": "petshop-admin-api"}'),
    ('bff-mobile-customer', '{"host": "http://petshop-bff-mobile:9997", "app-context": "petshop-bff-mobile"}');

create schema petshop_api

-- auto-generated definition

    create table address
    (
        id     serial       not null
            constraint petshop_api_address_pkey primary key,
        street varchar(255) not null,
        number varchar(255) not null
    )

    create
        unique index petshop_api_address_id_uindex
        on address (id)

    create table person
    (
        id          serial       not null
            constraint petshop_api_person_pkey primary key,
        document    varchar(255) not null unique,
        person_type varchar(255) not null,
        CONSTRAINT chk_person_type_value
            CHECK (person_type IN ('individual', 'legal'))
    )

    create
        unique index petshop_api_person_id_uindex
        on person (id)

    create table contract
    (
        id            serial       not null
            constraint petshop_api_contract_pkey primary key,
        name          varchar(255) not null,
        email         varchar(255) not null unique,
        date_created  timestamp default timezone('BRT'::text, now()),
        fk_id_address int          not null unique,
        fk_id_person  int          not null unique,
        FOREIGN KEY (fk_id_address) references address (id),
        FOREIGN KEY (fk_id_person) references person (id)
    )

    create
        unique index petshop_api_contract_id_uindex
        on contract (id)

    create table customer
    (
        id             serial       not null
            constraint petshop_api_customer_pkey primary key,
        name           varchar(255) not null,
        email          varchar(255) not null unique,
        date_created   timestamp default timezone('BRT'::text, now()),
        fk_id_address  int          not null unique,
        fk_id_person   int          not null unique,
        fk_id_contract int          not null,
        FOREIGN KEY (fk_id_contract) references contract (id),
        FOREIGN KEY (fk_id_address) references address (id),
        FOREIGN KEY (fk_id_person) references person (id)
    )

    create
        unique index petshop_api_customer_id_uindex
        on customer (id)

    create table phone
    (
        id           serial       not null
            constraint petshop_api_phone_pkey primary key,
        number       varchar(255) not null,
        code_area    varchar(255) not null,
        phone_type   varchar(255) not null,
        fk_id_person int          not null,
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
        id             serial       not null
            constraint petshop_api_employee_pkey primary key,
        name           varchar(255) not null,
        register       varchar(255) not null unique,
        date_created   timestamp             default timezone('BRT'::text, now()),
        active         bool         not null default true,
        fk_id_person   int          not null,
        fk_id_contract int          not null,
        FOREIGN KEY (fk_id_contract) references contract (id),
        FOREIGN KEY (fk_id_person) references person (id)
    )

    create
        unique index petshop_api_employee_id_uindex
        on employee (id)

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
INSERT INTO petshop_api.person (document, person_type)
VALUES ('38988657000181', 'legal');

INSERT INTO petshop_api.address (street, number)
VALUES ('Rua Jose Bonifácio', 1432);

INSERT INTO petshop_api.contract (name, email, date_created, fk_id_address, fk_id_person)
VALUES ('petshop juiz de fora', 'pet_jf@gmail.com', now(), 1, 1);

INSERT INTO petshop_api.phone (number, code_area, phone_type, fk_id_person)
VALUES ('912345674', '72', 'celular', 1);

-- first customer

INSERT INTO petshop_api.person (document, person_type)
VALUES ('22233344409', 'individual');

INSERT INTO petshop_api.address (street, number)
VALUES ('Rua Lechitz', 11);

INSERT INTO petshop_api.customer (name, fk_id_address, email, date_created, fk_id_person, fk_id_contract)
VALUES ('siclano', 2, 'siclano@gmail.com', now(), 2, 1);

INSERT INTO petshop_api.phone (number, code_area, phone_type, fk_id_person)
VALUES ('912345000', '72', 'celular', 2);

-- second customer

INSERT INTO petshop_api.address (street, number)
VALUES ('Av. Juiz de Fora', 1001);

INSERT INTO petshop_api.person (document, person_type)
VALUES ('38988657000182', 'legal');

INSERT INTO petshop_api.customer (name, fk_id_address, email, date_created, fk_id_person, fk_id_contract)
VALUES ('testando cnpj', 3, 'company@gmail.com', now(), 3, 1);

INSERT INTO petshop_api.phone (number, code_area, phone_type, fk_id_person)
VALUES ('900045678', '72', 'celular', 3);

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

INSERT INTO petshop_api.person(document, person_type)
VALUES ('63609931043', 'individual');

INSERT INTO petshop_api.employee(name, register, fk_id_person, fk_id_contract)
VALUES ('Fulana da Silva Sauro', 'FUNC-0001', 4, 1);

INSERT INTO petshop_api.person(document, person_type)
VALUES ('56689159051', 'individual');

INSERT INTO petshop_api.employee(name, register, fk_id_person, fk_id_contract)
VALUES ('Ciclano da Silva Sauro', 'FUNC-0002', 5, 1);

INSERT INTO petshop_api.person(document, person_type)
VALUES ('05740847036', 'individual');

INSERT INTO petshop_api.employee(name, register, fk_id_person, fk_id_contract)
VALUES ('Brave Vacinador', 'FUNC-0003', 6, 1);

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


