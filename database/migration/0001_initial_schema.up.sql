-- Create the credit_cards table
CREATE TABLE IF NOT EXISTS public.credit_cards (
    id SERIAL PRIMARY KEY,
    token UUID DEFAULT gen_random_uuid() NOT NULL,
    cardholder_name VARCHAR(255) NOT NULL,
    card_number_encrypted BYTEA NOT NULL,
    expiry_date_encrypted BYTEA NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_accessed TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create the update_modified_column() function
CREATE FUNCTION IF NOT EXISTS public.update_modified_column() RETURNS trigger
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
CREATE TRIGGER IF NOT EXISTS update_credit_cards_modtime
    BEFORE UPDATE ON public.credit_cards
    FOR EACH ROW
    EXECUTE FUNCTION public.update_modified_column();
