package response

import "github.com/minhhoccode111/realworld-fiber-clean/internal/entity"

// TagsResponse -.
type TagsResponse struct {
	Tags []entity.TagName `json:"tags"`

	entity.Pagination
}
