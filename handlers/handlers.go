package handlers

import (
	"lopa.to/sonimulus/controllers"
	"lopa.to/sonimulus/env"
)

type Handler struct {
	auth  *controllers.AuthController
	users *controllers.UsersController
	env   env.Env
}

func NewHandler(auth *controllers.AuthController, users *controllers.UsersController, e env.Env) *Handler {
	return &Handler{
		auth:  auth,
		users: users,
		env:   e,
	}
}
