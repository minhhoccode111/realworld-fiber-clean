package response

import "github.com/minhhoccode111/realworld-fiber-clean/internal/entity"

// CommentDetailResponse wraps a single comment detail.
type CommentDetailResponse struct {
	Comment *entity.CommentDetail `json:"comment"`
}

// CommentDetailsResponse wraps a list of comment details with pagination.
type CommentDetailsResponse struct {
	Comments []entity.CommentDetail `json:"comments"`

	entity.Pagination
}
