package customErrors

import (
	"errors"
	"strconv"
)

var (
	OOS           = errors.New("out of stock")
	EmptyField    = errors.New("invaild checkout")
	UnexpectedErr = errors.New("unexpected error")

	//Status codes
	Unauthorized        = errors.New("unexpected status code: 401 (Unauthorized)")
	Forbidden           = errors.New("Forbidden / Access denied (403)")
	NotFound            = errors.New("unexpected status code: 404 (Not Found)")
	Gone                = errors.New("unexpected status code: 410 (Gone)")
	TooManyRequests     = errors.New("unexpected status code: 429 (Too Many Requests)")
	InternalServerError = errors.New("unexpected status code: 500 (Internal Server Error)")
	ServiceUnavailable  = errors.New("unexpected status code: 503 (Service Unavailable)")
)

func StatusCodeNotExpectedError(statusCode int) error {
	var err error
	switch statusCode {
	case 401:
		err = Unauthorized
	case 403:
		err = Forbidden
	case 404:
		err = NotFound
	case 410:
		err = Gone
	case 429:
		err = TooManyRequests
	case 500:
		err = InternalServerError
	case 503:
		err = ServiceUnavailable
	default:
		err = errors.New("unexpected status code: " + strconv.Itoa(statusCode))
	}
	return err
}
