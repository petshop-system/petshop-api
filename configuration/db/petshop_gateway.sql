--- New Schema
create schema petshop_gateway

    create table router
    (
        id serial not null
            constraint petshop_api_gateway_pkey primary key,
        router  varchar(255) not null,
        configuration jsonb default '{"":""}'
    )

    create
        unique index petshop_api_gateway_id_uindex
        on router (id);

INSERT INTO petshop_gateway.router (router, configuration)
VALUES
    ('address', '{"host": "http://petshop-api:5001", "replace-old-app-context": "petshop-system", "replace-new-app-context": "petshop-api"}'),
    ('authentication', '{"host": "http://petshop-auth-api:5004", "replace-old-app-context": "petshop-system", "replace-new-app-context": "petshop-auth-api"}'),
    ('bff-desktop-service', '{"host": "http://petshop-bff-desktop:9998", "replace-old-app-context": "petshop-system/bff-desktop-service", "replace-new-app-context": "petshop-bff-desktop"}'),
    ('customer', '{"host": "http://petshop-api:5001", "replace-old-app-context": "petshop-system", "replace-new-app-context": "petshop-api"}'),
    ('employee', '{"host": "http://petshop-admin-api:5002", "replace-old-app-context": "petshop-system", "replace-new-app-context": "petshop-admin-api"}'),
    ('schedule', '{"host": "https://demo2908199.mockable.io", "replace-old-app-context": "petshop-system", "replace-new-app-context": "petshop-api"}'),
    ('schedule-request', '{"host": "http://petshop-message-api:5003", "replace-old-app-context": "petshop-system", "replace-new-app-context": "petshop-message-api"}'),
    ('service', '{"host": "http://petshop-admin-api:5002", "replace-old-app-context": "petshop-system", "replace-new-app-context": "petshop-admin-api"}');

