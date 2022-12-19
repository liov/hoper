# 更新行自动更新updated_at
CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
NEW.updated_at = now();
RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_table_name_update_at BEFORE UPDATE ON table_name FOR EACH ROW EXECUTE PROCEDURE  update_modified_column();