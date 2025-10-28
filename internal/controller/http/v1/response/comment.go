package response

import "github.com/minhhoccode111/realworld-fiber-clean/internal/entity"

type CommentDetailResponse struct {
	Comment entity.CommentDetail `json:"comment"`
}

type CommentDetailsResponse struct {
	Comments []entity.CommentDetail `json:"comments"`

	entity.Pagination
}
