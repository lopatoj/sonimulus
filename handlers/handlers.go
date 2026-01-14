package handlers

import (
	"lopa.to/sonimulus/config"
	"lopa.to/sonimulus/controllers"
)

type Handler struct {
	auth   *controllers.AuthController
	users  *controllers.UsersController
	config config.Config
}

func NewHandler(auth *controllers.AuthController, users *controllers.UsersController, config config.Config) *Handler {
	return &Handler{
		auth:   auth,
		users:  users,
		config: config,
	}
}
