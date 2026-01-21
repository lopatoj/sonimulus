package data

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"lopa.to/sonimulus/soundcloud"
)

// NewSoundCloudClient creates a new SoundCloud client.
func NewSoundCloudClient(serverURL string) (*soundcloud.ClientWithResponses, error) {
	client, err := soundcloud.NewClientWithResponses(serverURL,
		soundcloud.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			if token, ok := ctx.Value("access_token").(string); ok && token != "" {
				req.Header.Set("Authorization", "OAuth "+token)
				return nil
			}

			err := errors.New("failed to extract access token from context")
			slog.Error("Failed to extract access token from context")
			return err
		}),
	)
	if err != nil {
		return nil, err
	}

	slog.Info("SoundCloud client initialized")

	return client, nil
}
