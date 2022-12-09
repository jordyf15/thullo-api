package custom_errors

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
