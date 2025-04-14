package database

import (
	"fmt"
	"log"
	"time"

	"github.com/gaming-leaderboard/internal/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB struct {
	*sqlx.DB
	cache models.LeaderboardCache
}

func NewDB(host, user, password, dbname string, port int) (*DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	return &DB{db, models.LeaderboardCache{}}, nil
}

func (db *DB) UserExists(userID int64) (bool, error) {
	var exists bool
	start := time.Now()
	err := db.Get(&exists, `
	SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`,
		userID)
	if err != nil {
		return false, err
	}
	log.Println("UserExists query latency:", time.Since(start).Milliseconds(), "ms")
	return exists, nil
}

func (db *DB) SubmitScore(session models.GameSession) error {
	/* old implementations
		tx, err := db.Beginx()
		if err != nil {
			return err
		}
		defer tx.Rollback()

		// Insert game session
		start := time.Now()
		_, err = tx.NamedExec(`
	        INSERT INTO game_sessions (user_id, score, game_mode)
	        VALUES (:user_id, :score, :game_mode)`,
			session)
		if err != nil {
			return err
		}
		log.Println("Insert game session latency:", time.Since(start).Milliseconds(), "ms")

		// Update materialized view concurrently
		start = time.Now()
		_, err = tx.Exec(`
	        REFRESH MATERIALIZED VIEW CONCURRENTLY mv_leaderboard;
	    `)
		if err != nil {
			return err
		}
		log.Println("Update materialized view mv_leaderboard concurrently latency:", time.Since(start).Milliseconds(), "ms")
		return tx.Commit()
	*/

	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert into game_sessions
	_, err = tx.NamedExec(`
		INSERT INTO game_sessions (user_id, score, game_mode)
		VALUES (:user_id, :score, :game_mode)`,
		session)
	if err != nil {
		return err
	}

	// Update leaderboard with the new score
	_, err = tx.Exec(`
		INSERT INTO leaderboard (user_id, total_score)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE
		SET total_score = leaderboard.total_score + EXCLUDED.total_score`,
		session.UserID, session.Score)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (db *DB) GetTopPlayers(limit int) ([]models.Leaderboard, error) {
	/* old implementations
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
	*/
	db.cache.Mu.RLock()
	if time.Now().Before(db.cache.ExpiresAt) {
		defer db.cache.Mu.RUnlock()
		return db.cache.TopPlayers, nil
	}
	db.cache.Mu.RUnlock()
	var leaderboard []models.Leaderboard
	err := db.Select(&leaderboard, `
        SELECT u.username, l.total_score, 
			   RANK() OVER (ORDER BY l.total_score DESC) as rank
		FROM leaderboard l
		JOIN users u ON l.user_id = u.id
		ORDER BY l.total_score DESC
		LIMIT 10`)

	// Update cache with write lock
	db.cache.Mu.Lock()
	db.cache.TopPlayers = leaderboard
	db.cache.ExpiresAt = time.Now().Add(5 * time.Second)
	db.cache.Mu.Unlock()
	return leaderboard, err
}

func (db *DB) GetPlayerRank(userID int64) (*models.Leaderboard, error) {
	/* old implementation
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
	*/
	var leaderboard models.Leaderboard
	err := db.Get(&leaderboard, `
        SELECT (SELECT COUNT(*) + 1 FROM leaderboard WHERE total_score > l.total_score) as rank
		FROM leaderboard l
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
