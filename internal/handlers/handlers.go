package handlers

import (
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gaming-leaderboard/internal/database"
	"github.com/gaming-leaderboard/internal/models"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	db *database.DB
}

// NewHandler creates a new instance of Handler
func NewHandler(db *database.DB) *Handler {
	return &Handler{db: db}
}

// SubmitScore handles the submission of new game scores
func (h *Handler) SubmitScore(c *gin.Context) {
	var submission models.ScoreSubmission
	if err := c.ShouldBindJSON(&submission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Validate score
	if submission.Score < 0 || submission.Score > 50000 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid score. Score must be between 0 and 10000.",
		})
		return
	}

	session := models.GameSession{
		UserID:   submission.UserID,
		Score:    submission.Score,
		GameMode: []string{"solo", "team"}[rand.Intn(2)],
	}

	if err := h.db.SubmitScore(session); err != nil {
		if err.Error() == "user does not exist" {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Failed to submit score",
				"details": err.Error(),
			})
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to submit score",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Score submitted successfully",
		"user_id": submission.UserID,
		"score":   submission.Score,
	})
}

// GetLeaderboard retrieves the top players
func (h *Handler) GetLeaderboard(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	leaderboard, err := h.db.GetTopPlayers(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve leaderboard",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"leaderboard": leaderboard,
		"count":       len(leaderboard),
	})
}

// GetPlayerRank retrieves the rank for a specific player
func (h *Handler) GetPlayerRank(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	rank, err := h.db.GetPlayerRank(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Player not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, rank)
}
