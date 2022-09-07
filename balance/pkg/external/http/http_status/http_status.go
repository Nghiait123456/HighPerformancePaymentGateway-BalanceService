package http_status

import (
	"github.com/valyala/fasthttp"
)

// HTTP status codes were stolen from net/http.
const (
	StatusContinue           = fasthttp.StatusContinue           // RFC 7231, 6.2.1
	StatusSwitchingProtocols = fasthttp.StatusSwitchingProtocols // RFC 7231, 6.2.2
	StatusProcessing         = fasthttp.StatusProcessing         // RFC 2518, 10.1
	StatusEarlyHints         = fasthttp.StatusEarlyHints         // RFC 8297

	StatusOK                   = fasthttp.StatusOK                   // RFC 7231, 6.3.1
	StatusCreated              = fasthttp.StatusCreated              // RFC 7231, 6.3.2
	StatusAccepted             = fasthttp.StatusAccepted             // RFC 7231, 6.3.3
	StatusNonAuthoritativeInfo = fasthttp.StatusNonAuthoritativeInfo // RFC 7231, 6.3.4
	StatusNoContent            = fasthttp.StatusNoContent            // RFC 7231, 6.3.5
	StatusResetContent         = fasthttp.StatusResetContent         // RFC 7231, 6.3.6
	StatusPartialContent       = fasthttp.StatusPartialContent       // RFC 7233, 4.1
	StatusMultiStatus          = fasthttp.StatusMultiStatus          // RFC 4918, 11.1
	StatusAlreadyReported      = fasthttp.StatusAlreadyReported      // RFC 5842, 7.1
	StatusIMUsed               = fasthttp.StatusIMUsed               // RFC 3229, 10.4.1

	StatusMultipleChoices   = fasthttp.StatusMultipleChoices   // RFC 7231, 6.4.1
	StatusMovedPermanently  = fasthttp.StatusMovedPermanently  // RFC 7231, 6.4.2
	StatusFound             = fasthttp.StatusFound             // RFC 7231, 6.4.3
	StatusSeeOther          = fasthttp.StatusSeeOther          // RFC 7231, 6.4.4
	StatusNotModified       = fasthttp.StatusNotModified       // RFC 7232, 4.1
	StatusUseProxy          = fasthttp.StatusUseProxy          // RFC 7231, 6.4.5
	_                       = 306                              // RFC 7231, 6.4.6 (Unused)
	StatusTemporaryRedirect = fasthttp.StatusTemporaryRedirect // RFC 7231, 6.4.7
	StatusPermanentRedirect = fasthttp.StatusPermanentRedirect // RFC 7538, 3

	StatusBadRequest                   = fasthttp.StatusBadRequest                   // RFC 7231, 6.5.1
	StatusUnauthorized                 = fasthttp.StatusUnauthorized                 // RFC 7235, 3.1
	StatusPaymentRequired              = fasthttp.StatusPaymentRequired              // RFC 7231, 6.5.2
	StatusForbidden                    = fasthttp.StatusForbidden                    // RFC 7231, 6.5.3
	StatusNotFound                     = fasthttp.StatusNotFound                     // RFC 7231, 6.5.4
	StatusMethodNotAllowed             = fasthttp.StatusMethodNotAllowed             // RFC 7231, 6.5.5
	StatusNotAcceptable                = fasthttp.StatusNotAcceptable                // RFC 7231, 6.5.6
	StatusProxyAuthRequired            = fasthttp.StatusProxyAuthRequired            // RFC 7235, 3.2
	StatusRequestTimeout               = fasthttp.StatusRequestTimeout               // RFC 7231, 6.5.7
	StatusConflict                     = fasthttp.StatusConflict                     // RFC 7231, 6.5.8
	StatusGone                         = fasthttp.StatusGone                         // RFC 7231, 6.5.9
	StatusLengthRequired               = fasthttp.StatusLengthRequired               // RFC 7231, 6.5.10
	StatusPreconditionFailed           = fasthttp.StatusPreconditionFailed           // RFC 7232, 4.2
	StatusRequestEntityTooLarge        = fasthttp.StatusRequestEntityTooLarge        // RFC 7231, 6.5.11
	StatusRequestURITooLong            = fasthttp.StatusRequestURITooLong            // RFC 7231, 6.5.12
	StatusUnsupportedMediaType         = fasthttp.StatusUnsupportedMediaType         // RFC 7231, 6.5.13
	StatusRequestedRangeNotSatisfiable = fasthttp.StatusRequestedRangeNotSatisfiable // RFC 7233, 4.4
	StatusExpectationFailed            = fasthttp.StatusExpectationFailed            // RFC 7231, 6.5.14
	StatusTeapot                       = fasthttp.StatusTeapot                       // RFC 7168, 2.3.3
	StatusMisdirectedRequest           = fasthttp.StatusMisdirectedRequest           // RFC 7540, 9.1.2
	StatusUnprocessableEntity          = fasthttp.StatusUnprocessableEntity          // RFC 4918, 11.2
	StatusLocked                       = fasthttp.StatusLocked                       // RFC 4918, 11.3
	StatusFailedDependency             = fasthttp.StatusFailedDependency             // RFC 4918, 11.4
	StatusUpgradeRequired              = fasthttp.StatusUpgradeRequired              // RFC 7231, 6.5.15
	StatusPreconditionRequired         = fasthttp.StatusPreconditionRequired         // RFC 6585, 3
	StatusTooManyRequests              = fasthttp.StatusTooManyRequests              // RFC 6585, 4
	StatusRequestHeaderFieldsTooLarge  = fasthttp.StatusRequestHeaderFieldsTooLarge  // RFC 6585, 5
	StatusUnavailableForLegalReasons   = fasthttp.StatusUnavailableForLegalReasons   // RFC 7725, 3

	StatusInternalServerError           = fasthttp.StatusInternalServerError           // RFC 7231, 6.6.1
	StatusNotImplemented                = fasthttp.StatusNotImplemented                // RFC 7231, 6.6.2
	StatusBadGateway                    = fasthttp.StatusBadGateway                    // RFC 7231, 6.6.3
	StatusServiceUnavailable            = fasthttp.StatusServiceUnavailable            // RFC 7231, 6.6.4
	StatusGatewayTimeout                = fasthttp.StatusGatewayTimeout                // RFC 7231, 6.6.5
	StatusHTTPVersionNotSupported       = fasthttp.StatusHTTPVersionNotSupported       // RFC 7231, 6.6.6
	StatusVariantAlsoNegotiates         = fasthttp.StatusVariantAlsoNegotiates         // RFC 2295, 8.1
	StatusInsufficientStorage           = fasthttp.StatusInsufficientStorage           // RFC 4918, 11.5
	StatusLoopDetected                  = fasthttp.StatusLoopDetected                  // RFC 5842, 7.2
	StatusNotExtended                   = fasthttp.StatusNotExtended                   // RFC 2774, 7
	StatusNetworkAuthenticationRequired = fasthttp.StatusNetworkAuthenticationRequired // RFC 6585, 6
)
