package custom_errors

import "strings"

var (
	// general errors
	ErrUnknownErrorOccured = newErr(101, "Unknown error occured")
	ErrInvalidIDInPath     = newErr(102, "Invalid ID in path")
	ErrRecordNotFound      = newErr(103, "Record not found")
	ErrNotAuthorized       = newErr(104, "You are not authorized to perform this action")

	// User Errors
	ErrCurrentPasswordWrong          = newErr(201, "Wrong current password")
	ErrEmailAddressInvalid           = newErr(202, "Invalid email address")
	ErrEmailAddressAlreadyRegistered = newErr(203, "Email address already registered")
	ErrUsernameTooShort              = newErr(204, "Username is too short")
	ErrUsernameTooLong               = newErr(205, "Username is too long")
	ErrUsernameAlreadyExists         = newErr(206, "Username already exist")
	ErrUsernameInvalid               = newErr(207, "Username is invalid")
	ErrNameTooShort                  = newErr(208, "Name is too short")
	ErrNameTooLong                   = newErr(209, "Name is too long")
	ErrPasswordTooShort              = newErr(210, "Password is too short")
	ErrPasswordTooLong               = newErr(211, "Password is too long")
	ErrPasswordInvalid               = newErr(212, "Password is invalid")
	ErrImageFormatInvalid            = newErr(213, "Image format must be in JPEG format")
	ErrImageSizeTooLarge             = newErr(214, "Image size is too large")

	// token errors
	ErrMalformedRefreshToken   = newErr(301, "Refresh token is malformed")
	ErrInvalidRefreshToken     = newErr(302, "Invalid refresh token")
	ErrRefreshTokenNotFound    = newErr(303, "Refresh token not found")
	ErrMalformedAccessToken    = newErr(304, "Access token is malformed")
	ErrInvalidAccessToken      = newErr(305, "Invalid access token")
	ErrAccessTokenExpired      = newErr(306, "Access token expired")
	ErrGoogleOauthTokenExpired = newErr(307, "Google oauth token expired")

	// board errors
	ErrTitleEmpty                 = newErr(401, "Title is empty")
	ErrCoverEmpty                 = newErr(402, "Cover is empty")
	ErrMalformedCover             = newErr(403, "Cover is malformed")
	ErrUnknownUnsplashPhotoID     = newErr(404, "Unknown unsplash photo ID")
	ErrInvalidCoverSource         = newErr(405, "Source of cover is invalid")
	ErrUnsplashFocalPointYTooHigh = newErr(406, "Unsplash focal point Y is too high")
	ErrUnsplashFocalPointYTooLow  = newErr(407, "Unsplash focal point Y is too low")
	ErrInvalidVisibility          = newErr(408, "Visibility is invalid")
)

type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (err *Error) Error() string {
	return err.Message
}

func newErr(code int, message string) *Error {
	return &Error{Message: message, Code: code}
}

type MultipleErrors struct {
	Errors []error `json:"errors"`
}

func (multipleErr *MultipleErrors) Error() string {
	messages := make([]string, len(multipleErr.Errors))
	for i, error := range multipleErr.Errors {
		messages[i] = error.Error()
	}

	return strings.Join(messages, ", ")
}
