do $$
begin
    if not exists (select * from pg_type where typname = 'status_income_enum') then
create type status_income_enum as enum (
        'in_process',
        'finished',
        'canceled'
);
end if;
end$$;

create table if not exists incomes (
    id uuid primary key not null,
    branch_id uuid,
    comment varchar(255),
    courier_id uuid ,
    status status_income_enum,
    created_at timestamp default now() not null,
    updated_at timestamp default now() not null,
    deleted_at timestamp default null,

    foreign key (courier_id) references staff(id) on delete set null
    foreign key (branch_id) references branches(id) on delete set null
);

create table if not exists income_products (
    id uuid primary key not null,
    income_id uuid,
    product_id uuid,
    quantity int,
    price numeric,
    created_at timestamp default now() not null,
    updated_at timestamp default now() not null,
    deleted_at timestamp default null,
    foreign key (income_id) references incomes(id) on delete set null,
    foreign key (product_id) references products(id) on delete set null
);

