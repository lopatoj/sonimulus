package controllers

import "lopa.to/sonimulus/repository"

type UsersRepository interface {
	Create(name, email string, salt, hash []byte) (*repository.User, error)
	GetUserByEmail(email string) (*repository.User, error)
	GetUserById(id string) (*repository.User, error)
	GetUserByName(name string) (*repository.User, error)
}
