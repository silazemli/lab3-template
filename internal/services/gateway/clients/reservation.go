package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/silazemli/lab2-template/internal/services/reservation"
)

type ReservationClient struct {
	client  HTTPClient
	baseURL string
}

func NewReservationClient(client HTTPClient, baseURL string) *ReservationClient {
	return &ReservationClient{
		client:  client,
		baseURL: baseURL,
	}
}

func (reservationClient *ReservationClient) GetAllHotels() ([]reservation.Hotel, error) {
	URL := fmt.Sprintf("%s/%s", reservationClient.baseURL, "hotels")
	request, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return []reservation.Hotel{}, fmt.Errorf("failed to build request: %s", err)
	}
	response, err := reservationClient.client.Do(request)
	if err != nil {
		return []reservation.Hotel{}, fmt.Errorf("failed to make request: %s", err)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return []reservation.Hotel{}, fmt.Errorf("failed to read response body: %s", err)
	}
	defer response.Body.Close()
	switch response.StatusCode {
	case http.StatusOK:
		var hotels []reservation.Hotel
		if err := json.Unmarshal(body, &hotels); err != nil {
			return []reservation.Hotel{}, fmt.Errorf("failed to unmarshal response body: %s", err)
		}
		return hotels, nil
	case http.StatusInternalServerError, http.StatusNotFound, http.StatusBadRequest:
		return []reservation.Hotel{}, fmt.Errorf("server error: %s", err)
	default:
		return []reservation.Hotel{}, fmt.Errorf("unknown error: %s", err)
	}
}

func (reservationClient *ReservationClient) GetReservations(username string) ([]reservation.Reservation, error) {
	URL := fmt.Sprintf("%s/%s", reservationClient.baseURL, "reservations")
	request, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return []reservation.Reservation{}, fmt.Errorf("failed to build request: %w", err)
	}
	request.Header.Set("X-User-Name", username)
	response, err := reservationClient.client.Do(request)
	if err != nil {
		return []reservation.Reservation{}, fmt.Errorf("failed to make request: %w", err)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return []reservation.Reservation{}, fmt.Errorf("failed to read response body: %w", err)
	}
	defer response.Body.Close()
	switch response.StatusCode {
	case http.StatusOK:
		var reservations []reservation.Reservation
		if err := json.Unmarshal(body, &reservations); err != nil {
			return []reservation.Reservation{}, fmt.Errorf("failed to unmarshal response body: %w", err)
		}
		return reservations, nil
	case http.StatusInternalServerError, http.StatusNotFound, http.StatusBadRequest:
		return []reservation.Reservation{}, fmt.Errorf("server error: %w", err)
	default:
		return []reservation.Reservation{}, fmt.Errorf("unknown error: %w", err)
	}
}

func (reservationClient *ReservationClient) GetReservation(reservationUID string) (reservation.Reservation, error) {
	URL := fmt.Sprintf("%s/%s/%s", reservationClient.baseURL, "reservations", reservationUID)
	fmt.Println("reservation")
	fmt.Println(URL)
	request, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return reservation.Reservation{}, fmt.Errorf("failed to build request: %w", err)
	}
	response, err := reservationClient.client.Do(request)
	if err != nil {
		return reservation.Reservation{}, fmt.Errorf("failed to make request: %w", err)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return reservation.Reservation{}, fmt.Errorf("failed to read response body: %w", err)
	}
	defer response.Body.Close()
	fmt.Println(response.StatusCode)
	switch response.StatusCode {
	case http.StatusOK:
		var theReservation reservation.Reservation
		if err := json.Unmarshal(body, &theReservation); err != nil {
			return reservation.Reservation{}, fmt.Errorf("failed to unmarshal response body: %w", err)
		}
		return theReservation, nil
	case http.StatusInternalServerError, http.StatusNotFound, http.StatusBadRequest:
		return reservation.Reservation{}, fmt.Errorf("server error: %w", err)
	default:
		fmt.Println(err)
		return reservation.Reservation{}, fmt.Errorf("unknown error: %w", err)
	}
}

