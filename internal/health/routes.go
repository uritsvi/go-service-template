package health

import (
	"github.com/gin-gonic/gin"
)

// Register mounts health routes at /health.
func Register(router *gin.Engine, h *Handler) {
	g := router.Group("/health")
	g.GET("", h.Get)
}
