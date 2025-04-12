-- Remove test data in reverse order of dependencies
DELETE FROM leaderboard 
WHERE user_id IN (
    SELECT id FROM users 
    WHERE username IN (
        'ProGamer123', 'GameMaster', 'PixelWarrior', 'SpeedRunner', 'QuestHunter',
        'NinjaGamer', 'DragonSlayer', 'StarCollector', 'LegendHunter', 'MysteryPlayer',
        'VirtualHero', 'CyberKnight', 'WizardKing', 'BattleMaster', 'PixelPirate',
        'CosmicRacer', 'DungeonLord', 'StealthMaster', 'PowerPlayer', 'GalaxyWarrior',
        'EpicGamer', 'TitanSlayer', 'MythicHero', 'LegacyPlayer', 'EliteWarrior',
        'PhantomGamer', 'OmegaPlayer', 'UltimatePro', 'AlphaGamer', 'LegendaryPro'
    )
);

DELETE FROM game_sessions 
WHERE user_id IN (
    SELECT id FROM users 
    WHERE username IN (
        'ProGamer123', 'GameMaster', 'PixelWarrior', 'SpeedRunner', 'QuestHunter',
        'NinjaGamer', 'DragonSlayer', 'StarCollector', 'LegendHunter', 'MysteryPlayer',
        'VirtualHero', 'CyberKnight', 'WizardKing', 'BattleMaster', 'PixelPirate',
        'CosmicRacer', 'DungeonLord', 'StealthMaster', 'PowerPlayer', 'GalaxyWarrior',
        'EpicGamer', 'TitanSlayer', 'MythicHero', 'LegacyPlayer', 'EliteWarrior',
        'PhantomGamer', 'OmegaPlayer', 'UltimatePro', 'AlphaGamer', 'LegendaryPro'
    )
);

DELETE FROM users 
WHERE username IN (
    'ProGamer123', 'GameMaster', 'PixelWarrior', 'SpeedRunner', 'QuestHunter',
    'NinjaGamer', 'DragonSlayer', 'StarCollector', 'LegendHunter', 'MysteryPlayer',
    'VirtualHero', 'CyberKnight', 'WizardKing', 'BattleMaster', 'PixelPirate',
    'CosmicRacer', 'DungeonLord', 'StealthMaster', 'PowerPlayer', 'GalaxyWarrior',
    'EpicGamer', 'TitanSlayer', 'MythicHero', 'LegacyPlayer', 'EliteWarrior',
    'PhantomGamer', 'OmegaPlayer', 'UltimatePro', 'AlphaGamer', 'LegendaryPro'
);