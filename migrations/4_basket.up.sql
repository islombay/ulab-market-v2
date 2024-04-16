create table if not exists basket(
    user_id uuid not null,
    product_id uuid not null,
    quantity int default 0,

    created_at timestamp default now() not null,
    deleted_at timestamp default null,

    foreign key (user_id) references clients(id),
    foreign key (product_id) references products(id)
);