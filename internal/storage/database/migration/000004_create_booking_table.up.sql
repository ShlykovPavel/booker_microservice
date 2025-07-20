CREATE TABLE bookings
(
    id                SERIAL PRIMARY KEY,
    user_id           BIGINT                   NOT NULL,
    booking_entity_id BIGINT                   NOT NULL,
    start_time        TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time          TIMESTAMP WITH TIME ZONE NOT NULL,
    status            VARCHAR(20)              DEFAULT 'pending',
    created_at        TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_booking_entity FOREIGN KEY (booking_entity_id) REFERENCES booking_entities (id) ON DELETE CASCADE
);

CREATE TRIGGER update_bookings_updated_at
    BEFORE UPDATE
    ON bookings
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();