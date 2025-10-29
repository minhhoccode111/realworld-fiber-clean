package v1

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http/middleware"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http/v1/response"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/entity"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/util"
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
// @Router      /articles/{username}/comments [post]
// @Security    BearerAuth
func (r *V1) getProfile(ctx *fiber.Ctx) error {
	isAuth := ctx.Locals(middleware.CtxIsAuthKey).(bool)
	userId := ctx.Locals(middleware.CtxUserIdKey).(string)
	if userId == "" && isAuth {
		return errorResponse(ctx, http.StatusUnauthorized, "cannot authorize user in jwt")
	}

	username := ctx.Params("username")
	if username == "" {
		return errorResponse(ctx, http.StatusBadRequest, "username is required")
	}

	profile, err := r.p.Detail(ctx.UserContext(), userId, username)
	if err != nil {
		if strings.Contains(err.Error(), "notfound") {
			return errorResponse(ctx, http.StatusNotFound, "Profile not found")
		}

		r.l.Error(err, "http - v1 - getProfile - r.p.Detail")

		return errorResponse(ctx, http.StatusInternalServerError, "database problems")
	}

	return ctx.Status(200).JSON(response.ProfilePreviewResponse{
		Profile: profile,
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
// @Success     400 {object} response.CommentDetailsResponse
// @Success     401 {object} response.CommentDetailsResponse
// @Failure     500 {object} response.Error
// @Router      /articles/{slug}/comments [get]
// @Security    BearerAuth
func (r *V1) postFollowProfile(ctx *fiber.Ctx) error {
	isAuth := ctx.Locals(middleware.CtxIsAuthKey).(bool)
	userId := ctx.Locals(middleware.CtxUserIdKey).(string)
	if userId == "" && isAuth {
		return errorResponse(ctx, http.StatusUnauthorized, "cannot authorize user in jwt")
	}

	slug := ctx.Params("slug")
	if slug == "" {
		return errorResponse(ctx, http.StatusBadRequest, "slug is required")
	}

	_, _, _, limit, offset := util.SearchQueries(ctx)

	comments, total, err := r.c.List(
		ctx.UserContext(),
		userId,
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
// @Success     204 "No Content"
// @Failure     400 {object} response.Error
// @Failure     401 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /articles/{slug}/comments/{commentId} [delete]
// @Security    BearerAuth
func (r *V1) deleteFollowProfile(ctx *fiber.Ctx) error {
	userId := ctx.Locals(middleware.CtxUserIdKey).(string)
	if userId == "" {
		return errorResponse(ctx, http.StatusUnauthorized, "cannot authorize user in jwt")
	}

	slug := ctx.Params("slug")
	if slug == "" {
		return errorResponse(ctx, http.StatusBadRequest, "slug is required")
	}

	commentId := ctx.Params("commentId")
	if commentId == "" {
		return errorResponse(ctx, http.StatusBadRequest, "commentId is required")
	}

	err := r.c.Delete(ctx.UserContext(), userId, slug, commentId)
	if err != nil {
		if strings.Contains(err.Error(), "notfound") {
			return errorResponse(ctx, http.StatusNotFound, "Article/comment not found")
		}

		r.l.Error(err, "http - v1 - deleteComment - r.c.Delete")

		return errorResponse(ctx, http.StatusInternalServerError, "database problems")
	}

	return ctx.SendStatus(http.StatusNoContent)
}
