package response

import (
	"github.com/gin-gonic/gin"
)

// JSONSuccess sends a successful JSON response with the given status code and data.
func JSONSuccess(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}

// JSONError sends an error JSON response.
func JSONError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}
