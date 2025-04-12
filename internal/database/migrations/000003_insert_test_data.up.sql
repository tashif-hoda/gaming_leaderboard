-- Insert test users
INSERT INTO users (id, username, join_date)
SELECT id, username, join_date
FROM (VALUES
    (1, 'ProGamer123', TIMESTAMP '2025-01-01 10:00:00'),
    (2, 'GameMaster', TIMESTAMP '2025-01-01 11:00:00'),
    (3, 'PixelWarrior', TIMESTAMP '2025-01-01 12:00:00'),
    (4, 'SpeedRunner', TIMESTAMP '2025-01-01 13:00:00'),
    (5, 'QuestHunter', TIMESTAMP '2025-01-01 14:00:00'),
    (6, 'NinjaGamer', TIMESTAMP '2025-01-01 15:00:00'),
    (7, 'DragonSlayer', TIMESTAMP '2025-01-01 16:00:00'),
    (8, 'StarCollector', TIMESTAMP '2025-01-01 17:00:00'),
    (9, 'LegendHunter', TIMESTAMP '2025-01-01 18:00:00'),
    (10, 'MysteryPlayer', TIMESTAMP '2025-01-01 19:00:00'),
    (11, 'VirtualHero', TIMESTAMP '2025-01-01 20:00:00'),
    (12, 'CyberKnight', TIMESTAMP '2025-01-02 10:00:00'),
    (13, 'WizardKing', TIMESTAMP '2025-01-02 11:00:00'),
    (14, 'BattleMaster', TIMESTAMP '2025-01-02 12:00:00'),
    (15, 'PixelPirate', TIMESTAMP '2025-01-02 13:00:00'),
    (16, 'CosmicRacer', TIMESTAMP '2025-01-02 14:00:00'),
    (17, 'DungeonLord', TIMESTAMP '2025-01-02 15:00:00'),
    (18, 'StealthMaster', TIMESTAMP '2025-01-02 16:00:00'),
    (19, 'PowerPlayer', TIMESTAMP '2025-01-02 17:00:00'),
    (20, 'GalaxyWarrior', TIMESTAMP '2025-01-02 18:00:00'),
    (21, 'EpicGamer', TIMESTAMP '2025-01-02 19:00:00'),
    (22, 'TitanSlayer', TIMESTAMP '2025-01-02 20:00:00'),
    (23, 'MythicHero', TIMESTAMP '2025-01-03 10:00:00'),
    (24, 'LegacyPlayer', TIMESTAMP '2025-01-03 11:00:00'),
    (25, 'EliteWarrior', TIMESTAMP '2025-01-03 12:00:00'),
    (26, 'PhantomGamer', TIMESTAMP '2025-01-03 13:00:00'),
    (27, 'OmegaPlayer', TIMESTAMP '2025-01-03 14:00:00'),
    (28, 'UltimatePro', TIMESTAMP '2025-01-03 15:00:00'),
    (29, 'AlphaGamer', TIMESTAMP '2025-01-03 16:00:00'),
    (30, 'LegendaryPro', TIMESTAMP '2025-01-03 17:00:00')
) AS new_users(id, username, join_date)
WHERE NOT EXISTS (
    SELECT 1 FROM users WHERE username = new_users.username
);

-- Create score patterns for each user
WITH RECURSIVE score_patterns AS (
    SELECT 
        u.id as user_id,
        1 as session_num,
        floor((1000 + (u.id * 100)) * 1.0) as score,
        u.join_date + interval '1 hour' as session_time
    FROM users u
    WHERE u.id <= 30
    
    UNION ALL
    
    SELECT 
        sp.user_id,
        sp.session_num + 1,
        CASE 
            WHEN sp.session_num = 1 THEN floor((1000 + (sp.user_id * 100)) * 1.2)
            WHEN sp.session_num = 2 THEN floor((1000 + (sp.user_id * 100)) * 1.5)
        END,
        sp.session_time + interval '1 hour'
    FROM score_patterns sp
    WHERE sp.session_num < 3
)
-- Insert game sessions using the generated patterns
INSERT INTO game_sessions (user_id, score, game_mode, timestamp)
SELECT 
    user_id,
    score::integer,
    'default',
    session_time
FROM score_patterns
WHERE NOT EXISTS (
    SELECT 1 FROM game_sessions gs
    WHERE gs.user_id = score_patterns.user_id
    AND gs.timestamp = score_patterns.session_time
);

-- Initialize or update leaderboard entries
INSERT INTO leaderboard (user_id, total_score, rank)
SELECT 
    user_id,
    SUM(score) as total_score,
    ROW_NUMBER() OVER (ORDER BY SUM(score) DESC) as rank
FROM game_sessions
GROUP BY user_id
ON CONFLICT (user_id) DO UPDATE
SET total_score = EXCLUDED.total_score,
    rank = EXCLUDED.rank;