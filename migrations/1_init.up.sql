create table if not exists roles (
    id uuid primary key,
    name varchar(40) not null unique,
    description varchar(255),
    created_at timestamp default now() not null,
    updated_at timestamp default now() not null,
    deleted_at timestamp default null
);

create table if not exists permissions (
    id uuid primary key,
    name varchar(40) not null unique,
    description varchar(255),

    created_at timestamp default now() not null,
    updated_at timestamp default now() not null,
    deleted_at timestamp default null
);

create table if not exists permission_to_role (
    role_id uuid not null,
    permission_id uuid not null,

    created_at timestamp default now() not null,
    updated_at timestamp default now() not null,
    deleted_at timestamp default null,

    foreign key (role_id) references roles(id) on delete cascade,
    foreign key (permission_id) references permissions(id) on delete cascade
);

create table if not exists clients (
    id uuid primary key,
    name varchar(30) default null,
    phone_number varchar(12) unique,
    email varchar(255) unique,

    created_at timestamp default now() not null,
    updated_at timestamp default now() not null,
    deleted_at timestamp default null
);

create table if not exists staff (
    id uuid primary key,
    name varchar(30) not null,
    phone_number varchar(12),
    email varchar(255),
    password varchar(255),

    role_id uuid not null ,

    created_at timestamp default now() not null,
    updated_at timestamp default now() not null,
    deleted_at timestamp default null,

    foreign key(role_id) references roles(id) on delete set null
);

create table if not exists client_address (
    id uuid primary key,
    client_id uuid not null,
    long float not null,
    lat float not null,
    name varchar(30),

    created_at timestamp default now() not null,
    updated_at timestamp default now() not null,
    deleted_at timestamp default null,

    foreign key (client_id) references clients (id) on delete cascade
);

create table if not exists category (
    id uuid primary key,
    name_uz varchar(255) unique not null,
    name_ru varchar(255) unique not null,
    image varchar(255),
    parent_id uuid,

    created_at timestamp default now() not null,
    updated_at timestamp default now() not null,
    deleted_at timestamp default null
);