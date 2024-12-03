package payment

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type server struct {
	srv echo.Echo
	db  paymentStorage
}

func NewServer(db paymentStorage) server {
	srv := server{}
	srv.db = db
	srv.srv = *echo.New()
	api := srv.srv.Group("/api/payment")
	api.POST("", srv.PostPayment)         // +
	api.PATCH("/:uid", srv.CancelPayment) // +
	api.GET("/:uid", srv.GetPayment)

	srv.srv.GET("/manage/health", srv.HealthCheck)

	return srv
}

func (srv *server) Start() error {
	err := srv.srv.Start(":8060")
	if err != nil {
		return err
	}
	return nil
}

func (srv *server) PostPayment(ctx echo.Context) error {
	var thePayment Payment
	err := ctx.Bind(&thePayment)
	if err != nil {
		return err
	}
	err = srv.db.PostPayment(thePayment)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusCreated, echo.Map{})
}

func (srv *server) CancelPayment(ctx echo.Context) error {
	UID := ctx.Param("uid")
	err := srv.db.CancelPayment(UID)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, echo.Map{})
}

func (srv *server) GetPayment(ctx echo.Context) error {
	UID := ctx.Param("uid")
	payment, err := srv.db.GetPayment(UID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, Payment{})
	}
	return ctx.JSON(http.StatusOK, payment)
}

func (srv *server) HealthCheck(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, echo.Map{})
}
