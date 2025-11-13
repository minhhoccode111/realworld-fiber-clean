package v1

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http/middleware"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http/v1/response"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/entity"
)

// @Summary     Favorite Article
// @Description Favorite an article
// @ID          favorites-create
// @Tags  	    favorites
// @Accept      json
// @Produce     json
// @Param       slug path string true "Article slug"
// @Success     200 {object} response.ArticleDetailResponse
// @Failure     400 {object} response.Error
// @Failure     401 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /articles/{slug}/favorite [post]
// @Security    BearerAuth
func (r *V1) createFavorite(ctx *fiber.Ctx) error {
	userId := ctx.Locals(middleware.CtxUserIdKey).(string)
	if userId == "" {
		return errorResponse(ctx, http.StatusUnauthorized, "cannot authorize user in jwt")
	}

	slug := ctx.Params("slug")
	if slug == "" {
		return errorResponse(ctx, http.StatusBadRequest, "slug is required")
	}

	err := r.f.Create(ctx.UserContext(), userId, slug)
	if err != nil {
		if errors.Is(err, entity.ZeroRowsAffected) {
			return errorResponse(ctx, http.StatusBadRequest, "Article is already favorited")
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23502" { // not null constraint
			return errorResponse(ctx, http.StatusNotFound, "Article not found")
		}

		r.l.Error(err, "http - v1 - createFavorite - r.f.Create")

		return errorResponse(ctx, http.StatusInternalServerError, "database problems")
	}

	article, err := r.a.Detail(ctx.UserContext(), userId, slug)
	if err != nil {
		if errors.Is(err, entity.ErrNoRows) {
			return errorResponse(ctx, http.StatusNotFound, "Article not found")
		}

		r.l.Error(err, "http - v1 - deleteFavorite - r.f.Delete")

		return errorResponse(ctx, http.StatusInternalServerError, "database problems")
	}

	return ctx.Status(http.StatusOK).JSON(response.ArticleDetailResponse{
		Article: article,
	})
}

// @Summary     Unfavorite Article
// @Description Unfavorite an article
// @ID          favorites-delete
// @Tags        favorites
// @Produce     json
// @Param       slug path string true "Article slug"
// @Success     200 {object} response.ArticleDetailResponse
// @Failure     400 {object} response.Error
// @Failure     401 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /articles/{slug}/favorite [delete]
// @Security    BearerAuth
func (r *V1) deleteFavorite(ctx *fiber.Ctx) error {
	userId := ctx.Locals(middleware.CtxUserIdKey).(string)
	if userId == "" {
		return errorResponse(ctx, http.StatusUnauthorized, "cannot authorize user in jwt")
	}

	slug := ctx.Params("slug")
	if slug == "" {
		return errorResponse(ctx, http.StatusBadRequest, "slug is required")
	}

	err := r.f.Delete(ctx.UserContext(), userId, slug)
	if err != nil {
		if errors.Is(err, entity.ZeroRowsAffected) {
			return errorResponse(ctx, http.StatusBadRequest, "Article is already unfavorited")
		}

		r.l.Error(err, "http - v1 - deleteFavorite - r.c.Delete")

		return errorResponse(ctx, http.StatusInternalServerError, "database problems")
	}

	article, err := r.a.Detail(ctx.UserContext(), userId, slug)
	if err != nil {
		if errors.Is(err, entity.ErrNoRows) {
			return errorResponse(ctx, http.StatusNotFound, "Article not found")
		}

		r.l.Error(err, "http - v1 - deleteFavorite - r.a.Detail")

		return errorResponse(ctx, http.StatusInternalServerError, "database problems")
	}

	return ctx.Status(http.StatusOK).JSON(response.ArticleDetailResponse{
		Article: article,
	})
}
