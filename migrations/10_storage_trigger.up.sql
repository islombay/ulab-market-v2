CREATE OR REPLACE FUNCTION update_storage_quantity()
    RETURNS TRIGGER AS $$
BEGIN
    WITH summed_income AS (
        SELECT income_products.product_id, incomes.branch_id, SUM(income_products.quantity) AS total_income_quantity
        FROM income_products
                 JOIN incomes ON incomes.id = income_products.income_id
        GROUP BY income_products.product_id, incomes.branch_id
    ), summed_order AS (
        SELECT order_products.product_id, orders.branch_id, SUM(order_products.quantity) AS total_order_quantity
        FROM order_products
                 JOIN orders ON orders.id = order_products.order_id
        WHERE orders.status IN ('finished',
                                'in_process',
                                'picking',
                                'delivering')
        GROUP BY order_products.product_id, orders.branch_id
    ), final_quantities AS (
        SELECT COALESCE(si.product_id, so.product_id) AS product_id,
               COALESCE(si.branch_id, so.branch_id) AS branch_id,
               COALESCE(si.total_income_quantity, 0) - COALESCE(so.total_order_quantity, 0) AS net_quantity
        FROM summed_income si
                 FULL OUTER JOIN summed_order so ON si.product_id = so.product_id
                    -- and si.branch_id = so.branch_id was deleted
                    -- because, now, branches are ignored.
                    -- when more branches are added,
                    -- undo ignore the code below, and update function
                    -- AND si.branch_id = so.branch_id
    )
    INSERT INTO storage (id, product_id, branch_id, quantity)
    SELECT uuid_generate_v4(), product_id, branch_id, net_quantity FROM final_quantities
    ON CONFLICT (product_id, branch_id) DO UPDATE
        SET quantity = EXCLUDED.quantity;

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
