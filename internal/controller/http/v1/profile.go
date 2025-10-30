package v1

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http/middleware"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http/v1/response"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/entity"
)

// @Summary     Get profile
// @Description Get profile by username
// @ID          profiles-get-by-username
// @Tags  	    profiles
// @Produce     json
// @Param       username path string true "Username of the profile to get"
// @Success     200 {object} entity.ProfilePreview
// @Failure     400 {object} response.Error
// @Failure     401 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /profiles/{username} [get]
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
		if errors.Is(err, entity.ErrNoRows) {
			return errorResponse(ctx, http.StatusNotFound, "Profile not found")
		}

		r.l.Error(err, "http - v1 - getProfile - r.p.Detail")

		return errorResponse(ctx, http.StatusInternalServerError, "database problems")
	}

	return ctx.Status(200).JSON(response.ProfilePreviewResponse{
		Profile: profile,
	})
}

// @Summary     Follow user
// @Description Follow user by username
// @ID          profiles-follow
// @Tags        profiles
// @Produce     json
// @Param       username path string true "Username of the profile to follow"
// @Success     200 {object} response.ProfilePreviewResponse
// @Failure     400 {object} response.Error
// @Failure     401 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /profiles/{username}/follow [post]
// @Security    BearerAuth
func (r *V1) postFollowProfile(ctx *fiber.Ctx) error {
	userId := ctx.Locals(middleware.CtxUserIdKey).(string)
	if userId == "" {
		return errorResponse(ctx, http.StatusUnauthorized, "cannot authorize user in jwt")
	}

	username := ctx.Params("username")
	if username == "" {
		return errorResponse(ctx, http.StatusBadRequest, "username is required")
	}

	err := r.p.Follow(ctx.UserContext(), userId, username)
	if err != nil {
		if errors.Is(err, entity.ErrNoRows) {
			return errorResponse(ctx, http.StatusNotFound, "Profile not found")
		}

		r.l.Error(err, "http - v1 - postFollowProfile - r.p.Follow")

		return errorResponse(ctx, http.StatusInternalServerError, "database problems")
	}

	profile, err := r.p.Detail(ctx.UserContext(), userId, username)
	if err != nil {
		if errors.Is(err, entity.ErrNoRows) {
			return errorResponse(ctx, http.StatusNotFound, "Profile not found")
		}

		r.l.Error(err, "http - v1 - postFollowProfile - r.p.Detail")

		return errorResponse(ctx, http.StatusInternalServerError, "database problems")
	}

	return ctx.Status(200).JSON(response.ProfilePreviewResponse{
		Profile: profile,
	})
}

// @Summary     Unfollow user
// @Description Unfollow user by username
// @ID          profiles-unfollow
// @Tags        profiles
// @Produce     json
// @Param       username path string true "Username of the profile to unfollow"
// @Success     200 {object} response.ProfilePreviewResponse
// @Failure     400 {object} response.Error
// @Failure     401 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /profiles/{username}/follow [delete]
// @Security    BearerAuth
func (r *V1) deleteFollowProfile(ctx *fiber.Ctx) error {
	userId := ctx.Locals(middleware.CtxUserIdKey).(string)
	if userId == "" {
		return errorResponse(ctx, http.StatusUnauthorized, "cannot authorize user in jwt")
	}

	username := ctx.Params("username")
	if username == "" {
		return errorResponse(ctx, http.StatusBadRequest, "username is required")
	}

	err := r.p.Unfollow(ctx.UserContext(), userId, username)
	if err != nil {
		if errors.Is(err, entity.ErrNoRows) {
			return errorResponse(ctx, http.StatusNotFound, "Profile not found")
		}

		r.l.Error(err, "http - v1 - postFollowProfile - r.p.Follow")

		return errorResponse(ctx, http.StatusInternalServerError, "database problems")
	}

	profile, err := r.p.Detail(ctx.UserContext(), userId, username)
	if err != nil {
		if errors.Is(err, entity.ErrNoRows) {
			return errorResponse(ctx, http.StatusNotFound, "Profile not found")
		}

		r.l.Error(err, "http - v1 - postFollowProfile - r.p.Detail")

		return errorResponse(ctx, http.StatusInternalServerError, "database problems")
	}

	return ctx.Status(200).JSON(response.ProfilePreviewResponse{
		Profile: profile,
	})
}