func (reservationClient *ReservationClient) MakeReservation(theReservation reservation.Reservation) error {
	URL := fmt.Sprintf("%s/%s", reservationClient.baseURL, "reservations")
	body, err := json.Marshal(theReservation)
	if err != nil {
		return fmt.Errorf("failed to build request body: %w", err)
	}
	fmt.Println(body)
	request, err := http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}
	request.Header.Set("X-User_Name", theReservation.Username)
	request.Header.Set("Content-Type", "application/json")
	response, err := reservationClient.client.Do(request)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	switch response.StatusCode {
	case http.StatusCreated:
		return nil
	case http.StatusInternalServerError, http.StatusNotFound, http.StatusBadRequest:
		return fmt.Errorf("server error: %w", err)
	default:
		return fmt.Errorf("unknown error: %w", err)
	}
}

func (reservationClient *ReservationClient) CancelReservation(reservationUID string) error {
	URL := fmt.Sprintf("%s/%s/%s", reservationClient.baseURL, "reservations", reservationUID)
	request, err := http.NewRequest(http.MethodPatch, URL, nil)
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}
	response, err := reservationClient.client.Do(request)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	switch response.StatusCode {
	case http.StatusAccepted:
		return nil
	case http.StatusInternalServerError, http.StatusNotFound, http.StatusBadRequest:
		return fmt.Errorf("server error: %w", err)
	default:
		return fmt.Errorf("unknown error: %w", err)
	}
}

func (reservationClient *ReservationClient) GetHotelID(hotelUID string) (int, error) {
	URL := fmt.Sprintf("%s/%s/%s", reservationClient.baseURL, "hotels", hotelUID)
	request, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return -1, fmt.Errorf("failed to build request: %w", err)
	}
	response, err := reservationClient.client.Do(request)
	if err != nil {
		return -1, fmt.Errorf("failed to make request: %w", err)
	}
	defer response.Body.Close()
	switch response.StatusCode {
	case http.StatusOK:
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return -1, fmt.Errorf("failed to read response body: %w", err)
		}
		var hotelIDResponse struct {
			ID int `json:"id"`
		}
		if err := json.Unmarshal(body, &hotelIDResponse); err != nil {
			return -1, fmt.Errorf("failed to unmarshal response: %w", err)
		}
		return hotelIDResponse.ID, nil

	case http.StatusInternalServerError, http.StatusNotFound, http.StatusBadRequest:
		return -1, fmt.Errorf("server error: %w", err)
	default:
		return -1, fmt.Errorf("unknown error: %w", err)
	}
}

func (reservationClient ReservationClient) GetHotel(ID string) (reservation.Hotel, error) {
	URL := fmt.Sprintf("%s/%s/%s/%s", reservationClient.baseURL, "hotels", "hotel", ID)
	request, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return reservation.Hotel{}, fmt.Errorf("failed to build request: %w", err)
	}
	response, err := reservationClient.client.Do(request)
	if err != nil {
		return reservation.Hotel{}, fmt.Errorf("failed to make request: %w", err)
	}
	defer response.Body.Close()
	switch response.StatusCode {
	case http.StatusOK:
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return reservation.Hotel{}, fmt.Errorf("failed to read response body: %w", err)
		}
		var hotel reservation.Hotel
		if err := json.Unmarshal(body, &hotel); err != nil {
			return reservation.Hotel{}, fmt.Errorf("failed to unmarshal response: %w", err)
		}
		return hotel, nil

	case http.StatusInternalServerError, http.StatusNotFound, http.StatusBadRequest:
		return reservation.Hotel{}, fmt.Errorf("server error: %w", err)
	default:
		return reservation.Hotel{}, fmt.Errorf("unknown error: %w", err)
	}
}
