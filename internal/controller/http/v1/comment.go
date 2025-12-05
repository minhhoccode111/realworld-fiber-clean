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

// @Summary     Create Comment
// @Description Create Comment
// @ID          comments-create
// @Tags  	    comments
// @Accept      json
// @Produce     json
// @Param       request body request.CommentCreateRequest true "Create Comment"
// @Success     201 {object} response.CommentDetailResponse
// @Failure     400 {object} response.Error
// @Failure     401 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /articles/{slug}/comments [post]
// @Security    BearerAuth
func (r *V1) postComment(c *gin.Context) {
	var body request.CommentCreateRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		r.l.Error(err, "http - v1 - postCreateComment - c.ShouldBindJSON")

		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	body.Comment.Trim()

	if err := r.v.Struct(body.Comment); err != nil {
		r.l.Error(err, "http - v1 - postCreateComment - r.v.Struct")

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

	comment, err := r.c.Create(
		c.Request.Context(),
		slug,
		&entity.Comment{
			AuthorID: userID,
			Body:     body.Comment.Body,
		},
	)
	if err != nil {
		if errors.Is(err, entity.ErrNoRows) {
			errorResponse(c, http.StatusNotFound, "Article not found")
			return
		}

		r.l.Error(err, "http - v1 - postCreateComment - r.c.Create")

		errorResponse(c, http.StatusInternalServerError, "database problems")
		return
	}

	c.JSON(http.StatusCreated, response.CommentDetailResponse{
		Comment: comment,
	})
}

// @Summary     Get all comments
// @Description Get all comments of an article
// @ID          comments-get-all
// @Tags        comments
// @Produce     json
// @Param       limit      query uint64 false "Limit number of results"
// @Param       offset     query uint64 false "Offset for pagination"
// @Success     200 {object} response.CommentDetailsResponse
// @Failure     400 {object} response.Error
// @Failure     401 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /articles/{slug}/comments [get]
// @Security    BearerAuth
func (r *V1) getAllComments(c *gin.Context) {
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

	_, _, _, limit, offset := utils.SearchQueries(c)

	comments, total, err := r.c.List(
		c.Request.Context(),
		userID,
		slug,
		limit,
		offset,
	)
	if err != nil {
		r.l.Error(err, "http - v1 - getAllComments - r.c.List")

		errorResponse(c, http.StatusInternalServerError, "database problems")
		return
	}

	c.JSON(http.StatusOK, response.CommentDetailsResponse{
		Comments: comments,
		Pagination: entity.Pagination{
			Total:  total,
			Limit:  limit,
			Offset: offset,
		},
	})
}

// @Summary     Delete comment
// @Description Delete comment by id
// @ID          comment-delete-by-id
// @Tags        comments
// @Produce     json
// @Param       slug path string true "Article slug"
// @Param       commentID path string true "Comment ID"
// @Success     204 "No Content"
// @Failure     400 {object} response.Error
// @Failure     401 {object} response.Error
// @Failure     403 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /articles/{slug}/comments/{commentID} [delete]
// @Security    BearerAuth
func (r *V1) deleteComment(c *gin.Context) {
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

	commentID := c.Param("commentID")
	if commentID == "" {
		errorResponse(c, http.StatusBadRequest, "commentID is required")
		return
	}

	err := r.c.Delete(c.Request.Context(), userID, slug, commentID, userRole.(entity.Role))
	if err != nil {
		if errors.Is(err, entity.ErrForbidden) {
			errorResponse(
				c,
				http.StatusForbidden,
				"Only admin/author can delete this comment",
			)
			return
		}

		if errors.Is(err, entity.ErrNoEffect) {
			errorResponse(c, http.StatusNotFound, "Article/comment not found")
			return
		}

		r.l.Error(err, "http - v1 - deleteComment - r.c.Delete")

		errorResponse(c, http.StatusInternalServerError, "database problems")
		return
	}

	c.Status(http.StatusNoContent)
}
