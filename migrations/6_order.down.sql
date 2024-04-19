drop trigger set_order_products_total_price on order_products;
drop function calculate_order_products_total_price();

drop trigger update_order_total_after_delete on order_products;
drop function update_order_total_delete();

drop trigger update_order_total_after_insert_or_update on order_products;
drop function update_order_total();

drop table order_products;
drop table orders;

drop type payment_order_enum;
drop type status_order_enum;