package domain

import (
	"errors"
	"time"

	"github.com/xybor-x/snowflake"
	"github.com/xybor/todennus-backend/pkg/xstring"
	"golang.org/x/crypto/bcrypt"
)

const (
	MinimumDisplayNameLength = 3
	MaximumDisplayNameLength = 32

	MinimumUsernameLength = 6
	MaximumUsernameLength = 20

	MinimumPasswordLength = 8
	MaximumPassowrdLength = 32

	HashingCost = bcrypt.DefaultCost
)

type User struct {
	ID          int64
	DisplayName string
	Username    string
	HashedPass  string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UserDomain struct {
	Snowflake *snowflake.Node
}

func NewUserDomain(snowflake *snowflake.Node) (*UserDomain, error) {
	return &UserDomain{Snowflake: snowflake}, nil
}

func (domain *UserDomain) Create(username, password string) (User, error) {
	if err := domain.validateUsername(username); err != nil {
		return User{}, err
	}

	if err := domain.validatePassword(password); err != nil {
		return User{}, err
	}

	hasedPass, err := bcrypt.GenerateFromPassword([]byte(password), HashingCost)
	if err != nil {
		return User{}, Wrap(ErrUnknownCritical, err.Error())
	}

	return User{
		ID:          domain.Snowflake.Generate().Int64(),
		DisplayName: username,
		Username:    username,
		HashedPass:  string(hasedPass),
	}, nil
}

func (domain *UserDomain) Validate(hashedPassword, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}

		return false, Wrap(ErrUnknownRecoverable, err.Error())
	}

	return true, nil
}

func (domain *UserDomain) SetDisplayName(user *User, displayname string) error {
	if err := domain.validateDisplayName(displayname); err != nil {
		return err
	}

	user.DisplayName = displayname
	return nil
}

func (domain *UserDomain) validateDisplayName(displayname string) error {
	if len(displayname) > MaximumDisplayNameLength {
		return Wrap(ErrInvalidDisplayName, "require at most %d characters", MaximumDisplayNameLength)
	}

	if len(displayname) < MinimumDisplayNameLength {
		return Wrap(ErrInvalidDisplayName, "require at least %d characters", MinimumDisplayNameLength)
	}

	for _, c := range displayname {
		if !xstring.IsLetter(c) && !xstring.IsUnderscore(c) && !xstring.IsSpace(c) {
			return Wrap(ErrInvalidUsername, "got an invalid character %c", c)
		}
	}

	return nil

}

func (domain *UserDomain) validateUsername(username string) error {
	if len(username) > MaximumUsernameLength {
		return Wrap(ErrInvalidUsername, "require at most %d characters", MaximumUsernameLength)
	}

	if len(username) < MinimumUsernameLength {
		return Wrap(ErrInvalidUsername, "require at least %d characters", MinimumUsernameLength)
	}

	for _, c := range username {
		if !xstring.IsLetter(c) && !xstring.IsUnderscore(c) {
			return Wrap(ErrInvalidUsername, "got an invalid character %c", c)
		}
	}

	return nil
}

func (domain *UserDomain) validatePassword(password string) error {
	if len(password) > MaximumPassowrdLength {
		return Wrap(ErrInvalidPassword, "require at most %d characters", MaximumPassowrdLength)
	}

	if len(password) < MinimumPasswordLength {
		return Wrap(ErrInvalidPassword, "require at least %d characters", MinimumPasswordLength)
	}

	haveLowercase := false
	haveUppercase := false
	haveNumber := false
	haveSpecial := false

	for _, c := range password {
		switch {
		case xstring.IsLowerCaseLetter(c):
			haveLowercase = true
		case xstring.IsUpperCaseLetter(c):
			haveUppercase = true
		case xstring.IsNumber(c):
			haveNumber = true
		case xstring.IsSpecialCharacter(c):
			haveSpecial = true
		default:
			return Wrap(ErrInvalidPassword, "got an invalid character %c", c)
		}
	}

	if !haveLowercase {
		return Wrap(ErrInvalidPassword, "require at least a lowercase letter")
	}

	if !haveUppercase {
		return Wrap(ErrInvalidPassword, "require at least an uppercase letter")
	}

	if !haveNumber {
		return Wrap(ErrInvalidPassword, "require at least a number")
	}

	if !haveSpecial {
		return Wrap(ErrInvalidPassword, "require at least a special character")
	}

	return nil
}
