package middleware

import (
	"fmt"
	"net/http"

	"github.com/BarTar213/bartlomiej-tarczynski/models"
	"github.com/gin-gonic/gin"
)

func CheckContentLength(maxContentLength int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxContentLength {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, models.Response{Error: fmt.Sprintf("content length exceed limit %d bytes", maxContentLength)})
			return
		}
		c.Next()
	}
}
