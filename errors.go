package mediabrowser

import (
	"errors"
	"fmt"
)

// most 404 errors are from UserNotFound, so this generic error doesn't really need any detail.
type ErrNotFound error

// DetailedError is sometimes returned by Jellyfin.
// These details are only populated when mb.Verbose is true.
type DetailedError struct {
	Verbose bool
	Code    int
	Data    string `json:"data"`
	/*
	   Code     int
	   Type     string `json:"type"`
	   Title    string `json:"title"`
	   Detail   string `json:"detail"`
	   Instance string `json:"instance"
	*/
}

func (err DetailedError) Details() string {
	if !err.Verbose {
		return ""
	}
	msg := err.Data
	/*if err.Type != "" {
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
	}*/
	return msg + "\n"
}

type ErrUnauthorized struct {
	DetailedError
}

func (err ErrUnauthorized) Error() string {
	msg := fmt.Sprintf("%d Unauthorized, check credentials.", err.Code)
	if err.Verbose {
		msg += "\n" + err.Details()
	}
	return msg
}

type ErrForbidden struct {
	DetailedError
}

func (err ErrForbidden) Error() string {
	msg := "forbidden, the user may not have the correct permissions."
	if err.Verbose {
		msg += "\n" + err.Details()
	}
	return msg
}

var (
	NotFound ErrNotFound = errors.New("resource not found")
)

type ErrUnknown struct {
	DetailedError
}

func (err ErrUnknown) Error() string {
	msg := fmt.Sprintf("failed (code %d)", err.Code)
	if err.Verbose {
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
		err = ErrUnauthorized{
			DetailedError: DetailedError{
				Data:    data,
				Verbose: mb.Verbose,
				Code:    status,
			},
		}
	case 404:
		err = NotFound
		return
	case 403:
		err = ErrForbidden{
			DetailedError: DetailedError{
				Data:    data,
				Verbose: mb.Verbose,
				Code:    status,
			},
		}
	default:
		err = ErrUnknown{
			DetailedError: DetailedError{
				Data:    data,
				Verbose: mb.Verbose,
				Code:    status,
			},
		}
	}
	/*if mb.Verbose && data != "" {
		json.Unmarshal([]byte(data), &err)
	}*/
	return err
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
	if err.Verbose {
		msg += "\n" + err.Details()
	}
	return msg
}

type ErrNoPolicySupplied struct {
	DetailedError
}

func (err ErrNoPolicySupplied) Error() string {
	msg := "No (valid) policy was given."
	if err.Verbose {
		msg += "\n" + err.Details()
	}
	return msg
}
