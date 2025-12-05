package v1

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http/middleware"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http/v1/request"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http/v1/response"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/entity"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/utils"
	"github.com/minhhoccode111/realworld-fiber-clean/pkg/validatorx"
)

// @Summary     Register User
// @Description Register User
// @ID          users-register
// @Tags  	    users
// @Accept      json
// @Produce     json
// @Param       request body request.UserRegisterRequest true "Register User"
// @Success     201 {object} response.UserAuthResponse
// @Failure     400 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /users [post]
func (r *V1) postRegisterUser(c *gin.Context) {
	var body request.UserRegisterRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		r.l.Error(err, "http - v1 - postRegisterUser - c.ShouldBindJSON")

		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	body.User.Trim()

	if err := r.v.Struct(body.User); err != nil {
		r.l.Error(err, "http - v1 - postRegisterUser - r.v.Struct")

		errs := validatorx.ExtractErrors(err)
		errorResponse(c, http.StatusBadRequest, strings.Join(errs, "; "))
		return
	}

	u := &entity.User{
		Username: body.User.Username,
		Email:    body.User.Email,
		Password: body.User.Password,
	}
	// user.ID generated, user.Password hashed
	err := r.u.Register(c.Request.Context(), u)
	if err != nil {
		if errors.Is(err, entity.ErrConflict) {
			errorResponse(c, http.StatusConflict, "email/username already existed")
			return
		}

		r.l.Error(err, "http - v1 - postRegisterUser - r.u.Register")

		errorResponse(c, http.StatusInternalServerError, "database problems")
		return
	}

	token, err := utils.GenerateJWT(
		u.ID,
		u.Role.String(),
		r.cfg.JWT.Secret,
		r.cfg.JWT.Issuer,
		r.cfg.JWT.Expiration,
	)
	if err != nil {
		r.l.Error(err, "http - v1 - postRegisterUser - utils.GenerateJWT")

		errorResponse(c, http.StatusInternalServerError, "generate jwt error")
		return
	}

	c.JSON(http.StatusCreated, response.UserAuthResponse{
		User: response.NewUserAuth(u, token),
	})
}

// @Summary     Login User
// @Description Login User
// @ID          users-login
// @Tags  	    users
// @Accept      json
// @Param       request body request.UserLoginRequest true "Login User"
// @Produce     json
// @Success     200 {object} response.UserAuthResponse
// @Failure     400 {object} response.Error
// @Failure     401 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /users/login [post]
func (r *V1) postLoginUser(c *gin.Context) {
	var body request.UserLoginRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		r.l.Error(err, "http - v1 - postLoginUser - c.ShouldBindJSON")

		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	body.User.Trim()

	if err := r.v.Struct(body.User); err != nil {
		r.l.Error(err, "http - v1 - postLoginUser - r.v.Struct")

		errs := validatorx.ExtractErrors(err)
		errorResponse(c, http.StatusBadRequest, strings.Join(errs, "; "))
		return
	}

	u, err := r.u.Login(c.Request.Context(), &entity.User{
		Email:    body.User.Email,
		Password: body.User.Password,
	})
	if err != nil {
		if errors.Is(err, entity.ErrNoRows) {
			errorResponse(c, http.StatusUnauthorized, "incorrect email")
			return
		}

		if errors.Is(err, entity.ErrInvalidCredentials) {
			errorResponse(c, http.StatusUnauthorized, "incorrect password")
			return
		}

		r.l.Error(err, "http - v1 - postLoginUser - r.u.Login")

		errorResponse(c, http.StatusInternalServerError, "database problems")
		return
	}

	token, err := utils.GenerateJWT(
		u.ID,
		u.Role.String(),
		r.cfg.JWT.Secret,
		r.cfg.JWT.Issuer,
		r.cfg.JWT.Expiration,
	)
	if err != nil {
		r.l.Error(err, "http - v1 - postLoginUser - utils.GenerateJWT")

		errorResponse(c, http.StatusInternalServerError, "jwt problems")
		return
	}

	c.JSON(http.StatusOK, response.UserAuthResponse{
		User: response.NewUserAuth(u, token),
	})
}

