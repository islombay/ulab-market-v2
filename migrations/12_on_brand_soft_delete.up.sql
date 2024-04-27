CREATE OR REPLACE FUNCTION func_update_product_on_brand_delete()
    RETURNS TRIGGER AS $$
BEGIN
    -- Check if the icon is being soft deleted
    IF OLD.deleted_at IS NULL AND NEW.deleted_at IS NOT NULL THEN
        -- Set the icon_id to NULL in the category table for all categories using this icon
        UPDATE products
        SET brand_id = NULL
        WHERE brand_id = OLD.id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

drop trigger if exists  trg_update_products_on_brand_delete on products;
CREATE TRIGGER trg_update_products_on_brand_delete
    BEFORE UPDATE ON brands
    FOR EACH ROW
EXECUTE FUNCTION func_update_product_on_brand_delete();


CREATE OR REPLACE FUNCTION func_update_product_on_category_delete()
    RETURNS TRIGGER AS $$
BEGIN
    -- Check if the icon is being soft deleted
    IF OLD.deleted_at IS NULL AND NEW.deleted_at IS NOT NULL THEN
        -- Set the icon_id to NULL in the category table for all categories using this icon
        UPDATE products
        SET category_id = NULL
        WHERE category_id = OLD.id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

drop trigger if exists  trg_update_products_on_category_delete on products;
CREATE TRIGGER trg_update_products_on_category_delete
    BEFORE UPDATE ON category
    FOR EACH ROW
EXECUTE FUNCTION func_update_product_on_category_delete();