package entities

import "libs/backend/domain/user/valueobjects"

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
	CommonID             valueobjects.CommonID
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

// WithUserNickname adds the user's nickname to the struct
func WithUserNickname(nickname string) UserOption {
	return func(u *User) {
		u.Nickname = nickname
	}
}

// WithUserUsername adds the user's username to the struct
func WithUserUsername(username string) UserOption {
	return func(u *User) {
		u.Username = username
	}
}

// WithEmailAddress adds the user's email address to the struct
func WithEmailAddress(emailAddress string) UserOption {
	return func(u *User) {
		u.EmailAddress = emailAddress
	}
}

// WithEmailAddressVerified adds the user's email address verification status to the struct
func WithEmailAddressVerified(emailAddressVerified bool) UserOption {
	return func(u *User) {
		u.EmailAddressVerified = emailAddressVerified
	}
}

// WithPhoneNumber adds the user's phone number to the struct
func WithPhoneNumber(phoneNumber string) UserOption {
	return func(u *User) {
		u.PhoneNumber = phoneNumber
	}
}

// WithPhoneNumberVerified adds the user's phone number verification status to the struct
func WithPhoneNumberVerified(phoneNumberVerified bool) UserOption {
	return func(u *User) {
		u.PhoneNumberVerified = phoneNumberVerified
	}
}

// WithStrategy adds the user's authentication strategy to the struct
func WithStrategy(strategy string) UserOption {
	return func(u *User) {
		u.Strategy = strategy
	}
}

// WithCommonID adds the user's common ID to the struct (UUID)
func WithCommonID(commonID valueobjects.CommonID) UserOption {
	return func(u *User) {
		u.CommonID = commonID
	}
}

// WithMetadata adds the user's metadata to the struct
func WithMetadata(Metadata map[string]any) UserOption {
	return func(u *User) {
		u.Metadata = Metadata
	}
}

// NewUser is the constructor for a User struct
func NewUser(opts ...UserOption) User {
	u := User{
		Metadata: make(map[string]any),
	}

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