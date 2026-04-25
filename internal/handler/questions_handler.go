package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/areyoush/algoroulette/internal/model"
	"github.com/areyoush/algoroulette/internal/repository"
	"github.com/gin-gonic/gin"
)

type QuestionHandler struct {
	repo *repository.QuestionRepository
}

func NewQuestionHandler(repo *repository.QuestionRepository) *QuestionHandler {
	return &QuestionHandler{repo: repo}
}

func (h *QuestionHandler) GetRandom(c *gin.Context) {
	userID := c.GetInt("user_id")
	topic := c.Query("topic")
	difficulty := c.Query("difficulty")

	q, err := h.repo.GetRandom(userID, topic, difficulty)
	if err != nil {
		log.Printf("Error getting random question: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if q == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no matching question found"})
		return
	}
	c.JSON(http.StatusOK, q)
}

func (h *QuestionHandler) Create(c *gin.Context) {
	userID := c.GetInt("user_id")

	var q model.Question
	if err := c.ShouldBindJSON(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.Insert(userID, &q); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert question"})
		return
	}
	c.JSON(http.StatusCreated, q)
}

func (h *QuestionHandler) Import(c *gin.Context) {
	userID := c.GetInt("user_id")

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1048576)
	
	var questions []model.Question
	if err := c.ShouldBindJSON(&questions); err != nil {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "file too large or invalid json (max 1MB)"})
		return
	}
	if err := h.repo.InsertBatch(userID, questions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to import questions"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"imported": len(questions)})
}

func (h *QuestionHandler) ClearAll(c *gin.Context) {
	userID := c.GetInt("user_id")

	if err := h.repo.DeleteAllForUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to clear questions"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "your questions cleared"})
}

func (h *QuestionHandler) UpdateStatus(c *gin.Context) {
	userID := c.GetInt("user_id")
	questionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var body struct {
		Status *string `json:"status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.UpsertStatus(userID, questionID, body.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update status"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"updated": questionID})
}

func (h *QuestionHandler) UpdateBookmark(c *gin.Context) {
	userID := c.GetInt("user_id")
	questionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var body struct {
		Bookmarked bool `json:"bookmarked"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.UpsertBookmark(userID, questionID, body.Bookmarked); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update bookmark"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"updated": questionID})
}

func (h *QuestionHandler) UpdateNotes(c *gin.Context) {
	userID := c.GetInt("user_id")
	questionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var body struct {
		Notes *string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.UpsertNotes(userID, questionID, body.Notes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update notes"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"updated": questionID})
}