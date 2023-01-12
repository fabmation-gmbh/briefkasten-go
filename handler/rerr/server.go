package rerr

import "net/http"

// InternalServerError describes a problem on server side.
var InternalServerError = newErr(
	ecInternalServerError,
	EServer,
	"internal server error",
	http.StatusInternalServerError,
)
