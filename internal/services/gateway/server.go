package gateway

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/RohanPoojary/gomq"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	log "github.com/rs/zerolog/log"
	circuit "github.com/rubyist/circuitbreaker"
	"github.com/silazemli/lab3-template/internal/services/gateway/async"
	"github.com/silazemli/lab3-template/internal/services/gateway/clients"
	"github.com/silazemli/lab3-template/internal/services/loyalty"
	"github.com/silazemli/lab3-template/internal/services/payment"
	"github.com/silazemli/lab3-template/internal/services/reservation"
)

type Server struct {
	srv         echo.Echo
	cfg         Config
	broker      gomq.Broker
	reservation clients.ReservationClient
	payment     clients.PaymentClient
	loyalty     clients.LoyaltyClient
}

func NewServer() Server {
	srv := Server{}
	srv.srv = *echo.New()
	srv.cfg = *NewConfig()

	srv.loyalty = *clients.NewLoyaltyClient(circuit.NewHTTPClient(0, 10, nil), srv.cfg.LoyaltyService)
	srv.payment = *clients.NewPaymentClient(circuit.NewHTTPClient(0, 10, nil), srv.cfg.PaymentService)
	srv.reservation = *clients.NewReservationClient(circuit.NewHTTPClient(0, 10, nil), srv.cfg.ReservationService)

	srv.broker = gomq.NewAsyncBroker()
	async.LoyaltyDecrementRetry(srv.broker, &srv.loyalty)

	api := srv.srv.Group("/api/v1")
	api.GET("/hotels", srv.GetAllHotels)
	api.GET("/me", srv.GetUser)
	api.GET("/loyalty", srv.GetStatus)
	api.GET("/reservations", srv.GetAllReservations)
	api.GET("/reservations/:reservationUid", srv.GetReservation)
	api.POST("/reservations", srv.MakeReservation)
	api.DELETE("/reservations/:reservationUid", srv.CancelReservation)

	srv.srv.GET("/manage/health", srv.HealthCheck)

	return srv
}

func (srv *Server) Start() error {
	err := srv.srv.Start(":8080")
	if err != nil {
		return err
	}
	return nil
}

func (srv *Server) GetUser(ctx echo.Context) error {
	username := ctx.Request().Header.Get("X-User-Name")
	response := userInfoResponse{}

	reservations, err := srv.reservation.GetReservations(username) // create a list of reservations
	if err != nil {
		reservations = make([]reservation.Reservation, 1)
		reservations[0] = reservation.Reservation{}
	}
	reservationsResponse := make([]reservationResponse, len(reservations))
	for index, theReservation := range reservations {
		reservationsResponse[index] = srv.createReservationResponse(theReservation)
	}
	response.Reservations = reservationsResponse

	theLoyalty, err := srv.loyalty.GetUser(username) // create this specific loyalty response
	if err != nil {
		theLoyalty = loyalty.Loyalty{}
	}

	loyaltyResponse := createLoyaltyResponseNoCount(theLoyalty) // out of ideas for names
	if err != nil {
		loyaltyResponse.Discount = ""
		ctx.Response().Write([]byte("Loyalty Service Unavailable"))
	}
	response.Loyalty = loyaltyResponse

	return ctx.JSON(http.StatusOK, response)
}

func (srv *Server) GetAllHotels(ctx echo.Context) error {
	pageStr := ctx.QueryParam("page")
	sizeStr := ctx.QueryParam("size")

	page := 1
	size := 1

	if pageStr != "" {
		var err error
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 0 {
			log.Info().Msg(err.Error())
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid page number"})
		}
	}

	if sizeStr != "" {
		var err error
		size, err := strconv.Atoi(sizeStr)
		if err != nil || size < 1 || size > 100 {
			log.Info().Msg(err.Error())
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid size number"})
		}
	}

	hotels, err := srv.reservation.GetAllHotels()
	if err != nil {
		log.Info().Msg(err.Error())
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err})
	}

	start := (page - 1) * size
	end := start + size
	if end > len(hotels) {
		end = len(hotels)
	}
	if end-start < size {
		size = end - start
	}

	response := struct {
		Page   int                 `json:"page"`
		Size   int                 `json:"pageSize"`
		Total  int                 `json:"totalElements"`
		Hotels []reservation.Hotel `json:"items"`
	}{
		Page:   page,
		Size:   size,
		Total:  len(hotels),
		Hotels: hotels[start:end],
	}

	return ctx.JSON(http.StatusOK, response)
}

func (srv *Server) GetAllReservations(ctx echo.Context) error {
	username := ctx.Request().Header.Get("X-User-Name")
	reservations, err := srv.reservation.GetReservations(username)
	if err != nil {
		log.Info().Msg(err.Error())
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err})
	}
	response := make([]reservationResponse, len(reservations))
	for index, theReservation := range reservations {
		response[index] = srv.createReservationResponse(theReservation)
	}
	return ctx.JSON(http.StatusOK, response)
}

