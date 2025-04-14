-- Remove constraints in reverse order
ALTER TABLE game_sessions DROP CONSTRAINT IF EXISTS check_score_range;