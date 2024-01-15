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
        final_time     varchar(255) not null,
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
        booking                               date         not null,
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
VALUES ('Canino');

INSERT INTO petshop_api.species (name)
VALUES ('Felino');

INSERT INTO petshop_api.breed (name, fk_id_species)
VALUES ('Pastor Alemao', 1);

INSERT INTO petshop_api.pet (name, date_created, date_birthday, fk_id_customer, fk_id_breed, fk_id_contract)
VALUES ('Rex', now(), to_date('12/12/2016', 'dd/MM/yyyy'), 1, 1, 1);

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

INSERT INTO petshop_api.service_employee_attention_time(active, initial_time, final_time, fk_id_service, fk_id_contract, fk_id_employee)
VALUES (true, '9:00', '9:40', 1, 1, 1);

INSERT INTO petshop_api.service_employee_attention_time(active, initial_time, final_time, fk_id_service, fk_id_contract, fk_id_employee)
VALUES (true, '10:00', '10:40', 1, 1, 1);

INSERT INTO petshop_api.service_employee_attention_time(active, initial_time, final_time, fk_id_service, fk_id_contract, fk_id_employee)
VALUES (true, '11:00', '11:40', 2, 1, 1);

INSERT INTO petshop_api.service_employee_attention_time(active, initial_time, final_time, fk_id_service, fk_id_contract, fk_id_employee)
VALUES (true, '10:00', '12:40', 2, 1, 2);

INSERT INTO petshop_api.service_employee_attention_time(active, initial_time, final_time, fk_id_service, fk_id_contract, fk_id_employee)
VALUES (true, '13:00', '15:40', 2, 1, 2);

INSERT INTO petshop_api.service_employee_attention_time(active, initial_time, final_time, fk_id_service, fk_id_contract, fk_id_employee)
VALUES (true, '8:00', '8:30', 3, 1, 3);