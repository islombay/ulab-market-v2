do $$
    begin
    if not exists (select * from pg_type where typname = 'status_order_enum') then
    create type status_order_enum as enum (
        'in_process',
        'finished',
        'canceled'
    );
    end if;
end$$;

do $$
    begin
        if not exists (select * from pg_type where typname = 'payment_order_enum') then
            create type payment_order_enum as enum (
                'cash',
                'card'
                );
        end if;
    end$$;

create table if not exists orders (
    id uuid primary key ,
    user_id uuid,
    branch_id uuid,
    status status_order_enum default 'in_process',
    total_price numeric,
    payment_type payment_order_enum,

    created_at timestamp default now() not null,
    updated_at timestamp default now() not null,
    deleted_at timestamp default null,

    foreign key (user_id) references clients(id) on delete set null,
    foreign key (branch_id) references branches(id) on delete set null
);

create table if not exists order_products (
    id uuid primary key ,
    order_id uuid,
    product_id uuid ,
    quantity int,
    product_price numeric,
    total_price numeric,

    created_at timestamp default now() not null,
    updated_at timestamp default now() not null,
    deleted_at timestamp default null,

    foreign key (order_id) references orders(id) on delete set null,
    foreign key (product_id) references products(id) on delete  set null
);

create or replace function calculate_order_products_total_price()
returns trigger as
    $$
    begin
        new.total_price := new.product_price * new.quantity;
        return new;
    end;
$$
language plpgsql;

create trigger set_order_products_total_price
    before insert or update on order_products
    for each row
    execute function calculate_order_products_total_price();


CREATE OR REPLACE FUNCTION update_order_total() RETURNS TRIGGER AS $$
BEGIN
    -- Calculate the total price for the current order
    UPDATE orders SET total_price = (
        SELECT COALESCE(SUM(total_price), 0)
        FROM order_products
        WHERE order_id = NEW.order_id
    )
    WHERE id = NEW.order_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_order_total_after_insert_or_update
    AFTER INSERT OR UPDATE ON order_products
    FOR EACH ROW
EXECUTE FUNCTION update_order_total();

CREATE OR REPLACE FUNCTION update_order_total_delete() RETURNS TRIGGER AS $$
BEGIN
    -- Calculate the total price for the current order
    UPDATE orders SET total_price = (
        SELECT COALESCE(SUM(total_price), 0)
        FROM order_products
        WHERE order_id = OLD.order_id
    )
    WHERE id = OLD.order_id;

    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_order_total_after_delete
    AFTER DELETE ON order_products
    FOR EACH ROW
EXECUTE FUNCTION update_order_total_delete();