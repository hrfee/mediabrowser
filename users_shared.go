package mediabrowser

// Shared functions that work the same on Jellyfin & Emby.

import (
	"encoding/json"
	"fmt"
	"time"
)

// UserByName returns the user corresponding to the provided username.
func (mb *MediaBrowser) UserByName(username string, public bool) (User, int, error) {
	var match User
	find := func() (User, int, error) {
		users, status, err := mb.GetUsers(public)
		if err != nil || status != 200 {
			return User{}, status, err
		}
		for _, user := range users {
			if user.Name == username {
				return user, status, err
			}
		}
		return User{}, status, ErrUserNotFound{user: username}
	}
	match, status, err := find()
	if match.Name == "" {
		mb.CacheExpiry = time.Now()
		match, status, err = find()
	}
	return match, status, err
}

// SetPolicy sets the access policy for the user corresponding to the provided ID.
// No GetPolicy is provided because a User object includes Policy already.
func (mb *MediaBrowser) SetPolicy(userID string, policy Policy) (int, error) {
	url := fmt.Sprintf("%s/Users/%s/Policy", mb.Server, userID)
	data, status, err := mb.post(url, policy, true)
	if status == 400 {
		err = ErrNoPolicySupplied{}
		if mb.Verbose {
			json.Unmarshal([]byte(data), &err)
		}
	} else if customErr := mb.genericErr(status, data); customErr != nil {
		err = customErr
	}
	return status, err
}

// SetConfiguration sets the configuration (part of homescreen layout) for the user corresponding to the provided ID.
// No GetConfiguration is provided because a User object includes Configuration already.
func (mb *MediaBrowser) SetConfiguration(userID string, configuration Configuration) (int, error) {
	url := fmt.Sprintf("%s/Users/%s/Configuration", mb.Server, userID)
	data, status, err := mb.post(url, configuration, true)
	if customErr := mb.genericErr(status, data); customErr != nil {
		err = customErr
	}
	return status, err
}

// GetDisplayPreferences gets the displayPreferences (part of homescreen layout) for the user corresponding to the provided ID.
func (mb *MediaBrowser) GetDisplayPreferences(userID string) (map[string]interface{}, int, error) {
	url := fmt.Sprintf("%s/DisplayPreferences/usersettings?userId=%s&client=emby", mb.Server, userID)
	data, status, err := mb.get(url, nil)
	if customErr := mb.genericErr(status, data); customErr != nil {
		err = customErr
	}
	if err != nil || !(status == 204 || status == 200) {
		return nil, status, err
	}
	var displayprefs map[string]interface{}
	err = json.Unmarshal([]byte(data), &displayprefs)
	if err != nil {
		return nil, status, err
	}
	return displayprefs, status, nil
}

// SetDisplayPreferences sets the displayPreferences (part of homescreen layout) for the user corresponding to the provided ID.
func (mb *MediaBrowser) SetDisplayPreferences(userID string, displayprefs map[string]interface{}) (int, error) {
	url := fmt.Sprintf("%s/DisplayPreferences/usersettings?userId=%s&client=emby", mb.Server, userID)
	data, status, err := mb.post(url, displayprefs, true)
	if customErr := mb.genericErr(status, data); customErr != nil {
		err = customErr
	}
	if err != nil || !(status == 204 || status == 200) {
		return status, err
	}
	return status, nil
}

func (mb *MediaBrowser) SetPassword(userID, currentPw, newPw string) (int, error) {
	url := fmt.Sprintf("%s/Users/%s/Password", mb.Server, userID)
	_, status, err := mb.post(url, setPasswordRequest{
		Current:       currentPw,
		CurrentPw:     currentPw,
		New:           newPw,
		ResetPassword: false,
	}, false)
	return status, err
}
