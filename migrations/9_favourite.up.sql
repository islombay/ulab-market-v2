create table if not exists favourite (
    user_id uuid not null,
    product_id uuid not null,

    foreign key (user_id) references clients(id) on delete cascade ,
    foreign key (product_id) references products(id) on delete cascade
);