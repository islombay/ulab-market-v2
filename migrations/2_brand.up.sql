create table if not exists brands (
    id uuid primary key ,
    name varchar(30) not null unique,

    created_at timestamp default now() not null,
    updated_at timestamp default now() not null,
    deleted_at timestamp default null
);