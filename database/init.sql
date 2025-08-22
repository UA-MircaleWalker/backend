-- UA Card Battle Game Database Schema
-- PostgreSQL 14+ required for advanced JSON operations

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "btree_gin";

-- Set timezone
SET timezone = 'UTC';

-- Users table - Core user management
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    display_name VARCHAR(100),
    avatar_url VARCHAR(500),
    level INTEGER DEFAULT 1 CHECK (level >= 1),
    experience INTEGER DEFAULT 0 CHECK (experience >= 0),
    rank INTEGER DEFAULT 0,
    rank_points INTEGER DEFAULT 1000 CHECK (rank_points >= 0),
    is_active BOOLEAN DEFAULT true,
    last_login_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Cards table - Master card database with CardVariantID system
CREATE TABLE cards (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    card_number VARCHAR(20) NOT NULL, -- Base card number like "UA25BT-001"  
    card_variant_id VARCHAR(30) UNIQUE NOT NULL, -- Full variant ID like "UA25BT-001-SR★★★"
    name VARCHAR(100) NOT NULL,
    card_type VARCHAR(20) NOT NULL CHECK (card_type IN ('CHARACTER', 'FIELD', 'EVENT', 'AP')),
    color VARCHAR(10) NOT NULL CHECK (color IN ('RED', 'BLUE', 'GREEN', 'PURPLE', 'YELLOW')),
    work_code VARCHAR(6) NOT NULL, -- Work series code like "UA25BT"
    bp INTEGER CHECK (bp >= 0), -- Battle Points for character cards
    ap_cost INTEGER DEFAULT 0 CHECK (ap_cost >= 0), -- Action Points cost
    energy_cost JSONB DEFAULT '{}', -- Energy requirements {"red": 2, "blue": 1}
    energy_produce JSONB DEFAULT '{}', -- Energy production
    rarity VARCHAR(10) NOT NULL CHECK (rarity IN (
        'OBC', 'SP', 'PR', 'UR', 
        'SR_3', 'SR_2', 'SR_1', 'SR',
        'R_2', 'R_1', 'R',
        'U_3', 'U_2', 'U_1', 'U',
        'C_2', 'C_1', 'C'
    )),
    characteristics TEXT[] DEFAULT '{}', -- Card traits/characteristics
    effect_text TEXT DEFAULT '', -- Human-readable effect description
    trigger_effect VARCHAR(50) DEFAULT 'NIL' CHECK (trigger_effect IN (
        'DRAW_CARD', 'COLOR', 'ACTIVE_BP_3000', 'ADD_TO_HAND', 
        'RUSH_OR_ADD_TO_HAND', 'SPECIAL', 'FINAL', 'NIL'
    )),
    keywords TEXT[] DEFAULT '{}', -- Keywords like レイド, 狙い撃ち, etc.
    image_url VARCHAR(500),
    is_banned BOOLEAN DEFAULT false, -- For card balance management
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Card instances - represents specific card copies in collections/decks
CREATE TABLE card_instances (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    card_variant_id VARCHAR(30) NOT NULL, -- References cards.card_variant_id
    user_id UUID REFERENCES users(id) ON DELETE CASCADE, -- Owner (if in collection)
    deck_id UUID REFERENCES decks(id) ON DELETE CASCADE, -- Deck (if in deck)
    quantity INTEGER DEFAULT 1 CHECK (quantity >= 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT card_variant_exists CHECK (
        EXISTS (SELECT 1 FROM cards WHERE cards.card_variant_id = card_instances.card_variant_id)
    )
);

-- User decks (updated structure)
CREATE TABLE decks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT DEFAULT '',
    color VARCHAR(10) CHECK (color IN ('RED', 'BLUE', 'GREEN', 'PURPLE', 'YELLOW')), -- Deck color restriction
    is_active BOOLEAN DEFAULT false,
    is_public BOOLEAN DEFAULT false,
    total_cards INTEGER DEFAULT 0 CHECK (total_cards >= 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT valid_deck_size CHECK (total_cards = 50), -- Union Arena uses exactly 50 cards
    CONSTRAINT one_active_deck_per_user EXCLUDE (user_id WITH =) WHERE (is_active = true)
);

-- Games table - Game instances
CREATE TABLE games (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    player1_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    player2_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'WAITING' 
        CHECK (status IN ('WAITING', 'IN_PROGRESS', 'COMPLETED', 'ABANDONED')),
    current_turn INTEGER DEFAULT 1 CHECK (current_turn >= 1),
    phase VARCHAR(20) DEFAULT 'START' 
        CHECK (phase IN ('START', 'MOVE', 'MAIN', 'ATTACK', 'END')),
    active_player UUID NOT NULL, -- References either player1_id or player2_id
    game_state JSONB NOT NULL DEFAULT '{}', -- Complete game state
    winner UUID REFERENCES users(id), -- NULL for ongoing games
    game_mode VARCHAR(20) DEFAULT 'RANKED' CHECK (game_mode IN ('RANKED', 'CASUAL', 'FRIEND')),
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT valid_active_player CHECK (active_player = player1_id OR active_player = player2_id),
    CONSTRAINT valid_winner CHECK (winner IS NULL OR winner = player1_id OR winner = player2_id)
);

-- Game actions log - For replay and audit
CREATE TABLE game_actions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    player_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    action_type VARCHAR(30) NOT NULL CHECK (action_type IN (
        'DRAW_CARD', 'PLAY_CARD', 'ATTACK', 'BLOCK', 'ACTIVATE_EFFECT',
        'MOVE_CHARACTER', 'END_PHASE', 'END_TURN', 'SURRENDER'
    )),
    action_data JSONB NOT NULL DEFAULT '{}',
    turn INTEGER NOT NULL CHECK (turn >= 1),
    phase VARCHAR(20) NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    is_valid BOOLEAN DEFAULT true,
    error_msg TEXT,
    sequence_number INTEGER NOT NULL, -- For action ordering
    UNIQUE(game_id, sequence_number)
);

-- Game results - Final game outcomes
CREATE TABLE game_results (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    player1_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    player2_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    winner UUID REFERENCES users(id), -- NULL for draws
    game_duration INTEGER NOT NULL CHECK (game_duration > 0), -- in seconds
    total_turns INTEGER NOT NULL CHECK (total_turns > 0),
    end_reason VARCHAR(50) NOT NULL CHECK (end_reason IN (
        'NORMAL_WIN', 'SURRENDER', 'TIMEOUT', 'DECK_OUT', 'CONNECTION_LOST', 'DRAW'
    )),
    game_mode VARCHAR(20) DEFAULT 'RANKED',
    completed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Player statistics - Calculated stats for performance
CREATE TABLE player_stats (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    games_played INTEGER DEFAULT 0 CHECK (games_played >= 0),
    games_won INTEGER DEFAULT 0 CHECK (games_won >= 0),
    games_lost INTEGER DEFAULT 0 CHECK (games_lost >= 0),
    games_drawn INTEGER DEFAULT 0 CHECK (games_drawn >= 0),
    win_rate DECIMAL(5,4) DEFAULT 0.0000 CHECK (win_rate >= 0 AND win_rate <= 1),
    current_streak INTEGER DEFAULT 0, -- Can be negative for lose streaks
    best_streak INTEGER DEFAULT 0 CHECK (best_streak >= 0),
    worst_streak INTEGER DEFAULT 0 CHECK (worst_streak <= 0),
    total_game_time INTEGER DEFAULT 0 CHECK (total_game_time >= 0), -- in seconds
    avg_game_time INTEGER DEFAULT 0 CHECK (avg_game_time >= 0),
    rank_points INTEGER DEFAULT 1000 CHECK (rank_points >= 0),
    previous_rank INTEGER DEFAULT 0,
    current_rank INTEGER DEFAULT 0,
    peak_rank INTEGER DEFAULT 0,
    peak_rank_points INTEGER DEFAULT 1000,
    last_played TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Achievements system
CREATE TABLE achievements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT NOT NULL,
    icon_url VARCHAR(500),
    type VARCHAR(20) NOT NULL CHECK (type IN ('MILESTONE', 'STREAK', 'RANK', 'SPECIAL', 'SEASONAL')),
    condition JSONB NOT NULL, -- Condition logic for unlocking
    reward JSONB DEFAULT '{}', -- Rewards granted (XP, cosmetics, etc.)
    is_active BOOLEAN DEFAULT true,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- User achievements - Tracking unlocked achievements
CREATE TABLE user_achievements (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    achievement_id UUID REFERENCES achievements(id) ON DELETE CASCADE,
    progress INTEGER DEFAULT 0,
    unlocked_at TIMESTAMP WITH TIME ZONE,
    PRIMARY KEY (user_id, achievement_id)
);

-- Card usage analytics
CREATE TABLE card_usage_stats (
    card_id UUID REFERENCES cards(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    games_played INTEGER DEFAULT 0,
    games_won INTEGER DEFAULT 0,
    total_copies_used INTEGER DEFAULT 0,
    avg_turn_played DECIMAL(4,2),
    win_rate DECIMAL(5,4),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (card_id, date)
);

-- Daily game statistics
CREATE TABLE daily_stats (
    date DATE PRIMARY KEY,
    total_games INTEGER DEFAULT 0,
    unique_players INTEGER DEFAULT 0,
    new_registrations INTEGER DEFAULT 0,
    avg_game_duration INTEGER DEFAULT 0,
    ranked_games INTEGER DEFAULT 0,
    casual_games INTEGER DEFAULT 0,
    peak_concurrent_games INTEGER DEFAULT 0,
    peak_hour INTEGER,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Matchmaking history for analytics
CREATE TABLE matchmaking_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    player1_id UUID NOT NULL REFERENCES users(id),
    player2_id UUID NOT NULL REFERENCES users(id),
    mode VARCHAR(20) NOT NULL,
    wait_time_seconds INTEGER NOT NULL,
    rank_difference INTEGER,
    match_quality_score DECIMAL(3,2), -- 0.0 to 1.0
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for performance optimization

-- User indexes
CREATE INDEX idx_users_username ON users(username) WHERE is_active = true;
CREATE INDEX idx_users_email ON users(email) WHERE is_active = true;
CREATE INDEX idx_users_rank_points ON users(rank_points DESC) WHERE is_active = true;
CREATE INDEX idx_users_last_login ON users(last_login_at DESC) WHERE is_active = true;

-- Card indexes
CREATE INDEX idx_cards_number ON cards(card_number);
CREATE INDEX idx_cards_variant_id ON cards(card_variant_id);
CREATE INDEX idx_cards_type ON cards(card_type);
CREATE INDEX idx_cards_color ON cards(color);
CREATE INDEX idx_cards_work ON cards(work_code);
CREATE INDEX idx_cards_rarity ON cards(rarity);
CREATE INDEX idx_cards_keywords ON cards USING gin(keywords);
CREATE INDEX idx_cards_characteristics ON cards USING gin(characteristics);
CREATE INDEX idx_cards_bp ON cards(bp) WHERE bp IS NOT NULL;
CREATE INDEX idx_cards_ap_cost ON cards(ap_cost);
CREATE INDEX idx_cards_trigger_effect ON cards(trigger_effect);
CREATE INDEX idx_cards_compound_search ON cards(card_number, rarity); -- For variant searches

-- Card instances indexes
CREATE INDEX idx_instances_variant_id ON card_instances(card_variant_id);
CREATE INDEX idx_instances_user ON card_instances(user_id) WHERE user_id IS NOT NULL;
CREATE INDEX idx_instances_deck ON card_instances(deck_id) WHERE deck_id IS NOT NULL;
CREATE INDEX idx_instances_user_quantity ON card_instances(user_id, quantity) WHERE quantity > 0 AND user_id IS NOT NULL;

-- Deck indexes
CREATE INDEX idx_decks_user ON decks(user_id);
CREATE INDEX idx_decks_active ON decks(user_id, is_active) WHERE is_active = true;
CREATE INDEX idx_decks_color ON decks(color) WHERE color IS NOT NULL;
CREATE INDEX idx_decks_public ON decks(is_public) WHERE is_public = true;

-- Game indexes
CREATE INDEX idx_games_players ON games(player1_id, player2_id);
CREATE INDEX idx_games_status ON games(status);
CREATE INDEX idx_games_active_player ON games(active_player) WHERE status = 'IN_PROGRESS';
CREATE INDEX idx_games_created_at ON games(created_at DESC);
CREATE INDEX idx_games_mode ON games(game_mode);

-- Game actions indexes (partitioned by game_id for performance)
CREATE INDEX idx_actions_game ON game_actions(game_id, sequence_number);
CREATE INDEX idx_actions_player ON game_actions(player_id);
CREATE INDEX idx_actions_timestamp ON game_actions(timestamp DESC);
CREATE INDEX idx_actions_type ON game_actions(action_type);

-- Game results indexes
CREATE INDEX idx_results_game ON game_results(game_id);
CREATE INDEX idx_results_players ON game_results(player1_id, player2_id);
CREATE INDEX idx_results_winner ON game_results(winner) WHERE winner IS NOT NULL;
CREATE INDEX idx_results_completed_at ON game_results(completed_at DESC);
CREATE INDEX idx_results_mode ON game_results(game_mode);

-- Stats indexes
CREATE INDEX idx_stats_rank_points ON player_stats(rank_points DESC);
CREATE INDEX idx_stats_games_played ON player_stats(games_played DESC);
CREATE INDEX idx_stats_win_rate ON player_stats(win_rate DESC) WHERE games_played >= 10;
CREATE INDEX idx_stats_last_played ON player_stats(last_played DESC);

-- Analytics indexes
CREATE INDEX idx_card_usage_date ON card_usage_stats(date DESC);
CREATE INDEX idx_card_usage_win_rate ON card_usage_stats(win_rate DESC) WHERE games_played >= 10;
CREATE INDEX idx_daily_stats_date ON daily_stats(date DESC);

-- Triggers for automatic updates

-- Update user updated_at timestamp
CREATE OR REPLACE FUNCTION update_user_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_user_timestamp
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_user_timestamp();

-- Update deck updated_at timestamp
CREATE TRIGGER update_deck_timestamp
    BEFORE UPDATE ON decks
    FOR EACH ROW
    EXECUTE FUNCTION update_user_timestamp();

-- Update game updated_at timestamp
CREATE TRIGGER update_game_timestamp
    BEFORE UPDATE ON games
    FOR EACH ROW
    EXECUTE FUNCTION update_user_timestamp();

-- Function to recalculate player stats after game completion
CREATE OR REPLACE FUNCTION update_player_stats_after_game()
RETURNS TRIGGER AS $$
DECLARE
    p1_won BOOLEAN := (NEW.winner = NEW.player1_id);
    p2_won BOOLEAN := (NEW.winner = NEW.player2_id);
    is_draw BOOLEAN := (NEW.winner IS NULL);
BEGIN
    -- Update player1 stats
    INSERT INTO player_stats (user_id, games_played, games_won, games_lost, games_drawn, 
                             total_game_time, last_played, updated_at)
    VALUES (NEW.player1_id, 1, 
            CASE WHEN p1_won THEN 1 ELSE 0 END,
            CASE WHEN p2_won THEN 1 ELSE 0 END,
            CASE WHEN is_draw THEN 1 ELSE 0 END,
            NEW.game_duration, NEW.completed_at, NOW())
    ON CONFLICT (user_id) DO UPDATE SET
        games_played = player_stats.games_played + 1,
        games_won = player_stats.games_won + CASE WHEN p1_won THEN 1 ELSE 0 END,
        games_lost = player_stats.games_lost + CASE WHEN p2_won THEN 1 ELSE 0 END,
        games_drawn = player_stats.games_drawn + CASE WHEN is_draw THEN 1 ELSE 0 END,
        total_game_time = player_stats.total_game_time + NEW.game_duration,
        avg_game_time = (player_stats.total_game_time + NEW.game_duration) / (player_stats.games_played + 1),
        win_rate = CASE WHEN (player_stats.games_played + 1) > 0 
                   THEN (player_stats.games_won + CASE WHEN p1_won THEN 1 ELSE 0 END)::decimal / (player_stats.games_played + 1)
                   ELSE 0 END,
        last_played = NEW.completed_at,
        updated_at = NOW();

    -- Update player2 stats
    INSERT INTO player_stats (user_id, games_played, games_won, games_lost, games_drawn,
                             total_game_time, last_played, updated_at)
    VALUES (NEW.player2_id, 1,
            CASE WHEN p2_won THEN 1 ELSE 0 END,
            CASE WHEN p1_won THEN 1 ELSE 0 END,
            CASE WHEN is_draw THEN 1 ELSE 0 END,
            NEW.game_duration, NEW.completed_at, NOW())
    ON CONFLICT (user_id) DO UPDATE SET
        games_played = player_stats.games_played + 1,
        games_won = player_stats.games_won + CASE WHEN p2_won THEN 1 ELSE 0 END,
        games_lost = player_stats.games_lost + CASE WHEN p1_won THEN 1 ELSE 0 END,
        games_drawn = player_stats.games_drawn + CASE WHEN is_draw THEN 1 ELSE 0 END,
        total_game_time = player_stats.total_game_time + NEW.game_duration,
        avg_game_time = (player_stats.total_game_time + NEW.game_duration) / (player_stats.games_played + 1),
        win_rate = CASE WHEN (player_stats.games_played + 1) > 0 
                   THEN (player_stats.games_won + CASE WHEN p2_won THEN 1 ELSE 0 END)::decimal / (player_stats.games_played + 1)
                   ELSE 0 END,
        last_played = NEW.completed_at,
        updated_at = NOW();

    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_player_stats_trigger
    AFTER INSERT ON game_results
    FOR EACH ROW
    EXECUTE FUNCTION update_player_stats_after_game();

-- Views for common queries

-- Active leaderboard view
CREATE VIEW leaderboard_view AS
SELECT 
    u.id,
    u.username,
    u.display_name,
    u.avatar_url,
    u.level,
    ps.rank_points,
    ps.games_played,
    ps.games_won,
    ps.win_rate,
    ps.current_streak,
    ps.current_rank,
    ROW_NUMBER() OVER (ORDER BY ps.rank_points DESC, ps.games_won DESC) as position
FROM users u
JOIN player_stats ps ON u.id = ps.user_id
WHERE u.is_active = true 
  AND ps.games_played > 0
ORDER BY ps.rank_points DESC, ps.games_won DESC;

-- Recent games view
CREATE VIEW recent_games_view AS
SELECT 
    g.id as game_id,
    g.created_at,
    g.completed_at,
    g.game_mode,
    g.status,
    u1.username as player1_username,
    u2.username as player2_username,
    CASE 
        WHEN g.winner = g.player1_id THEN u1.username
        WHEN g.winner = g.player2_id THEN u2.username
        ELSE 'Draw'
    END as winner_username,
    gr.game_duration,
    gr.total_turns
FROM games g
JOIN users u1 ON g.player1_id = u1.id
JOIN users u2 ON g.player2_id = u2.id
LEFT JOIN game_results gr ON g.id = gr.game_id
WHERE g.status = 'COMPLETED'
ORDER BY g.completed_at DESC;

-- Popular cards view (updated for CardVariantID)
CREATE VIEW popular_cards_view AS
SELECT 
    c.id,
    c.card_number,
    c.card_variant_id,
    c.name,
    c.card_type,
    c.color,
    c.rarity,
    c.work_code,
    COALESCE(SUM(cus.games_played), 0) as times_used,
    COALESCE(AVG(cus.win_rate), 0) as avg_win_rate,
    COUNT(DISTINCT cus.date) as days_tracked
FROM cards c
LEFT JOIN card_usage_stats cus ON c.id = cus.card_id
WHERE c.is_banned = false
GROUP BY c.id, c.card_number, c.card_variant_id, c.name, c.card_type, c.color, c.rarity, c.work_code
HAVING COALESCE(SUM(cus.games_played), 0) > 0
ORDER BY times_used DESC;

-- Card variants view - shows all variants of each base card
CREATE VIEW card_variants_view AS
SELECT 
    c.card_number,
    c.name,
    c.card_type,
    c.color,
    c.work_code,
    COUNT(*) as variant_count,
    ARRAY_AGG(c.rarity ORDER BY 
        CASE c.rarity 
            WHEN 'OBC' THEN 10
            WHEN 'SP' THEN 9 WHEN 'PR' THEN 9
            WHEN 'UR' THEN 8
            WHEN 'SR_3' THEN 7 WHEN 'SR_2' THEN 6 WHEN 'SR_1' THEN 5 WHEN 'SR' THEN 4
            WHEN 'R_2' THEN 3 WHEN 'R_1' THEN 2 WHEN 'R' THEN 1
            WHEN 'U_3' THEN 0 WHEN 'U_2' THEN -1 WHEN 'U_1' THEN -2 WHEN 'U' THEN -3
            WHEN 'C_2' THEN -4 WHEN 'C_1' THEN -5 WHEN 'C' THEN -6
            ELSE -10
        END DESC
    ) as available_rarities,
    MAX(c.updated_at) as last_updated
FROM cards c
WHERE c.is_banned = false
GROUP BY c.card_number, c.name, c.card_type, c.color, c.work_code
ORDER BY c.card_number;

-- Insert some sample achievements
INSERT INTO achievements (name, description, type, condition, reward) VALUES
('First Victory', 'Win your first game', 'MILESTONE', '{"type": "games_won", "value": 1}', '{"experience": 100}'),
('Winning Streak', 'Win 5 games in a row', 'STREAK', '{"type": "win_streak", "value": 5}', '{"experience": 250, "title": "Streak Master"}'),
('Veteran Player', 'Play 100 games', 'MILESTONE', '{"type": "games_played", "value": 100}', '{"experience": 500, "avatar_frame": "veteran"}'),
('Rising Star', 'Reach 1500 rank points', 'RANK', '{"type": "rank_points", "value": 1500}', '{"experience": 300, "title": "Rising Star"}'),
('Champion', 'Reach 2000 rank points', 'RANK', '{"type": "rank_points", "value": 2000}', '{"experience": 500, "title": "Champion"}'),
('Card Collector', 'Own 200+ different cards', 'MILESTONE', '{"type": "unique_cards", "value": 200}', '{"experience": 400}'),
('Speed Demon', 'Win a game in under 5 minutes', 'SPECIAL', '{"type": "quick_win", "value": 300}', '{"experience": 150, "title": "Speed Demon"}'),
('Marathon Master', 'Win a game lasting over 30 minutes', 'SPECIAL', '{"type": "long_win", "value": 1800}', '{"experience": 200, "title": "Marathon Master"}'),
('Perfect Week', 'Win 7 games without losing in 7 days', 'STREAK', '{"type": "perfect_week", "value": 7}', '{"experience": 350, "title": "Perfect Week"}'),
('Deck Master', 'Win with 10 different deck compositions', 'SPECIAL', '{"type": "deck_variety", "value": 10}', '{"experience": 300, "deck_slot": 1}');

-- Sample work codes and their themes
COMMENT ON COLUMN cards.work_code IS 'Work series codes: UA25BT (25th Booster), UA25ST (25th Starter), etc.';
COMMENT ON COLUMN cards.card_variant_id IS 'Unique identifier combining card number and rarity (e.g., UA25BT-001-SR★★★)';
COMMENT ON TABLE card_instances IS 'Represents specific card copies owned by users or used in decks';

-- Add trigger to update deck total_cards when card_instances change
CREATE OR REPLACE FUNCTION update_deck_card_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' OR TG_OP = 'UPDATE' THEN
        IF NEW.deck_id IS NOT NULL THEN
            UPDATE decks SET 
                total_cards = (
                    SELECT COALESCE(SUM(quantity), 0) 
                    FROM card_instances 
                    WHERE deck_id = NEW.deck_id
                ),
                updated_at = NOW()
            WHERE id = NEW.deck_id;
        END IF;
    END IF;
    
    IF TG_OP = 'DELETE' OR TG_OP = 'UPDATE' THEN
        IF OLD.deck_id IS NOT NULL THEN
            UPDATE decks SET 
                total_cards = (
                    SELECT COALESCE(SUM(quantity), 0) 
                    FROM card_instances 
                    WHERE deck_id = OLD.deck_id
                ),
                updated_at = NOW()
            WHERE id = OLD.deck_id;
        END IF;
    END IF;
    
    RETURN COALESCE(NEW, OLD);
END;
$$ language 'plpgsql';

CREATE TRIGGER update_deck_card_count_trigger
    AFTER INSERT OR UPDATE OR DELETE ON card_instances
    FOR EACH ROW
    EXECUTE FUNCTION update_deck_card_count();

-- Performance monitoring query examples
COMMENT ON TABLE daily_stats IS 'Tracks daily game metrics for analytics and monitoring';
COMMENT ON TABLE card_usage_stats IS 'Tracks card popularity and balance metrics';
COMMENT ON TABLE matchmaking_history IS 'Records matchmaking quality metrics for tuning';

-- Partition game_actions by date for better performance (optional, for high-volume production)
-- This would be implemented when action volume becomes very high
-- CREATE TABLE game_actions_y2024m01 PARTITION OF game_actions FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');

COMMIT;