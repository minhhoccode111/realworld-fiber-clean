package response

import "github.com/minhhoccode111/realworld-fiber-clean/internal/entity"

// TagsResponse -.
type TagsResponse struct {
	Tags []entity.Tag `json:"tags"`

	entity.Pagination
}
