package response

import "github.com/minhhoccode111/realworld-fiber-clean/internal/entity"

// ArticleDetailResponse wraps a single article detail.
type ArticleDetailResponse struct {
	Article *entity.ArticleDetail `json:"article"`
}

// ArticlePreviewsResponse wraps a list of article previews with pagination.
type ArticlePreviewsResponse struct {
	Articles []entity.ArticlePreview `json:"articles"`

	entity.Pagination
}
