package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
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
// @Success     200 {object} response.ProfilePreviewResponse
// @Failure     400 {object} response.Error
// @Failure     401 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /profiles/{username} [get]
// @Security    BearerAuth
func (r *V1) getProfile(c *gin.Context) {
	isAuth := c.MustGet(string(middleware.CtxIsAuthKey)).(bool)

	userID := ""
	if isAuth {
		userID = c.MustGet(string(middleware.CtxUserIDKey)).(string)
		if userID == "" {
			errorResponse(c, http.StatusUnauthorized, "cannot authorize user in jwt")
			return
		}
	}

	username := c.Param("username")
	if username == "" {
		errorResponse(c, http.StatusBadRequest, "username is required")
		return
	}

	p, err := r.p.Detail(c.Request.Context(), userID, username)
	if err != nil {
		if errors.Is(err, entity.ErrNoRows) {
			errorResponse(c, http.StatusNotFound, "Profile not found")
			return
		}

		r.l.Error(err, "http - v1 - getProfile - r.p.Detail")

		errorResponse(c, http.StatusInternalServerError, "database problems")
		return
	}

	c.JSON(http.StatusOK, response.ProfilePreviewResponse{
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
func (r *V1) postFollowProfile(c *gin.Context) {
	userID := c.MustGet(string(middleware.CtxUserIDKey)).(string)
	if userID == "" {
		errorResponse(c, http.StatusUnauthorized, "cannot authorize user in jwt")
		return
	}

	username := c.Param("username")
	if username == "" {
		errorResponse(c, http.StatusBadRequest, "username is required")
		return
	}

	err := r.p.Follow(c.Request.Context(), userID, username)
	if err != nil {
		if errors.Is(err, entity.ErrNoRows) {
			errorResponse(c, http.StatusNotFound, "Profile not found")
			return
		}

		if errors.Is(err, entity.ErrNoEffect) {
			errorResponse(c, http.StatusBadRequest, "Profile is already followed")
			return
		}

		r.l.Error(err, "http - v1 - postFollowProfile - r.p.Follow")

		errorResponse(c, http.StatusInternalServerError, "database problems")
		return
	}

	p, err := r.p.Detail(c.Request.Context(), userID, username)
	if err != nil {
		if errors.Is(err, entity.ErrNoRows) {
			errorResponse(c, http.StatusNotFound, "Profile not found")
			return
		}

		r.l.Error(err, "http - v1 - postFollowProfile - r.p.Detail")

		errorResponse(c, http.StatusInternalServerError, "database problems")
		return
	}

	c.JSON(http.StatusOK, response.ProfilePreviewResponse{
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
func (r *V1) deleteFollowProfile(c *gin.Context) {
	userID := c.MustGet(string(middleware.CtxUserIDKey)).(string)
	if userID == "" {
		errorResponse(c, http.StatusUnauthorized, "cannot authorize user in jwt")
		return
	}

	username := c.Param("username")
	if username == "" {
		errorResponse(c, http.StatusBadRequest, "username is required")
		return
	}

	err := r.p.Unfollow(c.Request.Context(), userID, username)
	if err != nil {
		if errors.Is(err, entity.ErrNoRows) {
			errorResponse(c, http.StatusNotFound, "Profile not found")
			return
		}

		if errors.Is(err, entity.ErrNoEffect) {
			errorResponse(c, http.StatusBadRequest, "Profile is already unfollowed")
			return
		}

		r.l.Error(err, "http - v1 - postFollowProfile - r.p.Follow")

		errorResponse(c, http.StatusInternalServerError, "database problems")
		return
	}

	p, err := r.p.Detail(c.Request.Context(), userID, username)
	if err != nil {
		if errors.Is(err, entity.ErrNoRows) {
			errorResponse(c, http.StatusNotFound, "Profile not found")
			return
		}

		r.l.Error(err, "http - v1 - postFollowProfile - r.p.Detail")

		errorResponse(c, http.StatusInternalServerError, "database problems")
		return
	}

	c.JSON(http.StatusOK, response.ProfilePreviewResponse{
		Profile: p,
	})
}
