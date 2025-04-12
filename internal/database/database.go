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

	// Update or insert leaderboard entry
	_, err = tx.Exec(`
        INSERT INTO leaderboard (user_id, total_score, rank)
        SELECT user_id, SUM(score) as total_score, 0
        FROM game_sessions
        WHERE user_id = $1
        GROUP BY user_id
        ON CONFLICT (user_id) DO UPDATE
        SET total_score = (
            SELECT SUM(score)
            FROM game_sessions
            WHERE user_id = $1
        )`,
		session.UserID)
	if err != nil {
		return err
	}

	// Update ranks for all players
	_, err = tx.Exec(`
        WITH ranked_scores AS (
            SELECT id, user_id, total_score,
                   ROW_NUMBER() OVER (ORDER BY total_score DESC) as new_rank
            FROM leaderboard
        )
        UPDATE leaderboard l
        SET rank = r.new_rank
        FROM ranked_scores r
        WHERE l.id = r.id`)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (db *DB) GetTopPlayers(limit int) ([]models.Leaderboard, error) {
	var leaderboard []models.Leaderboard
	err := db.Select(&leaderboard, `
        SELECT l.*, u.username
        FROM leaderboard l
        JOIN users u ON l.user_id = u.id
        ORDER BY l.total_score DESC
        LIMIT $1`,
		limit)
	return leaderboard, err
}

func (db *DB) GetPlayerRank(userID int64) (*models.Leaderboard, error) {
	var leaderboard models.Leaderboard
	err := db.Get(&leaderboard, `
        SELECT l.*, u.username
        FROM leaderboard l
        JOIN users u ON l.user_id = u.id
        WHERE l.user_id = $1`,
		userID)
	if err != nil {
		return nil, err
	}
	return &leaderboard, nil
}
