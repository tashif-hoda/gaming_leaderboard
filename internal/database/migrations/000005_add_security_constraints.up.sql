-- Add score constraints
ALTER TABLE game_sessions ADD CONSTRAINT check_score_range CHECK (score >= 0 AND score <= 50000);