func (srv *Server) GetReservation(ctx echo.Context) error {
	username := ctx.Request().Header.Get("X-User-Name")
	reservationUID := ctx.Param("reservationUid")
	theReservation, err := srv.reservation.GetReservation(reservationUID)
	if err != nil {
		log.Info().Msg(err.Error())
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err})
	}
	if username != theReservation.Username {
		return ctx.JSON(http.StatusForbidden, echo.Map{"error": err})
	}
	response := srv.createReservationResponse(theReservation)
	return ctx.JSON(http.StatusOK, response)
}

func (srv *Server) GetStatus(ctx echo.Context) error {
	username := ctx.Request().Header.Get("X-User-Name")
	loyalty, err := srv.loyalty.GetUser(username)
	if err != nil {
		log.Info().Msg(err.Error())
		return ctx.JSON(http.StatusServiceUnavailable, echo.Map{"error": err})
	}
	response := createLoyaltyResponse(loyalty)

	return ctx.JSON(http.StatusOK, response)
}

func (srv *Server) MakeReservation(ctx echo.Context) error {
	body, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err) // parse request
	}
	var reservationRequest struct {
		HotelUID  string `json:"hotelUid"`
		StartDate string `json:"startDate"`
		EndDate   string `json:"endDate"`
	}
	if err := json.Unmarshal(body, &reservationRequest); err != nil {
		log.Info().Msg(err.Error())
		return fmt.Errorf("failed to unmarshal request body: %w", err)
	}

	hotelUID := reservationRequest.HotelUID
	hotelID, err := srv.reservation.GetHotelID(hotelUID) // getting hotel ID and hotel by ID for some reason
	if err != nil {
		log.Info().Msg(err.Error())
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err})
	}
	hotel, err := srv.reservation.GetHotel(strconv.Itoa(hotelID))
	if err != nil {
		log.Info().Msg(err.Error())
		return fmt.Errorf("hotel not found: %w", err)
	}

	dateLayout := "2006-01-02"
	startDate, err := time.Parse(dateLayout, reservationRequest.StartDate) // parsing time
	if err != nil {
		log.Info().Msg(err.Error())
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err})
	}
	endDate, err := time.Parse(dateLayout, reservationRequest.EndDate)
	if err != nil {
		log.Info().Msg(err.Error())
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err})
	}
	duration := int(endDate.Sub(startDate).Hours() / 24)
	if duration < 0 {
		log.Info().Msg(err.Error())
		return ctx.JSON(http.StatusBadRequest, echo.Map{})
	}

	username := ctx.Request().Header.Get("X-User-Name") // getting the discount
	user, err := srv.loyalty.GetUser(username)
	if err != nil {
		log.Info().Msg(err.Error())
		return ctx.JSON(http.StatusServiceUnavailable, echo.Map{"error": err})
	}
	discount := user.Discount

	price := duration * hotel.Price * (100 - discount) / 100 // calculating price
	thePayment := payment.Payment{
		PaymentUID: uuid.New().String(),
		Status:     "PAID",
		Price:      price,
	}
	err = srv.payment.CreatePayment(thePayment)
	if err != nil {
		log.Info().Msg(err.Error())
		return ctx.JSON(http.StatusServiceUnavailable, echo.Map{"error": err})
	}

	theReservation := reservation.Reservation{
		ReservationUID: uuid.New().String(),
		Username:       username,
		StartDate:      startDate.Format(dateLayout),
		EndDate:        endDate.Format(dateLayout),
		Status:         "PAID",
		HotelID:        hotelID,
		PaymentUID:     thePayment.PaymentUID,
	}

	err = srv.reservation.MakeReservation(theReservation)
	if err != nil {
		log.Info().Msg(err.Error())
		return ctx.JSON(http.StatusServiceUnavailable, echo.Map{"error": err})
	}

	err = srv.loyalty.IncrementCounter(username)
	if err != nil {
		log.Info().Msg(err.Error())
		srv.payment.CancelPayment(thePayment.PaymentUID)
		ctx.Response().Write([]byte("Loyalty Service Unavailable"))
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err})
	}

	return ctx.JSON(http.StatusOK, srv.createReservationCreatedResponse(theReservation))
}

func (srv *Server) CancelReservation(ctx echo.Context) error {
	reservationUID := ctx.Param("reservationUid")
	err := srv.reservation.CancelReservation(reservationUID)
	if err != nil {
		log.Info().Msg(err.Error())
		return ctx.JSON(http.StatusBadGateway, echo.Map{})
	}

	reservation, err := srv.reservation.GetReservation(reservationUID)
	if err != nil {
		log.Info().Msg(err.Error())
		return ctx.JSON(http.StatusBadGateway, echo.Map{})
	}

	err = srv.payment.CancelPayment(reservation.PaymentUID)
	if err != nil {
		log.Info().Msg(err.Error())
		return ctx.JSON(http.StatusBadGateway, echo.Map{})
	}

	username := ctx.Request().Header.Get("X-User-Name")
	err = srv.loyalty.DecrementCounter(username)
	if err != nil {
		log.Info().Msg(err.Error())
		return ctx.JSON(http.StatusBadGateway, echo.Map{})
	}
	return ctx.JSON(http.StatusNoContent, echo.Map{})
}

func (srv *Server) HealthCheck(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, echo.Map{})
}
