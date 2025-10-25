package util

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// TODO: tag can be a slice of strings
func SearchQueries(ctx *fiber.Ctx) (tag, author, favorited string, limit, offset uint64) {
	var err error
	tag = strings.TrimSpace(ctx.Query("tag"))
	author = strings.TrimSpace(ctx.Query("author"))
	favorited = strings.TrimSpace(ctx.Query("favorited"))

	limit, err = strconv.ParseUint(ctx.Query("limit", "10"), 10, 64)
	if err != nil {
		limit = 10
	}

	offset, err = strconv.ParseUint(ctx.Query("offset", "0"), 10, 64)
	if err != nil {
		offset = 0
	}

	return tag, author, favorited, limit, offset
}
