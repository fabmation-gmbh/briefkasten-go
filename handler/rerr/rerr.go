package rerr

import "fmt"

// ErrorType describes the type of the error, i.e., "OAuthException".
type ErrorType uint32

const (
	// EPermission errors are related to the permission system.
	EPermission ErrorType = iota + 1
	// ERequest errors are related to malformed requests.
	ERequest
	// EServer errors are related to a problem on server side.
	EServer
	// ENoErr errors are no real errors, only specific structured information.
	ENoErr
)

// ErrorCode describes the error uniquely.
type ErrorCode uint32

const (
	ecRequestMalformed    = ErrorCode(100)
	ecInternalServerError = ErrorCode(101)
	ecUnauthenticated     = ErrorCode(102)
)

// Error is an rerr (request/ REST API) error.
type Error struct {
	ErrType  ErrorType `json:"type"`
	Message  string    `json:"message"`
	Code     ErrorCode `json:"code"`
	httpCode int       `json:"-"`

	err    error
	logMsg string
}

// Unwrap returns the wrapped error, if any error has been wrapped.
func (e Error) Unwrap() error { return e.err }

// With allows to wrap an error within the defined one.
func (e Error) With(err error) Error {
	e.err = err

	return e
}

// WithLogMsg allows to add an optional (internal) log message to the error.
func (e Error) WithLogMsg(msg string) Error {
	e.logMsg = msg

	return e
}

func (e Error) LogMsg() string { return e.logMsg }

// Error implements the error interface.
func (e Error) Error() string {
	return fmt.Sprintf("[%d of %d] %s", e.Code, e.ErrType, e.Message)
}

func (e Error) HTTPCode() int { return e.httpCode }

func (e Error) isRerrError() {}

func (e Error) Is(err error) bool {
	_, ok := err.(interface{ isRerrError() })

	return ok
}

func newErr(id ErrorCode, group ErrorType, msg string, httpCode int) Error {
	return Error{
		Code:     id,
		ErrType:  group,
		Message:  msg,
		httpCode: httpCode,
	}
}

//go:generate enumer -type=ErrorType -trimprefix="E" -json
