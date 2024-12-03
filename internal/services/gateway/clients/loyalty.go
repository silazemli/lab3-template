package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/silazemli/lab2-template/internal/services/loyalty"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type LoyaltyClient struct {
	client  HTTPClient
	baseURL string
}

func NewLoyaltyClient(client HTTPClient, baseURL string) *LoyaltyClient {
	return &LoyaltyClient{
		client:  client,
		baseURL: baseURL,
	}
}

func (loyaltyClient *LoyaltyClient) GetUser(username string) (loyalty.Loyalty, error) {
	URL := fmt.Sprintf("%s/%s", loyaltyClient.baseURL, "me")
	request, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return loyalty.Loyalty{}, fmt.Errorf("failed to build request: %w", err)
	}
	request.Header.Set("X-User-Name", username)
	responce, err := loyaltyClient.client.Do(request)
	if err != nil {
		return loyalty.Loyalty{}, fmt.Errorf("failed to make request: %w", err)
	}
	body, err := io.ReadAll(responce.Body)
	if err != nil {
		return loyalty.Loyalty{}, fmt.Errorf("failed to read response body: %w", err)
	}
	defer responce.Body.Close()
	switch responce.StatusCode {
	case http.StatusOK:
		var user loyalty.Loyalty
		if err := json.Unmarshal(body, &user); err != nil {
			return loyalty.Loyalty{}, fmt.Errorf("failed to unmarshal response body: %w", err)
		}
		return user, nil
	case http.StatusInternalServerError, http.StatusBadRequest:
		return loyalty.Loyalty{}, fmt.Errorf("server error: %w", err)
	case http.StatusNotFound:
		return loyalty.Loyalty{}, fmt.Errorf("not found %w", err)
	default:
		return loyalty.Loyalty{}, fmt.Errorf("unknown error: %w", err)
	}
}

func (loyaltyClient *LoyaltyClient) GetStatus(username string) (string, error) {
	URL := loyaltyClient.baseURL
	fmt.Sprintln(URL)
	request, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return "UNKNOWN", fmt.Errorf("failed to build request: %w", err)
	}
	request.Header.Set("X-User-Name", username)
	responce, err := loyaltyClient.client.Do(request)
	if err != nil {
		return "UNKNOWN", fmt.Errorf("failed to make request: %w", err)
	}

	switch responce.StatusCode {
	case http.StatusOK:
		body, err := io.ReadAll(responce.Body)
		if err != nil {
			return "UNKNOWN", fmt.Errorf("failed to read response body: %w", err)
		}

		var model struct {
			Status string `json:"status"`
		}
		if err := json.Unmarshal(body, &model); err != nil {
			return "UNKNOWN", fmt.Errorf("failed to unmarshal response body: %w", err)
		}
		status := model.Status
		return status, nil

	case http.StatusInternalServerError, http.StatusBadRequest:
		return "UNKNOWN", fmt.Errorf("server error: %w", err)
	case http.StatusNotFound:
		return "UNKNOWN", fmt.Errorf("not found: %w", err)
	default:
		return "UNKNOWN", fmt.Errorf("unknown error: %w", err)
	}
}

func (loyaltyClient *LoyaltyClient) DecrementCounter(username string) error {
	URL := fmt.Sprintf("%s/%s", loyaltyClient.baseURL, "decrement")
	request, err := http.NewRequest(http.MethodPatch, URL, nil)
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}
	request.Header.Set("X-User-Name", username)
	responce, err := loyaltyClient.client.Do(request)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}

	switch responce.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusInternalServerError, http.StatusNotFound, http.StatusBadRequest:
		return fmt.Errorf("server error: %w", err)
	default:
		return fmt.Errorf("unknown error: %w", err)
	}
}

func (loyaltyClient *LoyaltyClient) IncrementCounter(username string) error {
	URL := fmt.Sprintf("%s/%s", loyaltyClient.baseURL, "increment")
	request, err := http.NewRequest(http.MethodPatch, URL, nil)
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}
	request.Header.Set("X-User-Name", username)
	responce, err := loyaltyClient.client.Do(request)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}

	switch responce.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusInternalServerError, http.StatusNotFound, http.StatusBadRequest:
		return fmt.Errorf("server error: %w", err)
	default:
		return fmt.Errorf("unknown error: %w", err)
	}
}
