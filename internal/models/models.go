package models

import (
	"sync"
	"time"
)

type User struct {
	ID       int64     `db:"id" json:"id"`
	Username string    `db:"username" json:"username"`
	JoinDate time.Time `db:"join_date" json:"join_date"`
}

type GameSession struct {
	ID        int64     `db:"id" json:"id"`
	UserID    int64     `db:"user_id" json:"user_id"`
	Score     int       `db:"score" json:"score"`
	GameMode  string    `db:"game_mode" json:"game_mode"`
	Timestamp time.Time `db:"timestamp" json:"timestamp"`
}

type Leaderboard struct {
	ID         int64  `db:"id" json:"id"`
	UserID     int64  `db:"user_id" json:"user_id"`
	TotalScore int    `db:"total_score" json:"total_score"`
	Rank       int    `db:"rank" json:"rank"`
	Username   string `db:"username" json:"username,omitempty"` // For join queries
}

type ScoreSubmission struct {
	UserID int64 `json:"user_id" binding:"required"`
	Score  int   `json:"score" binding:"required"`
}

type LeaderboardCache struct {
	Mu         sync.RWMutex
	TopPlayers []Leaderboard
	ExpiresAt  time.Time
}
