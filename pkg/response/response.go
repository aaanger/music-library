package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func JSON(c *gin.Context, data any) {
	c.JSON(http.StatusOK, data)
}
