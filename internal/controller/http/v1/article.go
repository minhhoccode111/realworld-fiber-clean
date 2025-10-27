package v1

import (
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http/middleware"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http/v1/request"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http/v1/response"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/entity"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/util"
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
// @Security    BearerAuth
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

	userId := ctx.Locals(middleware.CtxUserIdKey).(string)
	if userId == "" {
		return errorResponse(ctx, http.StatusUnauthorized, "cannot authorize user in jwt")
	}

	article, err := r.a.Create(ctx.UserContext(), entity.Article{
		AuthorId:    userId,
		Title:       body.Article.Title,
		Body:        body.Article.Body,
		Description: body.Article.Description,
	}, body.Article.TagList)
	if err != nil {
		r.l.Error(err, "http - v1 - postCreateArticle - r.a.Create")

		return errorResponse(ctx, http.StatusInternalServerError, "database problems")
	}

	return ctx.Status(http.StatusCreated).JSON(response.ArticleDetailResponse{
		Article: article,
	})
}

// @Summary     Get all articles
// @Description Get all articles (filter by author, favorited, tag)
// @ID          articles-get-all
// @Tags        articles
// @Produce     json
// @Param       limit      query uint64 false "Limit number of results"
// @Param       offset     query uint64 false "Offset for pagination"
// @Param       author     query string false "Filter by author username"
// @Param       favorited  query string false "Filter by favorited username"
// @Param       tag        query string false "Filter by tag"
// @Success     200 {object} response.ArticlePreviewsResponse
// @Failure     500 {object} response.Error
// @Router      /articles [get]
// @Security    BearerAuth
func (r *V1) getAllArticles(ctx *fiber.Ctx) error {
	isAuth := ctx.Locals(middleware.CtxIsAuthKey).(bool)
	userId := ctx.Locals(middleware.CtxUserIdKey).(string)
	if userId == "" && isAuth {
		return errorResponse(ctx, http.StatusUnauthorized, "cannot authorize user in jwt")
	}

	tag, author, favorited, limit, offset := util.SearchQueries(ctx)

	articles, total, err := r.a.List(
		ctx.UserContext(),
		false,
		userId,
		tag,
		author,
		favorited,
		limit,
		offset,
	)
	if err != nil {
		r.l.Error(err, "http - v1 - getAllArticles - r.a.List")

		return errorResponse(ctx, http.StatusInternalServerError, "database problems")
	}

	return ctx.Status(http.StatusOK).JSON(response.ArticlePreviewsResponse{
		Articles: articles,
		Pagination: entity.Pagination{
			Total:  total,
			Limit:  limit,
			Offset: offset,
		},
	})
}

// @Summary     Get feed articles
// @Description Get feed articles (from followed authors)
// @ID          articles-get-feed
// @Tags        articles
// @Produce     json
// @Param       limit      query uint64 false "Limit number of results"
// @Param       offset     query uint64 false "Offset for pagination"
// @Success     200 {object} response.ArticlePreviewsResponse
// @Failure     500 {object} response.Error
// @Router      /articles/feed [get]
// @Security    BearerAuth
func (r *V1) getFeedArticles(ctx *fiber.Ctx) error {
	userId, ok := ctx.Locals(middleware.CtxUserIdKey).(string)
	if !ok {
		return errorResponse(ctx, http.StatusUnauthorized, "cannot authorize user in jwt")
	}

	_, _, _, limit, offset := util.SearchQueries(ctx)

	articles, total, err := r.a.List(
		ctx.UserContext(),
		true,
		userId,
		"",
		"",
		"",
		limit,
		offset,
	)
	if err != nil {
		r.l.Error(err, "http - v1 - getAllArticles - r.a.List")

		return errorResponse(ctx, http.StatusInternalServerError, "database problems")
	}

	return ctx.Status(http.StatusOK).JSON(response.ArticlePreviewsResponse{
		Articles: articles,
		Pagination: entity.Pagination{
			Total:  total,
			Limit:  limit,
			Offset: offset,
		},
	})
}
