package response

import "github.com/minhhoccode111/realworld-fiber-clean/internal/entity"

// Tags -.
type Tags struct {
	Tags []entity.Tag `json:"tags"`

	entity.Pagination
}
