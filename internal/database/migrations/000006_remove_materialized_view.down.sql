-- Create materialized view for leaderboard
CREATE MATERIALIZED VIEW IF NOT EXISTS mv_leaderboard AS
SELECT 
    RANK() OVER (ORDER BY total_scores.total_score DESC) as rank,
    total_scores.user_id,
    total_scores.total_score,
    u.username
FROM (
    SELECT user_id, SUM(score) as total_score
    FROM game_sessions
    GROUP BY user_id
) total_scores
JOIN users u ON u.id = total_scores.user_id;

-- Create index on materialized view
CREATE UNIQUE INDEX IF NOT EXISTS idx_mv_leaderboard_user_id ON mv_leaderboard(user_id);
CREATE INDEX IF NOT EXISTS idx_mv_leaderboard_rank ON mv_leaderboard(rank);