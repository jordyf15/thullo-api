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

	// cover errors
	ErrMalformedCover             = newErr(401, "Cover is malformed")
	ErrUnknownUnsplashPhotoID     = newErr(402, "Unknown unsplash photo ID")
	ErrInvalidCoverSource         = newErr(403, "Source of cover is invalid")
	ErrUnsplashFocalPointYTooHigh = newErr(404, "Unsplash focal point Y is too high")
	ErrUnsplashFocalPointYTooLow  = newErr(405, "Unsplash focal point Y is too low")

	// Board errors
	ErrBoardTitleEmpty          = newErr(501, "Board title is empty")
	ErrBoardCoverEmpty          = newErr(502, "Board cover is empty")
	ErrBoardInvalidVisibility   = newErr(503, "Board visibility is invalid")
	ErrUserIsAlreadyBoardMember = newErr(504, "User is already a board member")
	ErrInvalidBoardMemberRole   = newErr(505, "Board member role is invalid")
	ErrBoardMustHaveAnAdmin     = newErr(506, "Board must have atleast one admin")

	// list errors
	ErrListTitleEmpty      = newErr(601, "List title is empty")
	ErrListPositionTooLow  = newErr(602, "List position is too low")
	ErrListPositionTooHigh = newErr(603, "List position is too high")

	// card errors
	ErrCardTitleEmpty = newErr(701, "Card title is empty")

	// comment errors
	ErrCommentEmpty = newErr(801, "Comment is empty")
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
