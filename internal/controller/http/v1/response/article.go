package response

import "github.com/minhhoccode111/realworld-fiber-clean/internal/entity"

type ArticleDetailResponse struct {
	Article entity.ArticleDetail `json:"article"`
}

type ArticlePreviewsResponse struct {
	Articles []entity.ArticlePreview `json:"articles"`

	entity.Pagination
}
