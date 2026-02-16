package health

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-service-template/internal/models"
	"go-service-template/internal/response"
)

var _ = models.Response{} // used by swag @Success annotation

// Handler handles health endpoints.
type Handler struct{}

// New creates a new health handler.
func New() *Handler {
	return &Handler{}
}

// Get returns the health status.
// Get godoc
// @Summary      Health check
// @Description  Returns service health status
// @Tags         health
// @Produce      json
// @Success      200  {object}  models.Response
// @Router       /health [get]
func (h *Handler) Get(c *gin.Context) {
	response.JSONSuccess(c, http.StatusOK, OKResponse())
}
