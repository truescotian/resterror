package main

// UserService represents a service for managing users.
type UserService interface {
	// Returns a user by ID.
	FindUserByID(id int) (*User, error)

	// Creates a new user.
	CreateUser(user *User) error
}

// TODO: the following should be in a package specific to the postgres db.
// This example is just to show that we can return ENOTFOUND for our application
// to operate on independent of the implementation of UserService.
/*
func (s *UserService) FindUserByID(id int) (*User, error) {
	var user myapp.User
	if err := s.QueryRowContext(ctx, `
		SELECT id, username
		FROM users
		WHERE id = $1
	`,
		id
	).Scan(
		&user.ID,
		&user.Username,
	); err == sql.ErrNoRows {
		return nil, &Error{Code: myapp.ENOTFOUND}
	} else if err {
		return nil, err
	}
	return &user, nil
}
*/
