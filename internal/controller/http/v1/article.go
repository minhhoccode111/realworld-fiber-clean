package v1

import (
	"github.com/gofiber/fiber/v2"
)

// @Summary     Create Article
// @Description Create Article
// @ID          articles-create
// @Tags  	    articles
// @Accept      json
// @Produce     json
// @Param       request body request.ArticleCreateRequest true "Create Article"
// @Success     200 {object} response.ArticleDetailResponse
// @Failure     400 {object} response.Error
// @Failure     401 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /articles [post]
func (r *V1) postCreateArticle(ctx *fiber.Ctx) error {
	return nil
}
