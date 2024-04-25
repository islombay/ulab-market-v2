DROP FUNCTION IF EXISTS update_storage_quantity CASCADE;

DROP TRIGGER IF EXISTS update_storage_after_income_products_change ON income_products;
DROP TRIGGER IF EXISTS update_storage_after_order_products_change ON order_products;
DROP TRIGGER IF EXISTS update_storage_after_order_change ON orders;
