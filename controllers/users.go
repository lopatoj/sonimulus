package controllers

// PublicUser represents a user without sensitive information.
type PublicUser struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UsersController represents a controller for managing users.
type UsersController struct {
	users UsersRepository
}

// NewUsersController creates a new UsersController.
func NewUsersController(users UsersRepository) *UsersController {
	return &UsersController{
		users: users,
	}
}

// GetUserById retrieves a user by ID.
func (uc *UsersController) GetUserById(id string) (PublicUser, error) {
	user, err := uc.users.GetUserById(id)
	if err != nil {
		return PublicUser{}, err
	}
	if user == nil {
		return PublicUser{}, ErrUserNotFound
	}

	return PublicUser{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

// GetUserByName retrieves a user by name.
func (uc *UsersController) GetUserByName(name string) (PublicUser, error) {
	user, err := uc.users.GetUserByName(name)
	if err != nil {
		return PublicUser{}, err
	}
	if user == nil {
		return PublicUser{}, ErrUserNotFound
	}

	return PublicUser{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
