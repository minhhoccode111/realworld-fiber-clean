package v1

import (
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http/v1/request"
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
	var body request.ArticleCreateRequest

	if err := ctx.BodyParser(&body); err != nil {
		r.l.Error(err, "http - v1 - postCreateArticle - ctx.BodyParser")

		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	body.Article.Trim()

	if err := r.v.Struct(body.Article); err != nil {
		r.l.Error(err, "http - v1 - postCreateArticle - r.v.Struct")
		if verrs, ok := err.(validator.ValidationErrors); ok {
			errors := make([]string, 0, len(verrs))
			for _, e := range verrs {
				switch e.Tag() {
				case "required":
					errors = append(errors, e.Field()+" is required")
				default:
					errors = append(errors, e.Field()+" is invalid")
				}
			}
			return errorResponse(ctx, http.StatusBadRequest, strings.Join(errors, "; "))
		}
		return errorResponse(ctx, http.StatusInternalServerError, "validation error")
	}

	return nil
}
