drop trigger if exists update_income_total_after_insert_or_update on income_products;
drop function if exists update_income_total();

drop trigger if exists set_income_products_total_price on income_products;
drop function if exists calculate_income_products_total_price();

drop table if exists income_products;
drop table if exists incomes;