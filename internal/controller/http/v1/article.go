package v1

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http/middleware"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http/v1/request"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http/v1/response"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/entity"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/utils"
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
func (r *V1) postArticle(c *gin.Context) {
	var body request.ArticleCreateRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		r.l.Error(err, "http - v1 - postCreateArticle - c.ShouldBindJSON")

		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	body.Article.Trim()

	if err := r.v.Struct(body.Article); err != nil {
		r.l.Error(err, "http - v1 - postCreateArticle - r.v.Struct")

		errs := validatorx.ExtractErrors(err)
		errorResponse(c, http.StatusBadRequest, strings.Join(errs, "; "))
		return
	}

	userID := c.MustGet(string(middleware.CtxUserIDKey)).(string)
	if userID == "" {
		errorResponse(c, http.StatusUnauthorized, "cannot authorize user in jwt")
		return
	}

	a, err := r.a.Create(c.Request.Context(), &entity.Article{
		AuthorID:    userID,
		Title:       body.Article.Title,
		Body:        body.Article.Body,
		Description: body.Article.Description,
	}, body.Article.TagList)
	if err != nil {
		r.l.Error(err, "http - v1 - postCreateArticle - r.a.Create")

		errorResponse(c, http.StatusInternalServerError, "database problems")
		return
	}

	c.JSON(http.StatusCreated, response.ArticleDetailResponse{
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
func (r *V1) getAllArticles(c *gin.Context) {
	isAuth := c.MustGet(string(middleware.CtxIsAuthKey)).(bool)

	userID := ""
	if isAuth {
		userID = c.MustGet(string(middleware.CtxUserIDKey)).(string)
		if userID == "" {
			errorResponse(c, http.StatusUnauthorized, "cannot authorize user in jwt")
			return
		}
	}

	tag, author, favorited, limit, offset := utils.SearchQueries(c) // Will need to update SearchQueries to take *gin.Context

	articles, total, err := r.a.List(
		c.Request.Context(),
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

		errorResponse(c, http.StatusInternalServerError, "database problems")
		return
	}

	c.JSON(http.StatusOK, response.ArticlePreviewsResponse{
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
func (r *V1) getFeedArticles(c *gin.Context) {
	userID := c.MustGet(string(middleware.CtxUserIDKey)).(string)
	if userID == "" {
		errorResponse(c, http.StatusUnauthorized, "cannot authorize user in jwt")
		return
	}

	_, _, _, limit, offset := utils.SearchQueries(c) // Will need to update SearchQueries to take *gin.Context

	articles, total, err := r.a.List(
		c.Request.Context(),
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

		errorResponse(c, http.StatusInternalServerError, "database problems")
		return
	}

	c.JSON(http.StatusOK, response.ArticlePreviewsResponse{
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
func (r *V1) getArticle(c *gin.Context) {
	isAuth := c.MustGet(string(middleware.CtxIsAuthKey)).(bool)

	userID := ""
	if isAuth {
		userID = c.MustGet(string(middleware.CtxUserIDKey)).(string)
		if userID == "" {
			errorResponse(c, http.StatusUnauthorized, "cannot authorize user in jwt")
			return
		}
	}

	slug := c.Param("slug")
	if slug == "" {
		errorResponse(c, http.StatusBadRequest, "slug is required")
		return
	}

	a, err := r.a.Detail(c.Request.Context(), userID, slug)
	if err != nil {
		if errors.Is(err, entity.ErrNoRows) {
			errorResponse(c, http.StatusNotFound, "Article not found")
			return
		}

		r.l.Error(err, "http - v1 - getArticle - r.a.Detail")

		errorResponse(c, http.StatusInternalServerError, "database problems")
		return
	}

	c.JSON(http.StatusOK, response.ArticleDetailResponse{
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
func (r *V1) putArticle(c *gin.Context) {
	var body request.ArticleUpdateRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		r.l.Error(err, "http - v1 - putArticle - c.ShouldBindJSON")

		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	body.Article.Trim()

	if err := r.v.Struct(body.Article); err != nil {
		r.l.Error(err, "http - v1 - putArticle - r.v.Struct")

		errs := validatorx.ExtractErrors(err)
		errorResponse(c, http.StatusBadRequest, strings.Join(errs, "; "))
		return
	}

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

	a, err := r.a.Update(c.Request.Context(), userID, slug, &entity.Article{
		Title:       body.Article.Title,
		Description: body.Article.Description,
		Body:        body.Article.Body,
	})
	if err != nil {
		if errors.Is(err, entity.ErrForbidden) {
			errorResponse(c, http.StatusForbidden, "Only article author can update it")
			return
		}

		if errors.Is(err, entity.ErrNoRows) {
			errorResponse(c, http.StatusNotFound, "Article not found")
			return
		}

		r.l.Error(err, "http - v1 - putArticle - r.a.Update")

		errorResponse(c, http.StatusInternalServerError, "database problems")
		return
	}

	c.JSON(http.StatusOK, response.ArticleDetailResponse{
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
func (r *V1) deleteArticle(c *gin.Context) {
	userID := c.MustGet(string(middleware.CtxUserIDKey)).(string)
	if userID == "" {
		errorResponse(c, http.StatusUnauthorized, "cannot authorize user in jwt")
		return
	}

	userRole, ok := c.Get(string(middleware.CtxUserRoleKey))
	if !ok {
		userRole = entity.UserRole // Default to UserRole if not found
	}

	slug := c.Param("slug")
	if slug == "" {
		errorResponse(c, http.StatusBadRequest, "slug is required")
		return
	}

	err := r.a.Delete(c.Request.Context(), userID, slug, userRole.(entity.Role))
	if err != nil {
		if errors.Is(err, entity.ErrForbidden) {
			errorResponse(
				c,
				http.StatusForbidden,
				"Only admin/author can delete this article",
			)
			return
		}

		if errors.Is(err, entity.ErrNoEffect) {
			errorResponse(c, http.StatusNotFound, "Article not found")
			return
		}

		r.l.Error(err, "http - v1 - deleteArticle - r.a.Delete")

		errorResponse(c, http.StatusInternalServerError, "database problems")
		return
	}

	c.Status(http.StatusNoContent)
}
