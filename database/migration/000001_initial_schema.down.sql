-- Drop the trigger
DROP TRIGGER IF EXISTS update_credit_cards_modtime ON public.credit_cards;

-- Drop the update_modified_column() function
DROP FUNCTION IF EXISTS public.update_modified_column;

-- Drop the credit_cards table
DROP TABLE IF EXISTS public.credit_cards;
