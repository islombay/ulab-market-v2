create table if not exists branches(
    id uuid primary key ,
    name varchar(100) not null ,

    open_time time not null,
    close_time time not null,

    created_at timestamp default now() not null ,
    updated_at timestamp default now() not null,
    deleted_at timestamp default null
);