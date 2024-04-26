create table if not exists incomes (
    id uuid primary key not null,
    branch_id uuid,
    total_price numeric,
    comment varchar(255),
    created_at timestamp default now() not null,
    updated_at timestamp default now() not null,
    deleted_at timestamp default null,

    foreign key (branch_id) references branches(id) on delete set null
);

create table if not exists income_products (
    id uuid primary key not null,
    income_id uuid,
    product_id uuid,
    quantity int,
    product_price numeric,
    total_price numeric,
    created_at timestamp default now() not null,
    updated_at timestamp default now() not null,
    deleted_at timestamp default null,
    foreign key (income_id) references incomes(id) on delete set null,
    foreign key (product_id) references products(id) on delete set null
);

create or replace function calculate_income_products_total_price()
returns trigger as
    $$
begin
        new.total_price := new.product_price * new.quantity;
return new;
end;
$$
language plpgsql;

drop trigger if exists set_income_products_total_price on income_products;
create trigger set_income_products_total_price
    before insert or update on income_products
                         for each row
                         execute function calculate_income_products_total_price();

CREATE OR REPLACE FUNCTION update_income_total() RETURNS TRIGGER AS $$
BEGIN
    -- Calculate the total price for the current income
UPDATE incomes SET total_price = (
    SELECT COALESCE(SUM(total_price), 0)
    FROM income_products
    WHERE income_id = NEW.income_id
)
WHERE id = NEW.income_id;

RETURN NEW;
END;
$$ LANGUAGE plpgsql;

drop trigger if exists update_income_total_after_insert_or_update on income_products;
CREATE TRIGGER update_income_total_after_insert_or_update
    AFTER INSERT OR UPDATE ON income_products
                        FOR EACH ROW
                        EXECUTE FUNCTION update_income_total();