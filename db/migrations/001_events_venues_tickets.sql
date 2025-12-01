-- +goose Up
CREATE TABLE venues (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    location VARCHAR(255) NOT NULL
);

Create table events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL UNIQUE,
    description TEXT NOT NULL,
    start_date TIMESTAMP NOT NULL,
    venue_id UUID NOT NULL REFERENCES venues(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE ticket_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- event_id UUID NOT NULL REFERENCES events(id), if we want to allow different ticket types for different events
    name VARCHAR(255) NOT NULL UNIQUE,
    display_name TEXT NOT NULL,
    price_cents INTEGER NOT NULL
);

-- seed data for ticket_types
INSERT INTO ticket_types (name, display_name, price_cents) VALUES ('vip', 'VIP', 10000);
INSERT INTO ticket_types (name, display_name, price_cents) VALUES ('ga', 'General Admission', 1000);
INSERT INTO ticket_types (name, display_name, price_cents) VALUES ('front_row', 'Front Row', 5000);

CREATE TYPE ticket_status AS ENUM ('available', 'sold');

CREATE TABLE tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES events(id),
    ticket_type_id UUID NOT NULL REFERENCES ticket_types(id),
    status ticket_status NOT NULL DEFAULT 'available',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- Trigger for events table
CREATE TRIGGER trigger_set_updated_at_events
BEFORE UPDATE ON events
FOR EACH ROW
EXECUTE FUNCTION set_updated_at_column();

-- Trigger for tickets table
CREATE TRIGGER trigger_set_updated_at_tickets
BEFORE UPDATE ON tickets
FOR EACH ROW
EXECUTE FUNCTION set_updated_at_column();

-- +goose Down
DROP TRIGGER trigger_set_updated_at_events ON events;
DROP TRIGGER trigger_set_updated_at_tickets ON tickets;
DROP FUNCTION set_updated_at_column;
DROP TABLE events;
DROP TABLE tickets;
DROP TABLE venues;
DROP TYPE ticket_status;
