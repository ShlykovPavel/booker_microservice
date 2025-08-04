-- Удаление внешнего ключа из таблицы bookings
ALTER TABLE bookings
    DROP CONSTRAINT fk_bookings_company_id;

-- Удаление столбца company_id из таблицы bookings
ALTER TABLE bookings
    DROP COLUMN company_id;

-- Удаление внешнего ключа из таблицы booking_types
ALTER TABLE booking_types
    DROP CONSTRAINT fk_booking_types_company_id;

-- Удаление столбца company_id из таблицы booking_types
ALTER TABLE booking_types
    DROP COLUMN company_id;

-- Удаление внешнего ключа из таблицы booking_entities
ALTER TABLE booking_entities
    DROP CONSTRAINT fk_booking_entities_company_id;

-- Удаление столбца company_id из таблицы booking_entities
ALTER TABLE booking_entities
    DROP COLUMN company_id;

-- Удаление таблицы companies
DROP TABLE companies;