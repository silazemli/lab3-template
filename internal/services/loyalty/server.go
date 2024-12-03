package loyalty

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type server struct {
	srv echo.Echo
	db  loyaltyStorage
}

func NewServer(db loyaltyStorage) server {
	srv := server{}
	srv.db = db
	srv.srv = *echo.New()
	api := srv.srv.Group("/api/loyalty")

	api.GET("/me", srv.GetUser)                   // +
	api.PATCH("/increment", srv.IncrementCounter) // +
	api.PATCH("/decrement", srv.DecrementCounter) // +

	srv.srv.GET("/manage/health", srv.HealthCheck)

	return srv
}

func (srv *server) Start() error {
	err := srv.srv.Start(":8050")
	if err != nil {
		return err
	}
	return nil
}

func (srv *server) GetUser(ctx echo.Context) error {
	username := ctx.Request().Header.Get("X-User-Name")
	user, err := srv.db.GetUser(username)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ctx.JSON(http.StatusNotFound, echo.Map{})
	}
	if err != nil {
		return ctx.JSON(http.StatusNotFound, echo.Map{"error": err})
	}
	return ctx.JSON(http.StatusOK, user)
}

func (srv *server) IncrementCounter(ctx echo.Context) error {
	username := ctx.Request().Header.Get("X-User-Name")
	err := srv.db.IncrementCounter(username)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{})
	}
	return nil
}

func (srv *server) DecrementCounter(ctx echo.Context) error {
	username := ctx.Request().Header.Get("X-User-Name")
	err := srv.db.DecrementCounter(username)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{})
	}
	return nil
}

func (srv *server) HealthCheck(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, echo.Map{})
}
