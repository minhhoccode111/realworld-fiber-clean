package utils

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// TODO: tag can be a slice of strings.
func SearchQueries(c *gin.Context) (tag, author, favorited string, limit, offset uint64) {
	var err error

	tag = strings.TrimSpace(c.Query("tag"))
	author = strings.TrimSpace(c.Query("author"))
	favorited = strings.TrimSpace(c.Query("favorited"))

	limit, err = strconv.ParseUint(c.DefaultQuery("limit", "10"), 10, 64)
	if err != nil || limit > 200 {
		limit = 10
	}

	offset, err = strconv.ParseUint(c.DefaultQuery("offset", "0"), 10, 64)
	if err != nil {
		offset = 0
	}

	return tag, author, favorited, limit, offset
}
