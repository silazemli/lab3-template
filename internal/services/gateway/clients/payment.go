package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/silazemli/lab2-template/internal/services/payment"
)

type PaymentClient struct {
	client  HTTPClient
	baseURL string
}

func NewPaymentClient(client HTTPClient, baseURL string) *PaymentClient {
	return &PaymentClient{
		client:  client,
		baseURL: baseURL,
	}
}

func (paymentClient *PaymentClient) CreatePayment(thePayment payment.Payment) error {
	URL := paymentClient.baseURL
	body, err := json.Marshal(thePayment)
	if err != nil {
		fmt.Println("failed to unmarshal")
		return fmt.Errorf("failed to build request body: %w", err)
	}
	request, err := http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("failed to build")
		return fmt.Errorf("failed to build request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	responce, err := paymentClient.client.Do(request)
	if err != nil {
		fmt.Println("failed to make request")
		return fmt.Errorf("failed to make request: %w", err)
	}
	fmt.Println(err)
	fmt.Println(responce.StatusCode)
	switch responce.StatusCode {
	case http.StatusCreated:
		return nil
	case http.StatusBadRequest, http.StatusInternalServerError:
		fmt.Println("server error")
		return fmt.Errorf("server error: %w", err)
	default:
		fmt.Println("wtf")
		fmt.Println(err)
		return fmt.Errorf("unknown error: %w", err)
	}
}

func (paymentClient *PaymentClient) CancelPayment(paymentUID string) error {
	URL := fmt.Sprintf("%s/%s", paymentClient.baseURL, paymentUID)
	request, err := http.NewRequest(http.MethodPatch, URL, nil)
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}

	response, err := paymentClient.client.Do(request)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	switch response.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusBadRequest, http.StatusInternalServerError:
		return fmt.Errorf("server error: %w", err)
	default:
		return fmt.Errorf("unknown error: %w", err)
	}
}

func (paymentClient PaymentClient) GetPayment(paymentUID string) (payment.Payment, error) {
	URL := fmt.Sprintf("%s/%s", paymentClient.baseURL, paymentUID)
	request, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return payment.Payment{}, fmt.Errorf("failed to build request: %w", err)
	}

	response, err := paymentClient.client.Do(request)
	if err != nil {
		return payment.Payment{}, fmt.Errorf("failed to make request: %w", err)
	}
	switch response.StatusCode {
	case http.StatusOK:
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return payment.Payment{}, fmt.Errorf("failed to read response body: %w", err)
		}
		var thePayment payment.Payment
		if err := json.Unmarshal(body, &thePayment); err != nil {
			return payment.Payment{}, fmt.Errorf("failed to unmarshal response: %w", err)
		}
		return thePayment, nil
	case http.StatusBadRequest, http.StatusInternalServerError:
		return payment.Payment{}, fmt.Errorf("server error: %w", err)
	default:
		return payment.Payment{}, fmt.Errorf("unknown error: %w", err)
	}
}
