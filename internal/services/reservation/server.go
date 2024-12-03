package reservation

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type server struct {
	srv echo.Echo
	rdb reservationStorage
	hdb hotelStorage
}

func NewServer(hdb hotelStorage, rdb reservationStorage) server {
	srv := server{}
	srv.rdb = rdb
	srv.hdb = hdb
	srv.srv = *echo.New()
	api := srv.srv.Group("api/reservation")
	api.GET("/hotels", srv.GetAllHotels)                              // +
	api.GET("/reservations", srv.GetAllReservations)                  // +
	api.GET("/reservations/:reservationUID", srv.GetReservation)      // +
	api.POST("/reservations", srv.MakeReservation)                    // +
	api.PATCH("/reservations/:reservationUID", srv.CancelReservation) // +
	api.GET("/hotels/:hotelUID", srv.GetHotelID)
	api.GET("/hotels/hotel/:ID", srv.GetHotel)

	srv.srv.GET("/manage/health", srv.HealthCheck)

	return srv
}

func (srv *server) Start() error {
	err := srv.srv.Start(":8070")
	if err != nil {
		return err
	}
	return nil
}

func (srv *server) GetAllReservations(ctx echo.Context) error {
	username := ctx.Request().Header.Get("X-User-Name")
	reservations, err := srv.rdb.GetReservations(username)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err})
	}
	return ctx.JSON(http.StatusOK, reservations)
}

func (srv *server) GetAllHotels(ctx echo.Context) error {
	hotels, err := srv.hdb.GetAll()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err})
	}
	return ctx.JSON(http.StatusOK, hotels)
}

func (srv *server) GetReservation(ctx echo.Context) error {
	reservation, err := srv.rdb.GetReservation(ctx.Param("reservationUID"))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ctx.JSON(http.StatusNotFound, echo.Map{})
	}
	if err != nil {
		return ctx.JSON(http.StatusNotFound, echo.Map{"error": err})
	}

	return ctx.JSON(http.StatusOK, reservation)
}

func (srv *server) MakeReservation(ctx echo.Context) error {
	reservation := Reservation{}
	err := ctx.Bind(&reservation)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err})
	}
	err = srv.rdb.MakeReservation(reservation)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err})
	}
	return ctx.JSON(http.StatusCreated, echo.Map{})
}

func (srv *server) CancelReservation(ctx echo.Context) error {
	reservationUID := ctx.Param("reservationUID")
	err := srv.rdb.CancelReservation(reservationUID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err})
	}
	return ctx.JSON(http.StatusAccepted, echo.Map{})
}

func (srv *server) GetHotelID(ctx echo.Context) error {
	hotelUID := ctx.Param("hotelUID")
	ID, err := srv.hdb.GetHotelID(hotelUID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"id": ID})
}

func (srv *server) GetHotel(ctx echo.Context) error {
	hotelUID := ctx.Param("ID")
	hotel, err := srv.hdb.GetHotel(hotelUID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err})
	}
	return ctx.JSON(http.StatusOK, hotel)
}

func (srv *server) HealthCheck(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, echo.Map{})
}
