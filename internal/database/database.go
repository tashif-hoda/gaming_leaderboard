package database

import (
	"fmt"

	"github.com/gaming-leaderboard/internal/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB struct {
	*sqlx.DB
}

func NewDB(host, user, password, dbname string, port int) (*DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func (db *DB) UserExists(userID int64) (bool, error) {
	var exists bool
	err := db.Get(&exists, `
		SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`,
		userID)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (db *DB) SubmitScore(session models.GameSession) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert game session
	_, err = tx.NamedExec(`
        INSERT INTO game_sessions (user_id, score, game_mode)
        VALUES (:user_id, :score, :game_mode)`,
		session)
	if err != nil {
		return err
	}

	// Update materialized view concurrently
	_, err = tx.Exec(`
        REFRESH MATERIALIZED VIEW CONCURRENTLY mv_leaderboard;
    `)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (db *DB) GetTopPlayers(limit int) ([]models.Leaderboard, error) {
	var leaderboard []models.Leaderboard
	err := db.Select(&leaderboard, `
        SELECT 
            user_id,
            total_score,
            rank,
            username
        FROM mv_leaderboard
        ORDER BY rank ASC
        LIMIT $1`,
		limit)
	return leaderboard, err
}

func (db *DB) GetPlayerRank(userID int64) (*models.Leaderboard, error) {
	var leaderboard models.Leaderboard
	err := db.Get(&leaderboard, `
        SELECT 
            user_id,
            total_score,
            rank,
            username
        FROM mv_leaderboard
        WHERE user_id = $1`,
		userID)
	if err != nil {
		return nil, err
	}
	return &leaderboard, nil
}

// Periodic refresh function for background updates
func (db *DB) RefreshLeaderboard() error {
	_, err := db.Exec(`REFRESH MATERIALIZED VIEW CONCURRENTLY mv_leaderboard`)
	return err
}
