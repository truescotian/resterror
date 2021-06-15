// This is an example scenario from the perspective of each consumer role. In these examples,
// we use a UserService interface defined in the application domain.

package main

// UserService represents a service for managing users.
type UserService interface {
	// Returns a user by ID.
	FindUserByID(id int) (*User, error)

	// Creates a new user.
	CreateUser(user *User) error
}

/*
  APPLICATION ROLE EXAMPLE

  Application role is typically concerned with simple error codes. For example, if our program
  attempts to fetch a User by ID, and it receives a "not found" error, it could re-attempt
  by searching by an email address.

  TODO: the following should be in a package specific to the postgres db.
  This example is just to show that we can return ENOTFOUND for our application
  to operate on independent of the implementation of UserService.
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
		  return nil, &Error{Code: ENOTFOUND}
	  } else if err {
		  return nil, err
	  }
	  return &user, nil
  }

  func main() {
    user, err := userService.FindUserByID(ctx, 100)
    if ErrorCode(err) == ENOTFOUND {
      // retry another method of finding user
    } else if err != nil {
      return err
    }
  }
*/

/*
  END USER ROLE

  End users expect actionable, human-readable messages. These can have additional
  constraints such as branding tone or internationalization but we'll focus on the basics.

  An example of end user messaging is field validation. Here we check to ensure that new
  users in our UserService have a username and it is unique.

  // CreateUser creates a new user in the system.
  // Returns EINVALID if the username is blank or already exists.
  // Returns ECONFLICT if the username is already in use.
  func (s *UserService) CreateUser(ctx context.Context, user *User) error {
	  // Validate username is non-blank.
	  if user.Username == "" {
		  return &Error{Code: EINVALID, Message: "Username is required."}
	  }

	  // Verify user does not already exist.
	  if s.usernameInUse(user.Username) {
		  return &Error{
			  Code: ECONFLICT,
			  Message: "Username is already in use. Please choose a different username.",
		  }
	  }

	  ...
  }
*/

/*
  OPERATOR ROLE

  We need to be able to provide all of the previous roles information plus a logical stack
  trace to our operator so they can debug issues. Go already provides a simple method,
  error.Error() to print error information, so we can utilize that.

  Logical Stack traces:
  Dumps a list of every function in the call stack from where an error occured. You can see
  this at work when you call panic(). These stack traces can be overwhelming and we only
  need a small subset of those lines to understand the context of our error. A logical
  stack trace contains only the layers that we as developers find to be important
  in describing the program flow. We will accomplish this using the Op and Err fields to wrap
  errors to provide context.

  Suppose we need to create additional roles for our new users. We can utilize Op and Err fields in our
  application Error to wrap this nested functionality.

  // CreateUser creates a new user in the system with a default role.
  func (s *UserService) CreateUser(ctx context.Context, user *myapp.User) error {
	  const op = "UserService.CreateUser"

	  // Perform validation...

	  // Insert user record.
	  if err := s.insertUser(ctx, user); err != nil {
		  return &Error{Op: op, Err: err}
	  }

	  // Insert default role.
	  if err := s.attachRole(ctx, user.ID, "default"); err != nil {
		  return &Error{Op: op, Err: err}
	  }
	  return nil
  }

  // insertUser inserts the user into the database.
  func (s *UserService) insertUser(ctx context.Context, user *myapp.User) error {
	  const op = "insertUser"
	  if _, err := s.db.Exec(`INSERT INTO users...`); err != nil {
		  return &Error{Op: op, Err: err}
	  }
	  return nil
  }

  // attachRole inserts a role record for a user in the database
  func (s *UserService) attachRole(ctx context.Context, id int, role string) error {
	  const op = "attachRole"
	  if _, err := s.db.Exec(`INSERT roles...`); err != nil {
		  return &Error{Op: op, Err: err}
	  }
	  return nil
  }

  Let’s assume we receive a Postgres syntax error inside our attachRole()
  function and Postgres returns the message: syntax error at or near "INSERT"

  Without context, we don’t know if this occurred in our insertUser()
  function or our attachRole() function. This is very  tedious to debug
  when your API executes 20+ SQL queries.

  However, because we are wrapping our errors, Error() will return:

  UserService.CreateUser: attachRole: syntax error at or near "INSERT"

  This lets us narrow down the errant query and make the fix.

*/
