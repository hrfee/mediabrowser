package mediabrowser

// Shared functions that work the same on Jellyfin & Emby.

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// GetUsers returns all (visible) users on the instance. If public, no authentication is needed but hidden users will not be visible.
func (mb *MediaBrowser) GetUsers(public bool) ([]User, error) {
	if !public && !mb.Authenticated {
		_, err := mb.Authenticate(mb.Username, mb.password)
		if err != nil {
			return []User{}, err
		}
	}
	if err := mb.syncUserCache(public); err != nil {
		return nil, err
	}
	return mb.userCache, nil
}

func (mb *MediaBrowser) syncUserCache(public bool) error {
	if mb.CacheExpiry.After(time.Now()) {
		return nil
	}

	syncStatus := make(chan error)

	go func(syncStatus chan error, mb *MediaBrowser) {
		mb.syncLock.Lock()
		alreadySyncing := mb.syncing
		// We're either already syncing or will be
		mb.syncing = true
		mb.syncLock.Unlock()
		if !alreadySyncing {
			var result []User
			var data string
			var status int
			var err error

			if public {
				url := fmt.Sprintf("%s/users/public", mb.Server)
				data, status, err = mb.get(url, nil)
			} else {
				url := fmt.Sprintf("%s/users", mb.Server)
				data, status, err = mb.get(url, mb.loginParams)
			}
			if customErr := mb.genericErr(status, data); customErr != nil {
				err = customErr
			}
			if err != nil {
				mb.syncLock.Lock()
				mb.syncing = false
				mb.syncLock.Unlock()
				syncStatus <- err
				return
			}
			err = json.Unmarshal([]byte(data), &result)
			if err != nil {
				mb.syncLock.Lock()
				mb.syncing = false
				mb.syncLock.Unlock()
				syncStatus <- err
				return
			}
			mb.userCache = result
			mb.usersByID = map[string]int{}
			mb.usersByName = map[string]int{}
			for i := range mb.userCache {
				mb.usersByID[mb.userCache[i].ID] = i
				// While usernames have case, Jellyfin (at least) counts identical usernames with different cases as identical.
				mb.usersByName[strings.ToLower(mb.userCache[i].Name)] = i
			}
			mb.CacheExpiry = time.Now().Add(time.Minute * time.Duration(mb.cacheLength))
			// Quirk
			if len(result) != 0 && result[0].ID[8] == '-' {
				mb.Hyphens = true
			}

			mb.syncLock.Lock()
			mb.syncing = false
			mb.syncLock.Unlock()
		} else {
			for mb.syncing {
				continue
			}
		}
		syncStatus <- nil
	}(syncStatus, mb)

	// Wait for completion
	err := <-syncStatus
	return err
}

// UserByID returns the user corresponding to the provided ID.
func (mb *MediaBrowser) UserByID(userID string, public bool) (User, error) {
	if userID == "" {
		return User{}, ErrUserNotFound{}
	}
	if u, err := mb.UserByIDFromCache(userID); err == nil {
		return u, err
	}
	// If the user isn't found in the cache then we update it
	if !mb.Authenticated {
		_, err := mb.Authenticate(mb.Username, mb.password)
		if err != nil {
			return User{}, err
		}
	}

	if public {
		_, err := mb.GetUsers(public)
		if err != nil {
			return User{}, err
		}
		if i, ok := mb.usersByID[userID]; ok {
			return mb.userCache[i], nil
		}
		return User{}, ErrUserNotFound{id: userID}
	}
	var result User
	var data string
	var status int
	var err error
	url := fmt.Sprintf("%s/users/%s", mb.Server, userID)
	data, status, err = mb.get(url, mb.loginParams)
	if (status == 404 && (mb.serverType == EmbyServer || data == "User not found")) || status == 400 {
		// 400 is really an "invalid ID", but we'll keep it as this for now.
		newErr := ErrUserNotFound{id: userID}
		if mb.Verbose {
			newErr.Verbose = true
			newErr.Code = status
			newErr.Data = data
		}
		err = newErr
	} else if customErr := mb.genericErr(status, data); customErr != nil {
		err = customErr
	}
	if err != nil {
		return User{}, err
	}
	json.Unmarshal([]byte(data), &result)
	return result, nil
}

// UserByIDFromCache searches only the local cache (reloading it if dated) for the user, rather than falling back to Jellyfin/Emby.
func (mb *MediaBrowser) UserByIDFromCache(userID string) (User, error) {
	if userID == "" {
		return User{}, ErrUserNotFound{}
	}
	if err := mb.syncUserCache(false); err != nil {
		return User{}, err
	}
	if i, ok := mb.usersByID[userID]; ok {
		return mb.userCache[i], nil
	}
	return User{}, ErrUserNotFound{id: userID}
}

// UserByName returns the user corresponding to the provided username.
func (mb *MediaBrowser) UserByName(username string, public bool) (User, error) {
	if username == "" {
		return User{}, ErrUserNotFound{}
	}
	if err := mb.syncUserCache(false); err != nil {
		return User{}, err
	}
	username = strings.ToLower(username)
	find := func() int {
		if i, ok := mb.usersByName[username]; ok {
			return i
		}
		return -1
	}

	idx := find()
	if idx == -1 {
		// Force-reload cache if not found
		mb.CacheExpiry = time.Now()
		idx := find()
		if idx == -1 {
			return User{}, ErrUserNotFound{user: username}
		}
	}
	return mb.userCache[idx], nil
}

// UserByNameFromCache searches only the local cache (reloading it if dated) for the user, rather than falling back to Jellyfin/Emby.
func (mb *MediaBrowser) UserByNameFromCache(username string) (User, error) {
	if username == "" {
		return User{}, ErrUserNotFound{}
	}
	if err := mb.syncUserCache(false); err != nil {
		return User{}, err
	}
	if i, ok := mb.usersByName[username]; ok {
		return mb.userCache[i], nil
	}
	return User{}, ErrUserNotFound{user: username}
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
