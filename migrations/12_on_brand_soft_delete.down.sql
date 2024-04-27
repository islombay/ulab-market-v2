drop trigger if exists  trg_update_products_on_brand_delete on products;

drop function if exists func_update_product_on_brand_delete();

drop trigger if exists trg_update_products_on_category_delete on products;

drop function if exists func_update_product_on_category_delete();