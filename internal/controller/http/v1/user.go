package v1

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/minhhoccode111/realworld-fiber-clean/internal/controller/http/common"
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
func (r *V1) postRegisterUser(ctx *fiber.Ctx) error {
	var body request.UserRegisterRequest

	if err := ctx.BodyParser(&body); err != nil {
		r.l.Error(err, "http - v1 - postRegisterUser - ctx.BodyParser")

		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	body.User.Trim()

	if err := r.v.Struct(body.User); err != nil {
		r.l.Error(err, "http - v1 - postRegisterUser - r.v.Struct")

		errs := validatorx.ExtractErrors(err)
		return errorResponse(ctx, http.StatusBadRequest, strings.Join(errs, "; "))
	}

	u := &entity.User{
		Username: body.User.Username,
		Email:    body.User.Email,
		Password: body.User.Password,
	}
	// user.ID generated, user.Password hashed
	err := r.u.Register(ctx.UserContext(), u)
	if err != nil {
		if errors.Is(err, entity.ErrConflict) {
			return errorResponse(ctx, http.StatusConflict, "email/username already existed")
		}

		r.l.Error(err, "http - v1 - postRegisterUser - r.u.Register")

		return errorResponse(ctx, http.StatusInternalServerError, "database problems")
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

		return errorResponse(ctx, http.StatusInternalServerError, "generate jwt error")
	}

	ctx.Cookie(NewJWTCookie(token, r.cfg.JWT.Expiration))

	return ctx.Status(http.StatusCreated).JSON(response.UserAuthResponse{
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
func (r *V1) postLoginUser(ctx *fiber.Ctx) error {
	var body request.UserLoginRequest

	if err := ctx.BodyParser(&body); err != nil {
		r.l.Error(err, "http - v1 - postLoginUser - ctx.BodyParser")

		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	body.User.Trim()

	if err := r.v.Struct(body.User); err != nil {
		r.l.Error(err, "http - v1 - postLoginUser - r.v.Struct")

		errs := validatorx.ExtractErrors(err)
		return errorResponse(ctx, http.StatusBadRequest, strings.Join(errs, "; "))
	}

	u, err := r.u.Login(ctx.UserContext(), &entity.User{
		Email:    body.User.Email,
		Password: body.User.Password,
	})
	if err != nil {
		if errors.Is(err, entity.ErrNoRows) {
			return errorResponse(ctx, http.StatusUnauthorized, "incorrect email")
		}

		if errors.Is(err, entity.ErrInvalidCredentials) {
			return errorResponse(ctx, http.StatusUnauthorized, "incorrect password")
		}

		r.l.Error(err, "http - v1 - postLoginUser - r.u.Login")

		return errorResponse(ctx, http.StatusInternalServerError, "database problems")
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

		return errorResponse(ctx, http.StatusInternalServerError, "jwt problems")
	}

	ctx.Cookie(NewJWTCookie(token, r.cfg.JWT.Expiration))

	return ctx.Status(http.StatusOK).JSON(response.UserAuthResponse{
		User: response.NewUserAuth(u, token),
	})
}

// @Summary     Logout User
// @Description Logout User by clearing JWT cookie
// @ID          users-logout
// @Tags  	    users
// @Success     204
// @Router      /users/logout [post]
func (r *V1) postLogoutUser(ctx *fiber.Ctx) error {
	ctx.Cookie(NewJWTCookie("", -time.Hour))
	return ctx.SendStatus(http.StatusNoContent)
}

// @Summary     Get current User
// @Description Get current User
// @ID          users-current
// @Tags  	    users
// @Produce     json
// @Success     200 {object} response.UserAuthResponse
// @Failure     500 {object} response.Error
// @Router      /user [get]
// @Security    BearerAuth
func (r *V1) getCurrentUser(ctx *fiber.Ctx) error {
	userID := ctx.Locals(common.CtxUserIDKey).(string)
	if userID == "" {
		return errorResponse(ctx, http.StatusInternalServerError, "cannot authorize user in jwt")
	}

	u, err := r.u.Current(ctx.UserContext(), userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorResponse(ctx, http.StatusNotFound, "userID in token not found")
		}

		r.l.Error(err, "http - v1 - getCurrentUser - r.u.Current")

		return errorResponse(ctx, http.StatusInternalServerError, "database problems")
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

		return errorResponse(ctx, http.StatusInternalServerError, "jwt problems")
	}

	ctx.Cookie(NewJWTCookie(token, r.cfg.JWT.Expiration))

	return ctx.Status(http.StatusOK).JSON(response.UserAuthResponse{
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
func (r *V1) putUpdateUser(ctx *fiber.Ctx) error {
	var body request.UserUpdateRequest

	if err := ctx.BodyParser(&body); err != nil {
		r.l.Error(err, "http - v1 - putUpdateUser - ctx.BodyParser")

		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	body.User.Trim()

	if body.User.IsAllEmpty() {
		return errorResponse(ctx, http.StatusBadRequest, "no field provided")
	}

	if err := r.v.Struct(body.User); err != nil {
		r.l.Error(err, "http - v1 - putUpdateUser - r.v.Struct")

		errs := validatorx.ExtractErrors(err)
		return errorResponse(ctx, http.StatusBadRequest, strings.Join(errs, "; "))
	}

	userID := ctx.Locals(common.CtxUserIDKey).(string)
	if userID == "" {
		return errorResponse(ctx, http.StatusInternalServerError, "cannot authorize user in jwt")
	}

	u, err := r.u.Update(ctx.UserContext(), body.User.NewUser(userID))
	if err != nil {
		if errors.Is(err, entity.ErrConflict) {
			return errorResponse(ctx, http.StatusConflict, "email/username alread existed")
		}

		r.l.Error(err, "http - v1 - putUpdateUser - r.u.Update")

		return errorResponse(ctx, http.StatusInternalServerError, "database problems")
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

		return errorResponse(ctx, http.StatusInternalServerError, "generate jwt error")
	}

	ctx.Cookie(NewJWTCookie(token, r.cfg.JWT.Expiration))

	return ctx.Status(http.StatusOK).JSON(response.UserAuthResponse{
		User: response.NewUserAuth(u, token),
	})
}
