package rerr

import "net/http"

// RequestMalformed describes that the request data (body, header, URI, ...) is malformed.
//
// NOTE: This error shall ONLY be used if information disclousure is acceptable.
var RequestMalformed = newErr(
	ecRequestMalformed,
	ERequest,
	"invalid request body/ header",
	http.StatusBadRequest,
)

// Unauthenticated describes that the client must authenticate itself
// before proceeding.
var Unauthenticated = newErr(
	ecUnauthenticated,
	ERequest,
	"unauthenticated",
	http.StatusBadRequest,
)
