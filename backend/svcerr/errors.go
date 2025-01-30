package svcerr

import "errors"

var (
	ErrStatusPageNotFound    = errors.New("status page was not found")
	ErrNoMonitorsFound       = errors.New("monitors were not found")
	ErrUserNotFound          = errors.New("user was not found")
	ErrNoIncidentsFound      = errors.New("incidents were not found")
	ErrNotificationsNotFound = errors.New("notifications channel was not found")
)

func IsNotFound(err error) bool {
	switch {
	case errors.Is(err, ErrNoMonitorsFound),
		errors.Is(err, ErrNotificationsNotFound),
		errors.Is(err, ErrUserNotFound),
		errors.Is(err, ErrNoIncidentsFound),
		errors.Is(err, ErrStatusPageNotFound):
		return true
	default:
		return false
	}
}
