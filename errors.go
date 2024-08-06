package mediabrowser

import (
	"encoding/json"
	"errors"
	"fmt"
)

// most 404 errors are from UserNotFound, so this generic error doesn't really need any detail.
type ErrNotFound error

type ErrUnauthorized struct {
	DetailedError
}

func (err ErrUnauthorized) Error() string {
	msg := "Unauthorized, check credentials."
	if err.verbose {
		msg += "\n" + err.Details()
	}
	return msg
}

type ErrForbidden struct {
	DetailedError
}

func (err ErrForbidden) Error() string {
	msg := "forbidden, the user may not have the correct permissions."
	if err.verbose {
		msg += "\n" + err.Details()
	}
	return msg
}

var (
	NotFound ErrNotFound = errors.New("Resource not found.")
)

type ErrUnknown struct {
	code int
	DetailedError
}

func (err ErrUnknown) Error() string {
	msg := fmt.Sprintf("failed (code %d)", err.code)
	if err.verbose {
		msg += "\n" + err.Details()
	}
	return msg
}

func (mb *MediaBrowser) genericErr(status int, data string) (err error) {
	switch status {
	case 200, 204, 201:
		err = nil
		return
	case 401, 400:
		err = ErrUnauthorized{}
	case 404:
		err = NotFound
		return
	case 403:
		err = ErrForbidden{}
	default:
		err = ErrUnknown{code: status}
	}
	if mb.Verbose && data != "" {
		json.Unmarshal([]byte(data), &err)
	}
	return err
}

// DetailedError is sometimes returned by Jellyfin.
// These details are only populated when mb.Verbose is true.
type DetailedError struct {
	verbose  bool
	Type     string `json:"type"`
	Title    string `json:"title"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
}

func (err DetailedError) Details() string {
	msg := ""
	if err.Type != "" {
		msg += "Type: " + err.Type
	}
	if err.Title != "" {
		msg += "\nTitle: " + err.Title
	}
	if err.Detail != "" {
		msg += "\nDetail: " + err.Detail
	}
	if err.Instance != "" {
		msg += "\nInstance: " + err.Instance
	}
	return msg + "\n"
}

type ErrUserNotFound struct {
	user, id string
	DetailedError
}

func (err ErrUserNotFound) Error() string {
	msg := ""
	if err.user != "" {
		msg += "User \"" + err.user + "\" not found."
	} else {
		msg += "User with ID \"" + err.id + "\" not found."
	}
	if err.verbose {
		msg += "\n" + err.Details()
	}
	return msg
}

type ErrNoPolicySupplied struct {
	DetailedError
}

func (err ErrNoPolicySupplied) Error() string {
	msg := "No (valid) policy was given."
	if err.verbose {
		msg += "\n" + err.Details()
	}
	return msg
}
