drop trigger if exists set_order_products_total_price on order_products;
drop function if exists calculate_order_products_total_price();

drop trigger if exists update_order_total_after_delete on order_products;
drop function if exists update_order_total_delete();

drop trigger if exists update_order_total_after_insert_or_update on order_products;
drop function if exists update_order_total();

drop table if exists order_products;
drop table if exists orders;

drop type if exists payment_order_enum;
drop type if exists status_order_enum;