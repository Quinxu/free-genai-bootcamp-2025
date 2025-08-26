package handlers

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	"lang-portal/internal/api/response"
	"lang-portal/internal/service"
)

// GroupHandler handles group-related routes
type GroupHandler struct {
	groupService *service.GroupService
}

// NewGroupHandler creates a new GroupHandler
func NewGroupHandler(groupService *service.GroupService) *GroupHandler {
	return &GroupHandler{groupService: groupService}
}

// RegisterRoutes registers group-related routes
func (h *GroupHandler) RegisterRoutes(r *gin.RouterGroup) {
	groups := r.Group("/groups")
	{
		groups.GET("", h.GetGroups)
		groups.GET("/:id", h.GetGroup)
		groups.GET("/:id/words", h.GetGroupWords)
		groups.GET("/:id/study_sessions", h.GetGroupStudySessions)
	}
}

// GetGroups handles GET /api/groups
func (h *GroupHandler) GetGroups(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "100"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 100
	}

	groups, err := h.groupService.GetGroups(page, perPage)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Success(c, groups)
}

// GetGroup handles GET /api/groups/:id
func (h *GroupHandler) GetGroup(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, errors.New("invalid group ID"))
		return
	}

	group, err := h.groupService.GetGroupByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(c, errors.New("group not found"))
			return
		}
		response.InternalError(c, err)
		return
	}

	response.Success(c, group)
}

// GetGroupWords handles GET /api/groups/:id/words
func (h *GroupHandler) GetGroupWords(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, errors.New("invalid group ID"))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "100"))

	words, err := h.groupService.GetGroupWordsPaginated(id, page, perPage)
	if err != nil {
		response.InternalError(c, err)
		return
	}
	response.Success(c, words)
}

// GetGroupStudySessions handles GET /api/groups/:id/study_sessions
func (h *GroupHandler) GetGroupStudySessions(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, errors.New("invalid group ID"))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "100"))

	sessions, err := h.groupService.GetGroupStudySessions(id, page, perPage)
	if err != nil {
		response.InternalError(c, err)
		return
	}
	response.Success(c, sessions)
}
