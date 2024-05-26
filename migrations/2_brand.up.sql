create table if not exists brands (
    id uuid primary key ,
    name varchar(30) not null,
    image varchar(255),

    created_at timestamp default now() not null,
    updated_at timestamp default now() not null,
    deleted_at timestamp default null
);
CREATE UNIQUE INDEX unique_name_brand_active ON brands (name) WHERE deleted_at IS NULL;