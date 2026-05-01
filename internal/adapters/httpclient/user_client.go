package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bikes2road/authentication/internal/domain"
	"github.com/bikes2road/authentication/internal/ports"
)

type userServiceClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewUserServiceClient crea una nueva instancia del cliente del servicio de usuarios
func NewUserServiceClient(baseURL string) ports.UserServiceClient {
	return &userServiceClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetUserByEmail obtiene un usuario por su email desde el servicio de usuarios
func (c *userServiceClient) GetUserByEmailOrNickName(ctx context.Context, emailOrNickName, password string) (*domain.User, error) {
	endpoint := fmt.Sprintf("%s/v1/verify-credentials", c.baseURL)

	requestBody := struct {
		EmailOrNickName string `json:"email_or_nick_name"`
		Password        string `json:"password"`
	}{
		EmailOrNickName: emailOrNickName,
		Password:        password,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Println(err)
		return nil, domain.ErrUserServiceUnavailable
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, domain.ErrUserNotFound
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, domain.ErrInvalidCredentials
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var user domain.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &user, nil
}

// GetUserByID obtiene un usuario por su ID desde el servicio de usuarios
func (c *userServiceClient) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	endpoint := fmt.Sprintf("%s/v1/%s", c.baseURL, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, domain.ErrUserServiceUnavailable
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, domain.ErrUserNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var user domain.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &user, nil
}
