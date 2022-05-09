package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/agandreev/crime-app-auth/internal/domain"
	"github.com/agandreev/crime-app-auth/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

var (
	ErrTimout = errors.New("the request execution timeout has expired")
)

type UserService interface {
	SignIn(ctx context.Context, user domain.User) (string, error)
	SignUp(ctx context.Context, user domain.User) error
}

type UserHandler struct {
	userService UserService

	log *logrus.Logger
}

func NewUserHandler(service UserService, log *logrus.Logger) *UserHandler {
	return &UserHandler{
		userService: service,
		log:         log,
	}
}

func (h UserHandler) InitRoutes(e *echo.Echo, timeout time.Duration) {
	users := e.Group(
		"users",
		middleware.RequestID(),
		middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: "method=${method}, uri=${uri}, status=${status}\n",
		}),
		middleware.Recover(),
		middleware.TimeoutWithConfig(middleware.TimeoutConfig{
			Skipper:      middleware.DefaultSkipper,
			ErrorMessage: ErrTimout.Error(),
			OnTimeoutRouteErrorHandler: func(err error, c echo.Context) {
				c.Error(err)
			},
			Timeout: timeout * time.Second,
		}),
	)
	users.POST("/register", h.SignUpHandler)
	users.POST("/login", h.SignInHandler)
}

// SignUpHandler godoc
// @Summary      provides signing up operation
// @Description  registers user in crime-app microservices ecosystem
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   body      domain.User  true  "User's consisted of login and password"
// @Success      201  {object}  domain.User
// @Failure 400 {object} echo.HTTPError
// @Failure 401 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Failure default {object} echo.HTTPError
// @Router       /users/register [post]
func (h UserHandler) SignUpHandler(c echo.Context) error {
	ctx := c.Request().Context()

	user := domain.User{}
	if err := c.Bind(&user); err != nil {
		h.log.Error(err)
		return echo.NewHTTPError(utils.StatusCode(err), err.Error())
	}

	if err := h.userService.SignUp(ctx, user); err != nil {
		h.log.Error(err)
		return echo.NewHTTPError(utils.StatusCode(err), err.Error())
	}

	user.Password = nil
	return c.JSON(http.StatusCreated, user)
}

// SignInHandler godoc
// @Summary      provides signing in operation
// @Description  authorize user in crime-app microservices ecosystem
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   body      domain.User  true  "User's consisted of login and password"
// @Success      200  {object}  string
// @Failure 400 {object} echo.HTTPError
// @Failure 401 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Failure default {object} echo.HTTPError
// @Router       /users/login [post]
func (h UserHandler) SignInHandler(c echo.Context) error {
	ctx := c.Request().Context()

	user := domain.User{}
	if err := c.Bind(&user); err != nil {
		h.log.Error(err)
		return echo.NewHTTPError(utils.StatusCode(err), err.Error())
	}

	token, err := h.userService.SignIn(ctx, user)
	if err != nil {
		h.log.Error(err)
		return echo.NewHTTPError(utils.StatusCode(err), err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}
