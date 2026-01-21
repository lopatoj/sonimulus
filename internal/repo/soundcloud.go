package repo

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"lopa.to/sonimulus/env"
	"lopa.to/sonimulus/soundcloud"
)

type SoundCloudRepository struct {
	client *soundcloud.ClientWithResponses
	env    env.Env
}

func NewSoundCloudRepository(client *soundcloud.ClientWithResponses, env env.Env) *SoundCloudRepository {
	return &SoundCloudRepository{client: client, env: env}
}

func (scr *SoundCloudRepository) GetMe(ctx context.Context) (*soundcloud.Me, error) {
	res, err := scr.client.GetMeWithResponse(ctx)
	if err != nil {
		slog.Error("Failed to get user", "error", err)
		return nil, err
	}
	if res.StatusCode() == http.StatusOK {
		return res.ApplicationjsonCharsetUtf8200, nil
	}
	return nil, errors.New("unauthorized")
}
