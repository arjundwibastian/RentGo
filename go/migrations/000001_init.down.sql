DROP TRIGGER IF EXISTS trg_users_updated_at ON users;
DROP TRIGGER IF EXISTS trg_vehicles_updated_at ON vehicles;
DROP TRIGGER IF EXISTS trg_bookings_updated_at ON bookings;


DROP TABLE IF EXISTS bookings;
DROP TABLE IF EXISTS vehicles;
DROP TABLE IF EXISTS users;



DROP FUNCTION IF EXISTS update_updated_at_column();
DROP TYPE IF EXISTS booking_status;
DROP TYPE IF EXISTS user_role;