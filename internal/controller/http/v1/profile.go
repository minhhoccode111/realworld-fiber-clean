package v1

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http/httpmeta"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http/v1/response"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/entity"
)

// @Summary     Get profile
// @Description Get profile by username
// @ID          profiles-get-by-username
// @Tags  	    profiles
// @Produce     json
// @Param       username path string true "Username of the profile to get"
// @Success     200 {object} response.ProfilePreviewResponse
// @Failure     400 {object} response.Error
// @Failure     401 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /profiles/{username} [get]
// @Security    BearerAuth
func (r *V1) getProfile(ctx *fiber.Ctx) error {
	isAuth, ok := ctx.Locals(httpmeta.CtxIsAuthKey).(bool)
	if !ok {
		isAuth = false
	}

	userID, ok := ctx.Locals(httpmeta.CtxUserIDKey).(string)
	if !ok || (userID == "" && isAuth) {
		return errorResponse(ctx, http.StatusUnauthorized, "cannot authorize user in jwt")
	}

	username := ctx.Params("username")
	if username == "" {
		return errorResponse(ctx, http.StatusBadRequest, "username is required")
	}

	p, err := r.p.Detail(ctx.UserContext(), userID, username)
	if err != nil {
		if errors.Is(err, entity.ErrNoRows) {
			return errorResponse(ctx, http.StatusNotFound, "Profile not found")
		}

		r.l.Error(err, "http - v1 - getProfile - r.p.Detail")

		return errorResponse(ctx, http.StatusInternalServerError, "database problems")
	}

	return ctx.Status(http.StatusOK).JSON(response.ProfilePreviewResponse{
		Profile: p,
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
	userID, ok := ctx.Locals(httpmeta.CtxUserIDKey).(string)
	if !ok || userID == "" {
		return errorResponse(ctx, http.StatusUnauthorized, "cannot authorize user in jwt")
	}

	username := ctx.Params("username")
	if username == "" {
		return errorResponse(ctx, http.StatusBadRequest, "username is required")
	}

	err := r.p.Follow(ctx.UserContext(), userID, username)
	if err != nil {
		if errors.Is(err, entity.ErrNoRows) {
			return errorResponse(ctx, http.StatusNotFound, "Profile not found")
		}

		if errors.Is(err, entity.ErrNoEffect) {
			return errorResponse(ctx, http.StatusNotFound, "Profile is already followed")
		}

		r.l.Error(err, "http - v1 - postFollowProfile - r.p.Follow")

		return errorResponse(ctx, http.StatusInternalServerError, "database problems")
	}

	p, err := r.p.Detail(ctx.UserContext(), userID, username)
	if err != nil {
		if errors.Is(err, entity.ErrNoRows) {
			return errorResponse(ctx, http.StatusNotFound, "Profile not found")
		}

		r.l.Error(err, "http - v1 - postFollowProfile - r.p.Detail")

		return errorResponse(ctx, http.StatusInternalServerError, "database problems")
	}

	return ctx.Status(http.StatusOK).JSON(response.ProfilePreviewResponse{
		Profile: p,
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
	userID, ok := ctx.Locals(httpmeta.CtxUserIDKey).(string)
	if !ok || userID == "" {
		return errorResponse(ctx, http.StatusUnauthorized, "cannot authorize user in jwt")
	}

	username := ctx.Params("username")
	if username == "" {
		return errorResponse(ctx, http.StatusBadRequest, "username is required")
	}

	err := r.p.Unfollow(ctx.UserContext(), userID, username)
	if err != nil {
		if errors.Is(err, entity.ErrNoRows) {
			return errorResponse(ctx, http.StatusNotFound, "Profile not found")
		}

		if errors.Is(err, entity.ErrNoEffect) {
			return errorResponse(ctx, http.StatusNotFound, "Profile is already unfollowed")
		}

		r.l.Error(err, "http - v1 - postFollowProfile - r.p.Follow")

		return errorResponse(ctx, http.StatusInternalServerError, "database problems")
	}

	p, err := r.p.Detail(ctx.UserContext(), userID, username)
	if err != nil {
		if errors.Is(err, entity.ErrNoRows) {
			return errorResponse(ctx, http.StatusNotFound, "Profile not found")
		}

		r.l.Error(err, "http - v1 - postFollowProfile - r.p.Detail")

		return errorResponse(ctx, http.StatusInternalServerError, "database problems")
	}

	return ctx.Status(http.StatusOK).JSON(response.ProfilePreviewResponse{
		Profile: p,
	})
}
