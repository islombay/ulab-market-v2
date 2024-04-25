CREATE OR REPLACE FUNCTION update_storage_quantity()
    RETURNS TRIGGER AS $$
BEGIN
    -- Calculate the new quantity for the storage
WITH summed_income as (
    SELECT income_products.product_id, incomes.branch_id, SUM(income_products.quantity) as total_income_quantity
    FROM income_products
             JOIN incomes ON incomes.id = income_products.income_id
    GROUP BY income_products.product_id, incomes.branch_id
), summed_order as (
    SELECT order_products.product_id, orders.branch_id, SUM(order_products.quantity) as total_order_quantity
    FROM order_products
             JOIN orders ON orders.id = order_products.order_id
    WHERE orders.status in ('finished', 'in_process')
    GROUP BY order_products.product_id, orders.branch_id
)
UPDATE storage
SET quantity = COALESCE(summed_income.total_income_quantity, 0) - COALESCE(summed_order.total_order_quantity, 0)
    FROM summed_income
             FULL OUTER JOIN summed_order ON summed_income.product_id = summed_order.product_id AND summed_income.branch_id = summed_order.branch_id
WHERE storage.product_id = COALESCE(summed_income.product_id, summed_order.product_id)
  AND storage.branch_id = COALESCE(summed_income.branch_id, summed_order.branch_id);

RETURN NULL; -- Since this is an AFTER trigger
END;
$$ LANGUAGE plpgsql;

-- Trigger on income_products
drop trigger if exists update_storage_after_income_products_change on income_products;
CREATE TRIGGER update_storage_after_income_products_change
    AFTER INSERT OR UPDATE OR DELETE ON income_products
    FOR EACH ROW EXECUTE FUNCTION update_storage_quantity();

-- Trigger on order_products when orders are finished
drop trigger if exists update_storage_after_order_products_change on order_products;
CREATE TRIGGER update_storage_after_order_products_change
    AFTER INSERT OR UPDATE OR DELETE ON order_products
    FOR EACH ROW EXECUTE FUNCTION update_storage_quantity();

drop trigger if exists update_storage_after_order_change on orders;
create trigger update_storage_after_order_change
    after insert or update or delete on orders
    for each row execute function update_storage_quantity();
