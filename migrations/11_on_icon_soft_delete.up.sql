CREATE OR REPLACE FUNCTION func_update_category_on_icon_delete()
    RETURNS TRIGGER AS $$
BEGIN
    -- Check if the icon is being soft deleted
    IF OLD.deleted_at IS NULL AND NEW.deleted_at IS NOT NULL THEN
        -- Set the icon_id to NULL in the category table for all categories using this icon
        UPDATE category
        SET icon_id = NULL
        WHERE icon_id = OLD.id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

drop trigger if exists  trg_update_category_on_icon_delete on icons_list;
CREATE TRIGGER trg_update_category_on_icon_delete
    BEFORE UPDATE ON icons_list
    FOR EACH ROW
EXECUTE FUNCTION func_update_category_on_icon_delete();