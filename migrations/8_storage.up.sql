create table if not exists storage (
    id uuid primary key not null,
    product_id uuid ,
    income_id uuid,
    branch_id uuid,
    total_price numeric,
    quantity int,

    created_at timestamp default now() not null,
    updated_at timestamp default now() not null,
    deleted_at timestamp default null,

    foreign key (product_id) references products(id) on delete  set null,
    foreign key (income_id) references incomes(id) on delete set null,
    foreign key (branch_id) references branches(id) on delete set null
);