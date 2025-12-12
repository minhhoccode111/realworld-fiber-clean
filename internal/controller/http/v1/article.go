package v1

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http/common"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http/v1/request"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http/v1/response"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/entity"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/utilities"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/validatorx"
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
func (r *V1) postArticle(ctx *fiber.Ctx) error {
	var body request.ArticleCreateRequest

	if err := ctx.BodyParser(&body); err != nil {
		r.l.Error(err, "http - v1 - postCreateArticle - ctx.BodyParser")

		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	body.Article.Trim()

	if err := r.v.Struct(body.Article); err != nil {
		r.l.Error(err, "http - v1 - postCreateArticle - r.v.Struct")

		errs := validatorx.ExtractErrors(err)

		return errorResponse(ctx, http.StatusBadRequest, strings.Join(errs, "; "))
	}

	userID, ok := ctx.Locals(common.CtxUserIDKey).(string)
	if !ok || userID == "" {
		return errorResponse(ctx, http.StatusUnauthorized, "cannot authorize user in jwt")
	}

	a, err := r.a.Create(ctx.UserContext(), &entity.Article{
		AuthorID:    userID,
		Title:       body.Article.Title,
		Body:        body.Article.Body,
		Description: body.Article.Description,
	}, body.Article.TagList)
	if err != nil {
		r.l.Error(err, "http - v1 - postCreateArticle - r.a.Create")

		return errorResponse(ctx, http.StatusInternalServerError, "database problems")
	}

	return ctx.Status(http.StatusCreated).JSON(response.ArticleDetailResponse{
		Article: a,
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
	isAuth, ok := ctx.Locals(common.CtxIsAuthKey).(bool)
	if !ok {
		isAuth = false
	}

	userID, ok := ctx.Locals(common.CtxUserIDKey).(string)
	if !ok {
		userID = ""
	}

	if userID == "" && isAuth {
		return errorResponse(ctx, http.StatusUnauthorized, "cannot authorize user in jwt")
	}

	tag, author, favorited, limit, offset := utilities.SearchQueries(ctx)

	articles, total, err := r.a.List(
		ctx.UserContext(),
		false,
		userID,
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
	userID := ctx.Locals(common.CtxUserIDKey).(string)
	if userID == "" {
		return errorResponse(ctx, http.StatusUnauthorized, "cannot authorize user in jwt")
	}

	_, _, _, limit, offset := utilities.SearchQueries(ctx)

	articles, total, err := r.a.List(
		ctx.UserContext(),
		true,
		userID,
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

// @Summary     Get article
// @Description Get article by slug
// @ID          articles-get-by-slug
// @Tags        articles
// @Produce     json
// @Param       slug path string true "Article slug"
// @Success     200 {object} response.ArticleDetailResponse
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /articles/{slug} [get]
// @Security    BearerAuth
func (r *V1) getArticle(ctx *fiber.Ctx) error {
	isAuth := ctx.Locals(common.CtxIsAuthKey).(bool)

	userID := ctx.Locals(common.CtxUserIDKey).(string)
	if userID == "" && isAuth {
		return errorResponse(ctx, http.StatusUnauthorized, "cannot authorize user in jwt")
	}

	slug := ctx.Params("slug")
	if slug == "" {
		return errorResponse(ctx, http.StatusBadRequest, "slug is required")
	}

	a, err := r.a.Detail(ctx.UserContext(), userID, slug)
	if err != nil {
		if errors.Is(err, entity.ErrNoRows) {
			return errorResponse(ctx, http.StatusNotFound, "Article not found")
		}

		r.l.Error(err, "http - v1 - getArticle - r.a.Detail")

		return errorResponse(ctx, http.StatusInternalServerError, "database problems")
	}

	return ctx.Status(http.StatusOK).JSON(response.ArticleDetailResponse{
		Article: a,
	})
}

// @Summary     Put article
// @Description Put article by slug
// @ID          articles-put-by-slug
// @Tags        articles
// @Produce     json
// @Param       slug path string true "Article slug"
// @Param       request body request.ArticleUpdateRequest true "Update Article"
// @Success     200 {object} response.ArticleDetailResponse
// @Failure     400 {object} response.Error
// @Failure     401 {object} response.Error
// @Failure     403 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /articles/{slug} [put]
// @Security    BearerAuth
func (r *V1) putArticle(ctx *fiber.Ctx) error {
	var body request.ArticleUpdateRequest

	if err := ctx.BodyParser(&body); err != nil {
		r.l.Error(err, "http - v1 - putArticle - ctx.BodyParser")

		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	body.Article.Trim()

	if err := r.v.Struct(body.Article); err != nil {
		r.l.Error(err, "http - v1 - putArticle - r.v.Struct")

		errs := validatorx.ExtractErrors(err)

		return errorResponse(ctx, http.StatusBadRequest, strings.Join(errs, "; "))
	}

	userID := ctx.Locals(common.CtxUserIDKey).(string)
	if userID == "" {
		return errorResponse(ctx, http.StatusUnauthorized, "cannot authorize user in jwt")
	}

	slug := ctx.Params("slug")
	if slug == "" {
		return errorResponse(ctx, http.StatusBadRequest, "slug is required")
	}

	a, err := r.a.Update(ctx.UserContext(), userID, slug, &entity.Article{
		Title:       body.Article.Title,
		Description: body.Article.Description,
		Body:        body.Article.Body,
	})
	if err != nil {
		if errors.Is(err, entity.ErrForbidden) {
			return errorResponse(ctx, http.StatusForbidden, "Only article author can update it")
		}

		if errors.Is(err, entity.ErrNoRows) {
			return errorResponse(ctx, http.StatusNotFound, "Article not found")
		}

		r.l.Error(err, "http - v1 - putArticle - r.a.Update")

		return errorResponse(ctx, http.StatusInternalServerError, "database problems")
	}

	return ctx.Status(http.StatusOK).JSON(response.ArticleDetailResponse{
		Article: a,
	})
}

// @Summary     Delete article
// @Description Delete article by slug
// @ID          articles-delete-by-slug
// @Tags        articles
// @Produce     json
// @Param       slug path string true "Article slug"
// @Success     204 "No Content"
// @Failure     400 {object} response.Error
// @Failure     401 {object} response.Error
// @Failure     403 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /articles/{slug} [delete]
// @Security    BearerAuth
func (r *V1) deleteArticle(ctx *fiber.Ctx) error {
	userID := ctx.Locals(common.CtxUserIDKey).(string)
	if userID == "" {
		return errorResponse(ctx, http.StatusUnauthorized, "cannot authorize user in jwt")
	}

	userRole, ok := ctx.Locals(common.CtxUserRoleKey).(entity.Role)
	if !ok {
		userRole = entity.UserRole
	}

	slug := ctx.Params("slug")
	if slug == "" {
		return errorResponse(ctx, http.StatusBadRequest, "slug is required")
	}

	err := r.a.Delete(ctx.UserContext(), userID, slug, userRole)
	if err != nil {
		if errors.Is(err, entity.ErrForbidden) {
			return errorResponse(
				ctx,
				http.StatusForbidden,
				"Only admin/author can delete this article",
			)
		}

		if errors.Is(err, entity.ErrNoEffect) {
			return errorResponse(ctx, http.StatusNotFound, "Article not found")
		}

		r.l.Error(err, "http - v1 - deleteArticle - r.a.Delete")

		return errorResponse(ctx, http.StatusInternalServerError, "database problems")
	}

	return ctx.SendStatus(http.StatusNoContent)
}
