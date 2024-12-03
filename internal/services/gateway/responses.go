package gateway

import (
	"strconv"

	"github.com/silazemli/lab2-template/internal/services/loyalty"
	"github.com/silazemli/lab2-template/internal/services/payment"
	"github.com/silazemli/lab2-template/internal/services/reservation"
)

type paymentResponse struct {
	Status string `json:"status"`
	Price  int    `json:"price"`
}

type reservationResponse struct {
	ReservationUID string          `json:"reservationUid"`
	Hotel          hotelResponse   `json:"hotel"`
	StartDate      string          `json:"startDate"`
	EndDate        string          `json:"endDate"`
	Status         string          `json:"status"`
	Payment        paymentResponse `json:"payment"`
}

type loyaltyResponse struct {
	Status           string `json:"status"`
	Discount         string `json:"discount"`
	ReservationCount int    `json:"reservationCount"`
}

type loyaltyResponseNoCount struct {
	Status   string `json:"status"`
	Discount string `json:"discount"`
}

type userInfoResponse struct {
	Reservations []reservationResponse  `json:"reservations"`
	Loyalty      loyaltyResponseNoCount `json:"loyalty"`
}

type hotelResponse struct {
	HotelUID    string `json:"hotelUid"`
	Name        string `json:"name"`
	FullAddress string `json:"fullAddress"`
	Stars       int    `json:"stars"`
}

type reservationCreatedResponse struct {
	ReservationUID string          `json:"reservationUid"`
	HotelUID       string          `json:"hotelUid"`
	StartDate      string          `json:"startDate"`
	EndDate        string          `json:"endDate"`
	Discount       string          `json:"discount"`
	Status         string          `json:"status"`
	Payment        paymentResponse `json:"payment"`
}

func (srv *server) createReservationResponse(theReservation reservation.Reservation) reservationResponse {
	response := reservationResponse{}
	response.ReservationUID = theReservation.ReservationUID
	response.StartDate = ymd(theReservation.StartDate)
	response.EndDate = ymd(theReservation.EndDate)
	response.Status = theReservation.Status

	hotel, err := srv.reservation.GetHotel(strconv.Itoa(theReservation.HotelID))
	if err != nil {
		return reservationResponse{}
	}
	response.Hotel = createHotelResponse(hotel)

	payment, err := srv.payment.GetPayment(theReservation.PaymentUID)
	if err != nil {
		return reservationResponse{}
	}
	response.Payment = createPaymentResponse(payment)

	return response
}

func createHotelResponse(hotel reservation.Hotel) hotelResponse {
	return hotelResponse{
		HotelUID:    hotel.HotelUID,
		Name:        hotel.Name,
		FullAddress: hotel.Country + ", " + hotel.City + ", " + hotel.Address,
		Stars:       hotel.Stars,
	}
}

func createPaymentResponse(thePayment payment.Payment) paymentResponse {
	return paymentResponse{
		Status: thePayment.Status,
		Price:  thePayment.Price,
	}
}

func createLoyaltyResponse(theLoyalty loyalty.Loyalty) loyaltyResponse {
	return loyaltyResponse{
		Status:           theLoyalty.Status,
		Discount:         strconv.Itoa(theLoyalty.Discount),
		ReservationCount: theLoyalty.ReservationCount,
	}
}

func createLoyaltyResponseNoCount(theLoyalty loyalty.Loyalty) loyaltyResponseNoCount {
	return loyaltyResponseNoCount{
		Status:   theLoyalty.Status,
		Discount: strconv.Itoa(theLoyalty.Discount),
	}
}

func (srv *server) createReservationCreatedResponse(theReservation reservation.Reservation) reservationCreatedResponse {
	response := reservationCreatedResponse{}
	response.ReservationUID = theReservation.ReservationUID
	response.StartDate = ymd(theReservation.StartDate)
	response.EndDate = ymd(theReservation.EndDate)
	response.Status = theReservation.Status

	hotel, err := srv.reservation.GetHotel(strconv.Itoa(theReservation.HotelID))
	if err != nil {
		return reservationCreatedResponse{}
	}
	response.HotelUID = hotel.HotelUID

	payment, err := srv.payment.GetPayment(theReservation.PaymentUID)
	if err != nil {
		return reservationCreatedResponse{}
	}
	response.Payment = createPaymentResponse(payment)

	loyalty, err := srv.loyalty.GetUser(theReservation.Username)
	if err != nil {
		return reservationCreatedResponse{}
	}
	response.Discount = strconv.Itoa(loyalty.Discount)

	return response
}

func ymd(date string) string {
	return date[0:10]
}
