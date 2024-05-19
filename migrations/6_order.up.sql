do $$
    begin
    if not exists (select * from pg_type where typname = 'status_order_enum') then
    create type status_order_enum as enum (
        'in_process',
        'finished',
        'canceled',
        'picking',
        'picked',
        'delivering'
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

do $$
    begin 
        if not exists (select * from pg_type where typname = 'delivery_order_enum') then
            create type delivery_order_enum as enum (
                'deliver'
            );
        end if;
end $$;

-- Create a function to generate the order ID
CREATE OR REPLACE FUNCTION generate_order_id() RETURNS TEXT AS $$
DECLARE
    unix_time BIGINT;
    random_letters TEXT;
BEGIN
    -- Get current Unix time
    unix_time := EXTRACT(EPOCH FROM NOW())::BIGINT;
    
    -- Remove first 3 digits
    unix_time := unix_time % 1000000000;
    
    -- Generate random letters
    random_letters := CHR(ASCII('a') + floor(random() * 26)) || CHR(ASCII('a') + floor(random() * 26));

    -- Concatenate letters and modified Unix time
    RETURN random_letters || unix_time;
END;
$$ LANGUAGE plpgsql;

create table if not exists orders (
    id uuid primary key default uuid_generate_v4(),
    order_id text unique not null default generate_order_id(),
    user_id uuid,

    client_first_name varchar(40) not null,
    client_last_name varchar(40) not null,
    client_phone_number varchar(12) not null,
    client_comment text,

    branch_id uuid,
    status status_order_enum default 'in_process',
    total_price numeric,
    payment_type payment_order_enum,
    
    delivery_type delivery_order_enum not null default 'deliver',
    delivery_addr_lat float not null,
    delivery_addr_long float not null,
    delivery_addr_name varchar(255) default null,

    picker_user_id uuid default null,
    picked_at timestamp default null,

    delivering_user_id uuid default null,
    delivered_at timestamp default null,

    created_at timestamp default now() not null,
    updated_at timestamp default now() not null,
    deleted_at timestamp default null,

    row_order serial,

    foreign key (user_id) references clients(id) on delete set null,
    foreign key (branch_id) references branches(id) on delete set null,
    foreign key (picker_user_id) references staff(id) on delete set null,
    foreign key (delivering_user_id) references staff(id) on delete set null
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

drop trigger if exists set_order_products_total_price on order_products;
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

drop trigger if exists update_order_total_after_insert_or_update on order_products;
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

drop trigger if exists update_order_total_after_delete on order_products;
CREATE TRIGGER update_order_total_after_delete
    AFTER DELETE ON order_products
    FOR EACH ROW
EXECUTE FUNCTION update_order_total_delete();