// func (r *V1) postLogoutUser(c *gin.Context) {
// 	c.Status(200) // Gin equivalent for ctx.SendStatus(200)
// }

// @Summary     Get current User
// @Description Get current User
// @ID          users-current
// @Tags  	    users
// @Produce     json
// @Success     200 {object} response.UserAuthResponse
// @Failure     500 {object} response.Error
// @Router      /user [get]
// @Security    BearerAuth
func (r *V1) getCurrentUser(c *gin.Context) {
	userID := c.MustGet(string(middleware.CtxUserIDKey)).(string)
	if userID == "" {
		errorResponse(c, http.StatusInternalServerError, "cannot authorize user in jwt")
		return
	}

	u, err := r.u.Current(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			errorResponse(c, http.StatusNotFound, "userID in token not found")
			return
		}

		r.l.Error(err, "http - v1 - getCurrentUser - r.u.Current")

		errorResponse(c, http.StatusInternalServerError, "database problems")
		return
	}

	token, err := utils.GenerateJWT(
		u.ID,
		u.Role.String(),
		r.cfg.JWT.Secret,
		r.cfg.JWT.Issuer,
		r.cfg.JWT.Expiration,
	)
	if err != nil {
		r.l.Error(err, "http - v1 - getCurrentUser - utils.GenerateJWT")

		errorResponse(c, http.StatusInternalServerError, "jwt problems")
		return
	}

	c.JSON(http.StatusOK, response.UserAuthResponse{
		User: response.NewUserAuth(u, token),
	})
}

// @Summary     Update User
// @Description Update User
// @ID          users-update
// @Tags  	    users
// @Accept      json
// @Produce     json
// @Param       request body request.UserUpdateRequest true "Update User"
// @Success     200 {object} response.UserAuthResponse
// @Failure     400 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /user [put]
// @Security    BearerAuth
func (r *V1) putUpdateUser(c *gin.Context) {
	var body request.UserUpdateRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		r.l.Error(err, "http - v1 - putUpdateUser - c.ShouldBindJSON")

		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	body.User.Trim()

	if body.User.IsAllEmpty() {
		errorResponse(c, http.StatusBadRequest, "no field provided")
		return
	}

	if err := r.v.Struct(body.User); err != nil {
		r.l.Error(err, "http - v1 - putUpdateUser - r.v.Struct")

		errs := validatorx.ExtractErrors(err)
		errorResponse(c, http.StatusBadRequest, strings.Join(errs, "; "))
		return
	}

	userID := c.MustGet(string(middleware.CtxUserIDKey)).(string)
	if userID == "" {
		errorResponse(c, http.StatusInternalServerError, "cannot authorize user in jwt")
		return
	}

	u, err := r.u.Update(c.Request.Context(), body.User.NewUser(userID))
	if err != nil {
		if errors.Is(err, entity.ErrConflict) {
			errorResponse(c, http.StatusConflict, "email/username alread existed")
			return
		}

		r.l.Error(err, "http - v1 - putUpdateUser - r.u.Update")

		errorResponse(c, http.StatusInternalServerError, "database problems")
		return
	}

	token, err := utils.GenerateJWT(
		u.ID,
		u.Role.String(),
		r.cfg.JWT.Secret,
		r.cfg.JWT.Issuer,
		r.cfg.JWT.Expiration,
	)
	if err != nil {
		r.l.Error(err, "http - v1 - putUpdateUser - utils.GenerateJWT")

		errorResponse(c, http.StatusInternalServerError, "generate jwt error")
		return
	}

	c.JSON(http.StatusOK, response.UserAuthResponse{
		User: response.NewUserAuth(u, token),
	})
}
