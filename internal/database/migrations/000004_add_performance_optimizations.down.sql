DROP INDEX IF EXISTS idx_mv_leaderboard_rank;
DROP INDEX IF EXISTS idx_mv_leaderboard_user_id;
DROP MATERIALIZED VIEW IF EXISTS mv_leaderboard;
DROP INDEX IF EXISTS idx_game_sessions_user_score;