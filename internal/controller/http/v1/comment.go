package v1

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http/httpmeta"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http/v1/request"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http/v1/response"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/entity"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/utilities"
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
func (r *V1) postComment(ctx *fiber.Ctx) error {
	var body request.CommentCreateRequest

	if err := ctx.BodyParser(&body); err != nil {
		r.l.Error(err, "http - v1 - postCreateComment - ctx.BodyParser")

		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	body.Comment.Trim()

	if err := r.v.Struct(body.Comment); err != nil {
		r.l.Error(err, "http - v1 - postCreateComment - r.v.Struct")

		errs := validatorx.ExtractErrors(err)

		return errorResponse(ctx, http.StatusBadRequest, strings.Join(errs, "; "))
	}

	userID, ok := ctx.Locals(httpmeta.CtxUserIDKey).(string)
	if !ok || userID == "" {
		return errorResponse(ctx, http.StatusUnauthorized, "cannot authorize user in jwt")
	}

	slug := ctx.Params("slug")
	if slug == "" {
		return errorResponse(ctx, http.StatusBadRequest, "slug is required")
	}

	c, err := r.c.Create(
		ctx.UserContext(),
		slug,
		&entity.Comment{
			AuthorID: userID,
			Body:     body.Comment.Body,
		},
	)
	if err != nil {
		if errors.Is(err, entity.ErrNoRows) {
			return errorResponse(ctx, http.StatusNotFound, "Article not found")
		}

		r.l.Error(err, "http - v1 - postCreateComment - r.c.Create")

		return errorResponse(ctx, http.StatusInternalServerError, "database problems")
	}

	return ctx.Status(http.StatusCreated).JSON(response.CommentDetailResponse{
		Comment: c,
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
func (r *V1) getAllComments(ctx *fiber.Ctx) error {
	isAuth, ok := ctx.Locals(httpmeta.CtxIsAuthKey).(bool)
	if !ok {
		isAuth = false
	}

	userID, ok := ctx.Locals(httpmeta.CtxUserIDKey).(string)
	if !ok || (userID == "" && isAuth) {
		return errorResponse(ctx, http.StatusUnauthorized, "cannot authorize user in jwt")
	}

	slug := ctx.Params("slug")
	if slug == "" {
		return errorResponse(ctx, http.StatusBadRequest, "slug is required")
	}

	limit, offset := utilities.PaginationQueries(ctx)

	comments, total, err := r.c.List(
		ctx.UserContext(),
		userID,
		slug,
		limit,
		offset,
	)
	if err != nil {
		r.l.Error(err, "http - v1 - getAllComments - r.c.List")

		return errorResponse(ctx, http.StatusInternalServerError, "database problems")
	}

	return ctx.Status(http.StatusOK).JSON(response.CommentDetailsResponse{
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
func (r *V1) deleteComment(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals(httpmeta.CtxUserIDKey).(string)
	if !ok || userID == "" {
		return errorResponse(ctx, http.StatusUnauthorized, "cannot authorize user in jwt")
	}

	userRole, ok := ctx.Locals(httpmeta.CtxUserRoleKey).(entity.Role)
	if !ok {
		userRole = entity.UserRole
	}

	slug := ctx.Params("slug")
	if slug == "" {
		return errorResponse(ctx, http.StatusBadRequest, "slug is required")
	}

	commentID := ctx.Params("commentID")
	if commentID == "" {
		return errorResponse(ctx, http.StatusBadRequest, "commentID is required")
	}

	err := r.c.Delete(ctx.UserContext(), userID, slug, commentID, userRole)
	if err != nil {
		if errors.Is(err, entity.ErrForbidden) {
			return errorResponse(
				ctx,
				http.StatusForbidden,
				"Only admin/author can delete this comment",
			)
		}

		if errors.Is(err, entity.ErrNoEffect) {
			return errorResponse(ctx, http.StatusNotFound, "Article/comment not found")
		}

		r.l.Error(err, "http - v1 - deleteComment - r.c.Delete")

		return errorResponse(ctx, http.StatusInternalServerError, "database problems")
	}

	return ctx.SendStatus(http.StatusNoContent)
}
