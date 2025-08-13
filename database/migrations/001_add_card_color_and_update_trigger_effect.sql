-- Migration: Add color field and update trigger_effect field in cards table
-- This migration adds the color system and updates trigger effects to use simple text

BEGIN;

-- Add color column to cards table
ALTER TABLE cards 
ADD COLUMN color VARCHAR(10) DEFAULT 'RED' 
CHECK (color IN ('RED', 'BLUE', 'GREEN', 'PURPLE', 'YELLOW'));

-- Update trigger_effect column to use simple text instead of JSONB
-- First create a temporary column
ALTER TABLE cards ADD COLUMN trigger_effect_new VARCHAR(50) DEFAULT 'NIL';

-- Update all existing records to use the new format
-- This assumes existing trigger_effect JSONB contains simple effect types
UPDATE cards 
SET trigger_effect_new = 'NIL' 
WHERE trigger_effect = '[]' OR trigger_effect = '{}' OR trigger_effect IS NULL;

-- Add check constraint for valid trigger effects
ALTER TABLE cards 
ADD CONSTRAINT valid_trigger_effect 
CHECK (trigger_effect_new IN (
    'DRAW_CARD', 'COLOR', 'ACTIVE_BP_3000', 'ADD_TO_HAND', 
    'RUSH_OR_ADD_TO_HAND', 'SPECIAL', 'FINAL', 'NIL'
));

-- Drop old column and rename new column
ALTER TABLE cards DROP COLUMN trigger_effect;
ALTER TABLE cards RENAME COLUMN trigger_effect_new TO trigger_effect;

-- Add index for color-based searches
CREATE INDEX idx_cards_color ON cards(color);

-- Update deck validation to enforce single-color deck requirement
-- Add a stored computed column to decks table to track deck color
ALTER TABLE decks 
ADD COLUMN deck_color VARCHAR(10) DEFAULT NULL;

-- Update deck size constraint to be exactly 50 cards (Union Arena rule)
ALTER TABLE decks 
DROP CONSTRAINT valid_deck_size;

ALTER TABLE decks 
ADD CONSTRAINT valid_deck_size CHECK (total_cards = 50);

-- Add constraint to ensure all cards in a deck are the same color (enforced at application level)
-- This will be handled by application validation in the deck building service

-- Add comment explaining the color system
COMMENT ON COLUMN cards.color IS 'Card color: RED, BLUE, GREEN, PURPLE, YELLOW. Each deck must contain only one color.';
COMMENT ON COLUMN cards.trigger_effect IS 'Trigger effect type: DRAW_CARD, COLOR (color-specific effect), ACTIVE_BP_3000, ADD_TO_HAND, RUSH_OR_ADD_TO_HAND, SPECIAL, FINAL, NIL';
COMMENT ON COLUMN decks.deck_color IS 'The color of all cards in this deck. NULL for incomplete decks.';

-- Update existing cards with default colors (this should be updated with real data)
-- For demo purposes, we'll distribute colors somewhat evenly
UPDATE cards SET color = 'RED' WHERE MOD(EXTRACT(MICROSECONDS FROM created_at)::INTEGER, 5) = 0;
UPDATE cards SET color = 'BLUE' WHERE MOD(EXTRACT(MICROSECONDS FROM created_at)::INTEGER, 5) = 1;
UPDATE cards SET color = 'GREEN' WHERE MOD(EXTRACT(MICROSECONDS FROM created_at)::INTEGER, 5) = 2;
UPDATE cards SET color = 'PURPLE' WHERE MOD(EXTRACT(MICROSECONDS FROM created_at)::INTEGER, 5) = 3;
UPDATE cards SET color = 'YELLOW' WHERE MOD(EXTRACT(MICROSECONDS FROM created_at)::INTEGER, 5) = 4;

-- Create a view for color-specific card queries
CREATE VIEW cards_by_color_view AS
SELECT 
    color,
    COUNT(*) as total_cards,
    COUNT(*) FILTER (WHERE card_type = 'CHARACTER') as character_cards,
    COUNT(*) FILTER (WHERE card_type = 'EVENT') as event_cards,
    COUNT(*) FILTER (WHERE card_type = 'FIELD') as field_cards,
    COUNT(*) FILTER (WHERE card_type = 'AP') as ap_cards,
    AVG(bp) FILTER (WHERE bp IS NOT NULL) as avg_bp,
    AVG(ap_cost) as avg_ap_cost
FROM cards
WHERE NOT is_banned
GROUP BY color
ORDER BY color;

COMMENT ON VIEW cards_by_color_view IS 'Summary statistics of cards grouped by color for deck building analysis';

COMMIT;