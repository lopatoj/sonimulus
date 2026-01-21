package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"lopa.to/sonimulus/env"
	"lopa.to/sonimulus/internal/repo"
	"lopa.to/sonimulus/soundcloud"
)

var (
	ErrStateNotFound = errors.New("State not found")
)

type StateStorer interface {
	CreateState(ctx context.Context, verifier string) (state string, err error)
	GetVerifier(ctx context.Context, state string) (verifier string, found bool, err error)
}

type SessionStorer interface {
	CreateSession(ctx context.Context, sessionData repo.Session) (sessionID string, err error)
	GetSession(ctx context.Context, sessionID string) (session repo.Session, found bool, err error)
	DeleteSession(ctx context.Context, sessionID string) (found bool, err error)
}

type SoundCloudProvider interface {
	GetMe(ctx context.Context) (*soundcloud.Me, error)
}

type UserProvider interface {
	FindByKey(ctx context.Context, key repo.UserKey, value string) (user repo.User, found bool, err error)
	Create(ctx context.Context, id int64, username string) (user repo.User, err error)
}

// AuthController handles authentication-related operations.
type AuthController struct {
	env      env.Env
	config   oauth2.Config
	states   StateStorer
	sessions SessionStorer
	sc       SoundCloudProvider
	users    UserProvider
}

// NewAuthController creates a new instance of AuthController.
func NewAuthController(
	e env.Env,
	stateRepo StateStorer,
	sessionRepo SessionStorer,
	scRepo SoundCloudProvider,
	userRepo UserProvider,
) *AuthController {
	c := oauth2.Config{
		ClientID:     e.Soundcloud.ClientID,
		ClientSecret: e.Soundcloud.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  e.Soundcloud.AuthURL,
			TokenURL: e.Soundcloud.TokenURL,
		},
		RedirectURL: fmt.Sprintf(
			"%s:%d%s%s",
			e.Server.URL,
			e.Server.Port,
			e.Server.Route,
			e.Server.RedirectRoute,
		),
	}

	return &AuthController{
		env:      e,
		config:   c,
		states:   stateRepo,
		sessions: sessionRepo,
		sc:       scRepo,
		users:    userRepo,
	}
}

func (ac *AuthController) CreateAuthURL(ctx context.Context) (string, error) {
	verifier := oauth2.GenerateVerifier()
	state, err := ac.states.CreateState(ctx, verifier)
	if err != nil {
		slog.Error("Saving state to Redis", "error", err)
		return "", err
	}

	opt := oauth2.S256ChallengeOption(verifier)
	return ac.config.AuthCodeURL(state, opt), nil
}

func (ac *AuthController) ObtainToken(
	ctx context.Context,
	code string,
	state string,
) (string, error) {
	verifier, found, err := ac.states.GetVerifier(ctx, state)
	if err != nil {
		slog.Error("Retrieving state from Redis", "error", err)
		return "", err
	}

	if !found {
		slog.Error("State not found in Redis")
		return "", ErrStateNotFound
	}

	opt := oauth2.VerifierOption(verifier)
	token, err := ac.config.Exchange(ctx, code, opt)
	if err != nil {
		slog.Error("Exchanging code for token", "error", err)
		return "", err
	}

	ctx = context.WithValue(ctx, "access_token", token.AccessToken)
	me, err := ac.sc.GetMe(ctx)
	if err != nil {
		slog.Error("Getting user from provider", "error", err)
		return "", err
	}

	id, err := idFromUrn(*me.Urn)
	if err != nil {
		slog.Error("Parsing user ID from URN", "error", err)
		return "", err
	}

	slog.Info("Recieved user data", "id", id)

	user, found, err := ac.users.FindByKey(ctx, repo.UserKeyID, strconv.FormatInt(id, 10))
	if err != nil {
		slog.Error("Finding user", "error", err)
		return "", err
	}

	if !found {
		user, err = ac.users.Create(ctx, id, *me.Username)
		if err != nil {
			slog.Error("Creating user", "error", err)
			return "", err
		}
	}

	session := repo.Session{
		Token: token,
		User:  user,
	}

	sessionID, err := ac.sessions.CreateSession(ctx, session)
	if err != nil {
		slog.Error("Creating session", "error", err)
		return "", err
	}

	return sessionID, nil
}

func (ac *AuthController) GetSession(ctx context.Context, sessionID string) (session repo.Session, found bool, err error) {
	session, found, err = ac.sessions.GetSession(ctx, sessionID)
	if err != nil {
		slog.Error("Getting session", "error", err)
		return session, false, err
	}
	if !found {
		return session, false, nil
	}

	if session.Token.Expiry.Before(time.Now()) {
		token, err := ac.config.TokenSource(ctx, &oauth2.Token{
			RefreshToken: session.Token.RefreshToken,
		}).Token()
		if err != nil {
			slog.Error("Refreshing token", "error", err)
			return session, false, err
		}

		session.Token = token
	}

	return session, true, nil
}

func (ac *AuthController) DeleteSession(ctx context.Context, sessionID string) (bool, error) {
	return ac.sessions.DeleteSession(ctx, sessionID)
}

// idFromUrn is a helper for converting soundcloud user urns to ID integers.
//
// res, err := idFromUrn("soundcloud:user:1") // 1, nil
func idFromUrn(urn string) (int64, error) {
	id := urn[strings.LastIndex(urn, ":")+1:]
	return strconv.ParseInt(id, 10, 64)
}
