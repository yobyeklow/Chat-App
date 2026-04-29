--Delete trigger
DROP TRIGGER IF EXISTS set_user_update_at ON users;
--Delete trigger function
DROP FUNCTION IF EXISTS update_user_updated_at;
--Drop table
DROP TABLE IF EXISTS users;
