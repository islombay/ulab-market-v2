create table if not exists incomes (
    id uuid primary key not null,
    branch_id uuid,
    comment varchar(255),
    courier_id uuid ,
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

