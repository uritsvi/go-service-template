package health

import (
	"github.com/gin-gonic/gin"

	healthconfig "go-service-template/internal/health/config"
)

// Register mounts health routes on the given router group or engine.
// cfg can be nil; if provided, BasePath is used for the group path.
func Register(router *gin.Engine, h *Handler, cfg *healthconfig.Config) {
	g := router.Group(cfg.BasePath)
	{
		g.GET("", h.Get)
	}
}
