package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http/v1/response"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/entity"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/utils"
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
func (r *V1) getTags(c *gin.Context) {
	_, _, _, limit, offset := utils.SearchQueries(c)

	tags, total, err := r.tag.List(c.Request.Context(), limit, offset)
	if err != nil {
		r.l.Error(err, "http - v1 - getTags")

		errorResponse(c, http.StatusInternalServerError, "database problems")
		return
	}

	c.JSON(http.StatusOK, response.TagsResponse{
		Tags: tags,
		Pagination: entity.Pagination{
			Limit:  limit,
			Offset: offset,
			Total:  total,
		},
	})
}
