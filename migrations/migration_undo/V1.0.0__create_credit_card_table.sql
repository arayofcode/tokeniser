-- Drop the trigger
DROP TRIGGER IF EXISTS update_credit_cards_modtime ON credit_cards;

-- Drop the update_modified_column() function
DROP FUNCTION IF EXISTS update_modified_column;

-- Drop the credit_cards table
DROP TABLE IF EXISTS credit_cards;

-- Drop the extension pgcrypto
DROP EXTENSION IF EXISTS pgcrypto;

-- We use this extension to generate uuid
CREATE EXTENSION IF NOT EXISTS pgcrypto;
