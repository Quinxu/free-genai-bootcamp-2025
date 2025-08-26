package handlers

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	"lang-portal/internal/api/response"
	"lang-portal/internal/service"
)

// WordHandler handles word-related routes
type WordHandler struct {
	wordService *service.WordService
}

// NewWordHandler creates a new WordHandler
func NewWordHandler(wordService *service.WordService) *WordHandler {
	return &WordHandler{wordService: wordService}
}

// RegisterRoutes registers word-related routes
func (h *WordHandler) RegisterRoutes(r *gin.RouterGroup) {
	words := r.Group("/words")
	{
		words.GET("", h.GetWords)
		words.GET("/:id", h.GetWord)
	}
}

// GetWords handles GET /api/words
func (h *WordHandler) GetWords(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "100"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 100
	}

	words, err := h.wordService.GetWords(page, perPage)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	// Flatten stats per spec
	type listItem struct {
		Chinese      string `json:"chinese"`
		English      string `json:"english"`
		CorrectCount int    `json:"correct_count"`
		WrongCount   int    `json:"wrong_count"`
	}
	items := make([]listItem, 0, len(words.Items))
	for _, w := range words.Items {
		items = append(items, listItem{
			Chinese:      w.Chinese,
			English:      w.English,
			CorrectCount: w.Stats.CorrectCount,
			WrongCount:   w.Stats.WrongCount,
		})
	}

	response.Success(c, gin.H{
		"items":      items,
		"pagination": words.Pagination,
	})
}

// GetWord handles GET /api/words/:id
func (h *WordHandler) GetWord(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, errors.New("invalid word ID"))
		return
	}

	word, err := h.wordService.GetWordByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(c, errors.New("word not found"))
			return
		}
		response.InternalError(c, err)
		return
	}
	groups, err := h.wordService.GetGroupsForWord(id)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, gin.H{
		"english": word.English,
		"stats": gin.H{
			"correct_count": word.Stats.CorrectCount,
			"wrong_count":   word.Stats.WrongCount,
		},
		"groups": groups,
	})
}
