package utilities

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// PaginationQueries extracts limit and offset from fiber.Ctx queries for pagination.
func PaginationQueries(ctx *fiber.Ctx) (limit, offset uint64) {
	var err error

	limit, err = strconv.ParseUint(ctx.Query("limit", "10"), 10, 64)
	if err != nil || limit > 200 {
		limit = 10
	}

	offset, err = strconv.ParseUint(ctx.Query("offset", "0"), 10, 64)
	if err != nil {
		offset = 0
	}

	return limit, offset
}
