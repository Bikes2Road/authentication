package supabase

import (
	"fmt"
	"log"

	supabase "github.com/supabase-community/supabase-go"
)

// ClientConfig holds the Supabase client configuration
type ClientConfig struct {
	URL    string
	APIKey string
}

// NewClient creates a new Supabase client connection
func NewClient(cfg ClientConfig) (*supabase.Client, error) {
	if cfg.URL == "" {
		return nil, fmt.Errorf("supabase URL is required")
	}
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("supabase API key is required")
	}

	client, err := supabase.NewClient(cfg.URL, cfg.APIKey, &supabase.ClientOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize supabase client: %w", err)
	}

	log.Println("Supabase client initialized successfully")

	return client, nil
}
