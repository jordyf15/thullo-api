package models

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/jordyf15/thullo-api/custom_errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	minUsernameLength = 3
	maxUsernameLength = 30
	minNameLength     = 1
	maxNameLength     = 30
	minPasswordLength = 8
	maxPasswordLength = 30
)

var (
	passwordRegexes = []*regexp.Regexp{
		regexp.MustCompile(".*[A-Z].*"),          // must contain uppercase letter
		regexp.MustCompile(".*[a-z].*"),          // must contain lowercase letter
		regexp.MustCompile(".*[0-9].*"),          // must contain digit
		regexp.MustCompile(`.*[^A-Za-z0-9\s].*`), // must contain special character
	}
	emailRegex      = regexp.MustCompile("\\A[\\w+\\-.]+@[a-z\\d\\-.]+\\.[a-z]+\\z")
	usernameRegexes = []regexPair{
		{regex: regexp.MustCompile("^[a-z0-9._]+$"), shouldMatch: true},    // characters allowed
		{regex: regexp.MustCompile("^[^_.].+$"), shouldMatch: true},        // must not start with a fullstop or underscore
		{regex: regexp.MustCompile("\\A.*[_.]{2}.*$"), shouldMatch: false}, // must not have consecutive fullstops/unserscores
		{regex: regexp.MustCompile("\\A.*[^_.]$"), shouldMatch: true},      // must not end with a fullstop or underscore
	}
)

type regexPair struct {
	regex       *regexp.Regexp
	shouldMatch bool
}

type User struct {
	ID                primitive.ObjectID `bson:"_id" json:"id"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"encrypted_password" json:"-"`
	Password          string             `bson:"-" json:"-"`
	Username          string             `bson:"username" json:"username"`
	Name              string             `bson:"name" json:"name"`
	Bio               string             `bson:"bio" json:"bio"`
	Images            Images             `bson:"images" json:"images"`
	UpdatedAt         time.Time          `bson:"updated_at" json:"-"`
	CreatedAt         time.Time          `bson:"created_at" json:"created_at"`
}

func (user *User) Initials() string {
	names := strings.Split(user.Name, " ")
	switch len(names) {
	case 0:
		return ""
	case 1:
		return strings.ToUpper(user.Name[0:1] + user.Name[0:1])
	default:
		return strings.ToUpper(names[0][0:1] + names[1][0:1])
	}
}

func (user *User) SetPassword(newPassword string) error {
	if len(newPassword) < minPasswordLength {
		return custom_errors.ErrPasswordTooShort
	} else if len(newPassword) > 30 {
		return custom_errors.ErrPasswordTooLong
	}

	for _, regex := range passwordRegexes {
		if !regex.MatchString(newPassword) {
			return custom_errors.ErrPasswordInvalid
		}
	}

	hashedNewPassword, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	user.Password = ""
	user.EncryptedPassword = string(hashedNewPassword)

	return nil
}

func (user *User) VerifyFields() []error {
	user.Name = strings.Join(strings.Fields(user.Name), " ")
	user.Username = strings.ToLower(user.Username)

	for _, img := range user.Images {
		img.URL = ""
	}
	errors := make([]error, 0)

	if !emailRegex.MatchString(user.Email) {
		errors = append(errors, custom_errors.ErrEmailAddressInvalid)
	}

	for _, regexPair := range usernameRegexes {
		if regexPair.regex.MatchString(user.Username) != regexPair.shouldMatch {
			errors = append(errors, custom_errors.ErrUsernameInvalid)
			break
		}
	}

	if len(user.Username) < minUsernameLength {
		errors = append(errors, custom_errors.ErrUsernameTooShort)
	}

	if len(user.Username) > maxUsernameLength {
		errors = append(errors, custom_errors.ErrUsernameTooLong)
	}

	if len(user.Name) < minNameLength {
		errors = append(errors, custom_errors.ErrNameTooShort)
	}

	if len(user.Name) > maxNameLength {
		errors = append(errors, custom_errors.ErrNameTooLong)
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

func (user *User) ImagePath(image *Image) string {
	if len(image.ID) == 0 {
		return "uploads/users/default_profile_picture.png"
	}
	return fmt.Sprintf("uploads/users/%s/%s", user.ID.Hex(), image.ID)
}
