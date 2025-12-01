-- +goose Up
CREATE TABLE purchases (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    total_cents INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Add purchase_id to tickets table
ALTER TABLE tickets 
ADD COLUMN purchase_id UUID REFERENCES purchases(id) ON DELETE SET NULL;

CREATE TRIGGER trigger_set_updated_at_purchases
BEFORE UPDATE ON purchases
FOR EACH ROW
EXECUTE FUNCTION set_updated_at_column();

-- +goose Down
DROP TRIGGER trigger_set_updated_at_purchases ON purchases;
ALTER TABLE tickets DROP COLUMN purchase_id;
DROP TABLE purchases;