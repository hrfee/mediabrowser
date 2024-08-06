package mediabrowser

// Shared functions that work the same on Jellyfin & Emby.

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// UserByName returns the user corresponding to the provided username.
func (mb *MediaBrowser) UserByName(username string, public bool) (User, error) {
	username = strings.ToLower(username)
	var match User
	find := func() (User, error) {
		users, err := mb.GetUsers(public)
		if err != nil {
			return User{}, err
		}
		for _, user := range users {
			if strings.ToLower(user.Name) == username {
				return user, err
			}
		}
		return User{}, ErrUserNotFound{user: username}
	}
	match, err := find()
	if match.Name == "" {
		mb.CacheExpiry = time.Now()
		match, err = find()
	}
	return match, err
}

// SetPolicy sets the access policy for the user corresponding to the provided ID.
// No GetPolicy is provided because a User object includes Policy already.
func (mb *MediaBrowser) SetPolicy(userID string, policy Policy) error {
	url := fmt.Sprintf("%s/Users/%s/Policy", mb.Server, userID)
	DeNullPolicy(&policy)
	data, status, err := mb.post(url, policy, true)
	if status == 400 {
		err = ErrNoPolicySupplied{}
		if mb.Verbose {
			json.Unmarshal([]byte(data), &err)
		}
	} else if customErr := mb.genericErr(status, data); customErr != nil {
		err = customErr
	}
	return err
}

// SetConfiguration sets the configuration (part of homescreen layout) for the user corresponding to the provided ID.
// No GetConfiguration is provided because a User object includes Configuration already.
func (mb *MediaBrowser) SetConfiguration(userID string, configuration Configuration) error {
	url := fmt.Sprintf("%s/Users/%s/Configuration", mb.Server, userID)
	DeNullConfiguration(&configuration)
	data, status, err := mb.post(url, configuration, true)
	if customErr := mb.genericErr(status, data); customErr != nil {
		err = customErr
	}
	return err
}

// GetDisplayPreferences gets the displayPreferences (part of homescreen layout) for the user corresponding to the provided ID.
func (mb *MediaBrowser) GetDisplayPreferences(userID string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/DisplayPreferences/usersettings?userId=%s&client=emby", mb.Server, userID)
	data, status, err := mb.get(url, nil)
	if customErr := mb.genericErr(status, data); customErr != nil {
		err = customErr
	}
	if err != nil {
		return nil, err
	}
	var displayprefs map[string]interface{}
	err = json.Unmarshal([]byte(data), &displayprefs)
	if err != nil {
		return nil, err
	}
	return displayprefs, nil
}

// SetDisplayPreferences sets the displayPreferences (part of homescreen layout) for the user corresponding to the provided ID.
func (mb *MediaBrowser) SetDisplayPreferences(userID string, displayprefs map[string]interface{}) error {
	url := fmt.Sprintf("%s/DisplayPreferences/usersettings?userId=%s&client=emby", mb.Server, userID)
	data, status, err := mb.post(url, displayprefs, true)
	if customErr := mb.genericErr(status, data); customErr != nil {
		err = customErr
	}
	if err != nil {
		return err
	}
	return nil
}

// SetPassword sets the password for a user given a userID, the old password, and the new one. Requires admin authentication or authentication as the target user.
func (mb *MediaBrowser) SetPassword(userID, currentPw, newPw string) error {
	url := fmt.Sprintf("%s/Users/%s/Password", mb.Server, userID)
	data, status, err := mb.post(url, setPasswordRequest{
		Current:       currentPw,
		CurrentPw:     currentPw,
		New:           newPw,
		ResetPassword: false,
	}, true)
	if customErr := mb.genericErr(status, data); customErr != nil {
		err = customErr
	}
	return err
}

// ResetPasswordAdmin resets the given user ID's password, allowing one to then change it without knowing the previous password.
func (mb *MediaBrowser) ResetPasswordAdmin(userID string) error {
	url := fmt.Sprintf("%s/Users/%s/Password", mb.Server, userID)
	data, status, err := mb.post(url, map[string]bool{
		"ResetPassword": true,
	}, true)
	if customErr := mb.genericErr(status, data); customErr != nil {
		err = customErr
	}
	return err
}
