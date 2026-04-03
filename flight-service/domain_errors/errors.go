package domain_errors

import "errors"

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrInvalidDate         = errors.New("invalid date format provided")
	ErrFlightNotFound      = errors.New("no flights found matching criteria")
	ErrDatabaseQuery       = errors.New("database query failed unexpectedly")
	ErrInvalidID           = errors.New("invalid or malformed instance ID provided")
	ErrAirportNotFound     = errors.New("airport not found")
)
