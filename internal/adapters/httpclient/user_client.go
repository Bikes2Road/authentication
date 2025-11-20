package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
func (c *userServiceClient) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	// URL encode del email
	encodedEmail := url.QueryEscape(email)
	endpoint := fmt.Sprintf("%s/api/v1/users/email/%s", c.baseURL, encodedEmail)

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

// GetUserByID obtiene un usuario por su ID desde el servicio de usuarios
func (c *userServiceClient) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	endpoint := fmt.Sprintf("%s/api/v1/users/%s", c.baseURL, id)

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
