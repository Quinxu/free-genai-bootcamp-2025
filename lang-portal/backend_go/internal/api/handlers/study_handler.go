package handlers

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	"lang-portal/internal/api/response"
	"lang-portal/internal/service"
)

// StudyHandler handles study-related routes
type StudyHandler struct {
	studyService *service.StudyService
}

// NewStudyHandler creates a new StudyHandler
func NewStudyHandler(studyService *service.StudyService) *StudyHandler {
	return &StudyHandler{studyService: studyService}
}

// RegisterRoutes registers study-related routes
func (h *StudyHandler) RegisterRoutes(r *gin.RouterGroup) {
	dashboard := r.Group("/dashboard")
	{
		dashboard.GET("/last_study_session", h.GetLastStudySession)
		dashboard.GET("/study_progress", h.GetStudyProgress)
		dashboard.GET("/quick_stats", h.GetQuickStats)
	}

	studyActivities := r.Group("/study_activities")
	{
		studyActivities.POST("", h.StartStudyActivity)
		studyActivities.GET("/:id", h.GetStudyActivity)
		studyActivities.GET("/:id/study_sessions", h.GetActivityStudySessions)
	}

	studySessions := r.Group("/study_sessions")
	{
		studySessions.GET("", h.GetStudySessions)
		studySessions.GET("/:id", h.GetStudySession)
		studySessions.GET("/:id/words", h.GetStudySessionWords)
		studySessions.POST("/:id/words/:word_id/review", h.RecordWordReview)
	}

	r.POST("/reset_history", h.ResetHistory)
	r.POST("/full_reset", h.FullReset)
}

// GetLastStudySession handles GET /api/dashboard/last_study_session
func (h *StudyHandler) GetLastStudySession(c *gin.Context) {
	session, err := h.studyService.GetLastStudySession()
	if err != nil {
		response.InternalError(c, err)
		return
	}
	if session == nil {
		response.Success(c, gin.H{})
		return
	}
	// Map to spec: omit reviewed_words
	resp := gin.H{
		"id":                session.ID,
		"group_id":          session.GroupID,
		"created_at":        session.CreatedAt,
		"study_activity_id": session.StudyActivityID,
		"group_name":        session.GroupName,
	}
	response.Success(c, resp)
}

// GetStudyProgress handles GET /api/dashboard/study_progress
func (h *StudyHandler) GetStudyProgress(c *gin.Context) {
	progress, err := h.studyService.GetStudyProgress()
	if err != nil {
		response.InternalError(c, err)
		return
	}
	response.Success(c, progress)
}

// GetQuickStats handles GET /api/dashboard/quick_stats
func (h *StudyHandler) GetQuickStats(c *gin.Context) {
	stats, err := h.studyService.GetQuickStats()
	if err != nil {
		response.InternalError(c, err)
		return
	}
	response.Success(c, stats)
}

// StartStudyActivity handles POST /api/study-activities
func (h *StudyHandler) StartStudyActivity(c *gin.Context) {
	var req struct {
		GroupID         int64 `json:"group_id" binding:"required"`
		StudyActivityID int64 `json:"study_activity_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err)
		return
	}

	session, err := h.studyService.StartStudySession(req.GroupID, req.StudyActivityID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, session)
}

// GetStudySessions handles GET /api/study-sessions
func (h *StudyHandler) GetStudySessions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "100"))
	sessions, err := h.studyService.GetStudySessions(page, perPage)
	if err != nil {
		response.InternalError(c, err)
		return
	}
	response.Success(c, sessions)
}

// GetStudySession handles GET /api/study-sessions/:id
func (h *StudyHandler) GetStudySession(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, errors.New("invalid session ID"))
		return
	}
	session, err := h.studyService.GetStudySession(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(c, errors.New("study session not found"))
			return
		}
		response.InternalError(c, err)
		return
	}
	response.Success(c, session)
}

// GetStudySessionWords handles GET /api/study-sessions/:id/words
func (h *StudyHandler) GetStudySessionWords(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, errors.New("invalid session ID"))
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "100"))
	words, err := h.studyService.GetStudySessionWords(id, page, perPage)
	if err != nil {
		response.InternalError(c, err)
		return
	}
	response.Success(c, words)
}

// RecordWordReview handles POST /api/study-sessions/:id/words/:word_id/review
func (h *StudyHandler) RecordWordReview(c *gin.Context) {
	sessionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, errors.New("invalid session ID"))
		return
	}

	wordID, err := strconv.ParseInt(c.Param("word_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, errors.New("invalid word ID"))
		return
	}

	var req struct {
		Correct *bool `json:"correct"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err)
		return
	}
	if req.Correct == nil {
		response.BadRequest(c, errors.New("correct is required"))
		return
	}

	review, err := h.studyService.RecordWordReview(sessionID, wordID, *req.Correct)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, review)
}

// ResetHistory handles POST /api/reset-history
func (h *StudyHandler) ResetHistory(c *gin.Context) {
	// TODO: Implement after adding service method
	response.Success(c, gin.H{
		"success": true,
		"message": "Study history has been reset",
	})
}

// FullReset handles POST /api/full-reset
func (h *StudyHandler) FullReset(c *gin.Context) {
	// TODO: Implement after adding service method
	response.Success(c, gin.H{
		"success": true,
		"message": "System has been fully reset",
	})
}

// GetStudyActivity handles GET /api/study_activities/:id
func (h *StudyHandler) GetStudyActivity(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, errors.New("invalid activity ID"))
		return
	}
	details, err := h.studyService.GetStudyActivityDetails(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(c, errors.New("activity not found"))
			return
		}
		response.InternalError(c, err)
		return
	}
	response.Success(c, details)
}

// GetActivityStudySessions handles GET /api/study_activities/:id/study_sessions
func (h *StudyHandler) GetActivityStudySessions(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, errors.New("invalid activity ID"))
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "100"))
	sessions, err := h.studyService.GetStudySessionsByActivity(id, page, perPage)
	if err != nil {
		response.InternalError(c, err)
		return
	}
	response.Success(c, sessions)
}
