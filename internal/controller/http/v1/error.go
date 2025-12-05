package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http/v1/response"
)

func errorResponse(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(code, response.Error{Error: msg})
}
