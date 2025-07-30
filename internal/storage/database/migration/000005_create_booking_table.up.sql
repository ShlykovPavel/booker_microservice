CREATE TABLE companies
(
    id           SERIAL PRIMARY KEY,
    company_id   INT NOT NULL UNIQUE,
    company_name VARCHAR NOT NULL
);
ALTER TABLE booking_entities
    ADD COLUMN company_id INT;

ALTER TABLE booking_entities
    ADD CONSTRAINT fk_booking_entities_company_id

        FOREIGN KEY (company_id)
            REFERENCES companies (company_id)
            ON DELETE SET NULL;

ALTER TABLE booking_types
    ADD COLUMN company_id INT;

ALTER TABLE booking_types
    ADD CONSTRAINT fk_booking_types_company_id

        FOREIGN KEY (company_id)
            REFERENCES companies (company_id)
            ON DELETE SET NULL;

ALTER TABLE bookings
    ADD COLUMN company_id INT;

ALTER TABLE bookings
    ADD CONSTRAINT fk_bookings_company_id

        FOREIGN KEY (company_id)
            REFERENCES companies (company_id)
            ON DELETE SET NULL;