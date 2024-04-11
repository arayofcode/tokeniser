-- Create the credit_cards table
CREATE TABLE IF NOT EXISTS credit_cards (
    id SERIAL PRIMARY KEY,
    token UUID DEFAULT gen_random_uuid() NOT NULL,
    cardholder_name VARCHAR(255) NOT NULL CHECK (cardholder_name <> ''),
    card_number_encrypted BYTEA NOT NULL CHECK (card_number_encrypted <> ''),
    expiry_date_encrypted BYTEA NOT NULL CHECK (expiry_date_encrypted <> ''),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_accessed TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create the update_modified_column() function
CREATE OR REPLACE FUNCTION update_modified_column() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF ROW(NEW.*) IS DISTINCT FROM ROW(OLD.*) THEN
        NEW.updated_at = CURRENT_TIMESTAMP;
    END IF;
    RETURN NEW;
END;
$$;

-- Create a trigger to automatically update the updated_at column
CREATE OR REPLACE TRIGGER update_credit_cards_modtime
    BEFORE UPDATE ON credit_cards
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();
