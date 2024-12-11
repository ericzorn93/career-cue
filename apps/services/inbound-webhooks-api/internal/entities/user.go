package entities

// User is the internal domain representation
// of an authenticated user
type User struct {
	FirstName            string
	LastName             string
	Nickname             string
	Username             string
	EmailAddress         string
	EmailAddressVerified bool
	PhoneNumber          string
	PhoneNumberVerified  bool
	Strategy             string
	Metadata             map[string]any
}

// UserOption allows us to configure User
type UserOption func(*User)

// WithUserFirstName adds the user's first name to the struct
func WithUserFirstName(firstName string) UserOption {
	return func(u *User) {
		u.FirstName = firstName
	}
}

// WithUserLastName adds the user's last name to the struct
func WithUserLastName(lastName string) UserOption {
	return func(u *User) {
		u.LastName = lastName
	}
}

// NewUser is the constructor for a User struct
func NewUser(opts ...UserOption) User {
	u := User{}

	for _, opt := range opts {
		opt(&u)
	}

	return u
}

// GetFullName will return the concatonated user's
// first and last name
func (u User) GetFullName() string {
	return u.FirstName + " " + u.LastName
}
