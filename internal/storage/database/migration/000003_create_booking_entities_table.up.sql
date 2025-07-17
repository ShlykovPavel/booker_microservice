CREATE TABLE booking_entities
(
    id              SERIAL PRIMARY KEY,
    booking_type_id BIGINT       NOT NULL,
    name            VARCHAR(100) NOT NULL,
    description     TEXT,
    status          VARCHAR(20) DEFAULT 'available',
    parent_id       BIGINT       NULL,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_booking_type FOREIGN KEY (booking_type_id) REFERENCES booking_types (id) ON DELETE CASCADE
);

CREATE TRIGGER update_booking_entities_updated_at
    BEFORE UPDATE
    ON booking_entities
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();