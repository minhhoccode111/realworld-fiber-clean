package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
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
func (r *V1) createFavorite(c *gin.Context) {
	userID := c.MustGet(string(middleware.CtxUserIDKey)).(string)
	if userID == "" {
		errorResponse(c, http.StatusUnauthorized, "cannot authorize user in jwt")
		return
	}

	slug := c.Param("slug")
	if slug == "" {
		errorResponse(c, http.StatusBadRequest, "slug is required")
		return
	}

	err := r.f.Create(c.Request.Context(), userID, slug)
	if err != nil {
		if errors.Is(err, entity.ErrNoRows) {
			errorResponse(c, http.StatusNotFound, "Article not found")
			return
		}

		if errors.Is(err, entity.ErrNoEffect) {
			errorResponse(c, http.StatusBadRequest, "Article is already favorited")
			return
		}

		r.l.Error(err, "http - v1 - createFavorite - r.f.Create")

		errorResponse(c, http.StatusInternalServerError, "database problems")
		return
	}

	a, err := r.a.Detail(c.Request.Context(), userID, slug)
	if err != nil {
		if errors.Is(err, entity.ErrNoRows) {
			errorResponse(c, http.StatusNotFound, "Article not found")
			return
		}

		r.l.Error(err, "http - v1 - deleteFavorite - r.f.Delete") // Typo here, should be createFavorite

		errorResponse(c, http.StatusInternalServerError, "database problems")
		return
	}

	c.JSON(http.StatusOK, response.ArticleDetailResponse{
		Article: a,
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
func (r *V1) deleteFavorite(c *gin.Context) {
	userID := c.MustGet(string(middleware.CtxUserIDKey)).(string)
	if userID == "" {
		errorResponse(c, http.StatusUnauthorized, "cannot authorize user in jwt")
		return
	}

	slug := c.Param("slug")
	if slug == "" {
		errorResponse(c, http.StatusBadRequest, "slug is required")
		return
	}

	err := r.f.Delete(c.Request.Context(), userID, slug)
	if err != nil {
		if errors.Is(err, entity.ErrNoEffect) {
			errorResponse(c, http.StatusBadRequest, "Article is already unfavorited")
			return
		}

		r.l.Error(err, "http - v1 - deleteFavorite - r.c.Delete")

		errorResponse(c, http.StatusInternalServerError, "database problems")
		return
	}

	a, err := r.a.Detail(c.Request.Context(), userID, slug)
	if err != nil {
		if errors.Is(err, entity.ErrNoRows) {
			errorResponse(c, http.StatusNotFound, "Article not found")
			return
		}

		r.l.Error(err, "http - v1 - deleteFavorite - r.a.Detail")

		errorResponse(c, http.StatusInternalServerError, "database problems")
		return
	}

	c.JSON(http.StatusOK, response.ArticleDetailResponse{
		Article: a,
	})
}
