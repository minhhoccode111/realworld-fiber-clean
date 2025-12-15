package response

import "github.com/minhhoccode111/realworld-fiber-clean/internal/entity"

// TagsResponse wraps tag names with pagination metadata.
type TagsResponse struct {
	Tags []entity.TagName `json:"tags"`

	entity.Pagination
}
