-- Add index for total_score to improve leaderboard query performance
CREATE INDEX IF NOT EXISTS idx_leaderboard_total_score ON leaderboard(total_score DESC);

-- Add index for user lookups
CREATE INDEX IF NOT EXISTS idx_leaderboard_user_id ON leaderboard(user_id);