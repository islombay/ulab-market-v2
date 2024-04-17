create table if not exists products (
    id uuid primary key,
    articul varchar(20) not null ,

    name_uz varchar(50) not null default '',
    name_ru varchar(50) not null default '',

    description_uz text not null default '',
    description_ru text not null default '',

--     income_price numeric not null default 0,
    outcome_price numeric not null,

    quantity int not null default 0,

    category_id uuid default null,
    brand_id uuid default null,

    rating numeric default 0 not null ,
    status varchar(30) default 'active' not null,
    -- types of statuses:
    --  active - they are in sell, clients can see them
    --  inactive - the are not in sell, clients cannot see them even if everything is done
    --  blocked - ?
    --  archive - ?

    main_image varchar(255) default null,

    created_at timestamp default now(),
    updated_at timestamp default now(),
    deleted_at timestamp default null,

    foreign key (category_id) references category(id) on delete set null,
    foreign key (brand_id) references brands(id) on delete set null
);

create table if not exists product_image_files (
    id uuid primary key ,
    product_id uuid not null,
    media_file varchar(255) not null,

    created_at timestamp default now() not null,
    updated_at timestamp default now() not null,
    deleted_at timestamp default null,
    foreign key (product_id) references products(id) on delete cascade
);

create table if not exists product_video_files (
    id uuid primary key ,
    product_id uuid not null,
    media_file varchar(255) not null,

    created_at timestamp default now() not null,
    updated_at timestamp default now() not null,
    deleted_at timestamp default null,

    foreign key (product_id) references products(id) on delete cascade
);