package v1

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http/v1/response"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/entity"
)

// @Summary     Get tags
// @Description Get all tags of all articles with pagination
// @ID          tags
// @Tags  	    tags
// @Produce     json
// @Param       limit  query int false "Number of items to return"  minimum(1)  default(10)
// @Param       offset query int false "Number of items to skip"    minimum(0)  default(0)
// @Success     200 {object} response.TagsResponse
// @Failure     500 {object} response.Error
// @Router      /tags [get]
func (r *V1) getTags(ctx *fiber.Ctx) error {
	limit, err := strconv.ParseUint(ctx.Query("limit", "10"), 10, 64)
	if err != nil {
		limit = 10
	}

	offset, err := strconv.ParseUint(ctx.Query("offset", "0"), 10, 64)
	if err != nil {
		offset = 0
	}

	tags, total, err := r.tag.List(ctx.UserContext(), limit, offset)
	if err != nil {
		r.l.Error(err, "http - v1 - getTags")

		return errorResponse(ctx, http.StatusInternalServerError, "database problems")
	}

	return ctx.Status(http.StatusOK).JSON(response.TagsResponse{
		Tags: tags,
		Pagination: entity.Pagination{
			Limit:  limit,
			Offset: offset,
			Total:  total,
		},
	})
}
