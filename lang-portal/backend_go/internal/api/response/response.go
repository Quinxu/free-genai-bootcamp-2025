package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Error response structure
type ErrorResponse struct {
	Error string `json:"error"`
}

// Success sends a successful JSON response
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

// Error sends an error JSON response
func Error(c *gin.Context, status int, err error) {
	c.JSON(status, ErrorResponse{Error: err.Error()})
}

// BadRequest sends a 400 bad request response
func BadRequest(c *gin.Context, err error) {
	Error(c, http.StatusBadRequest, err)
}

// NotFound sends a 404 not found response
func NotFound(c *gin.Context, err error) {
	Error(c, http.StatusNotFound, err)
}

// InternalError sends a 500 internal server error response
func InternalError(c *gin.Context, err error) {
	Error(c, http.StatusInternalServerError, err)
}
