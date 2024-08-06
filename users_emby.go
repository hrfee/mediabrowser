package mediabrowser

// Almost identical to jfapi, with the most notable change being the password workaround.

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func embyDeleteUser(emby *MediaBrowser, userID string) error {
	url := fmt.Sprintf("%s/Users/%s", emby.Server, userID)
	req, _ := http.NewRequest("DELETE", url, nil)
	for name, value := range emby.header {
		req.Header.Add(name, value)
	}
	resp, err := emby.httpClient.Do(req)
	defer emby.timeoutHandler()
	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		err = ErrUserNotFound{id: userID}
	} else if customErr := emby.genericErr(resp.StatusCode, ""); customErr != nil {
		err = customErr
	}
	return err
}

func embyGetUsers(emby *MediaBrowser, public bool) ([]User, error) {
	var result []User
	var data string
	var status int
	var err error
	if time.Now().After(emby.CacheExpiry) {
		if public {
			url := fmt.Sprintf("%s/users/public", emby.Server)
			data, status, err = emby.get(url, nil)
		} else {
			url := fmt.Sprintf("%s/users", emby.Server)
			data, status, err = emby.get(url, emby.loginParams)
		}
		if customErr := emby.genericErr(status, ""); customErr != nil {
			err = customErr
		}
		if err != nil {
			return nil, err
		}
		err := json.Unmarshal([]byte(data), &result)
		if err != nil {
			return nil, err
		}
		emby.userCache = result
		emby.CacheExpiry = time.Now().Add(time.Minute * time.Duration(emby.cacheLength))
		if result[0].ID[8] == '-' {
			emby.Hyphens = true
		}
		return result, nil
	}
	return emby.userCache, nil
}

func embyUserByID(emby *MediaBrowser, userID string, public bool) (User, error) {
	if emby.CacheExpiry.After(time.Now()) {
		for _, user := range emby.userCache {
			if user.ID == userID {
				return user, nil
			}
		}
		// If the user isn't found in the cache then we update it
	}
	if public {
		users, err := emby.GetUsers(public)
		if err != nil {
			return User{}, err
		}
		for _, user := range users {
			if user.ID == userID {
				return user, nil
			}
		}
		return User{}, ErrUserNotFound{id: userID}
	}
	var result User
	var data string
	var status int
	var err error
	url := fmt.Sprintf("%s/users/%s", emby.Server, userID)
	data, status, err = emby.get(url, emby.loginParams)
	if status == 404 || status == 400 {
		err = ErrUserNotFound{id: userID}
	} else if customErr := emby.genericErr(status, ""); customErr != nil {
		err = customErr
	}
	if err != nil {
		return User{}, err
	}
	json.Unmarshal([]byte(data), &result)
	return result, nil
}

// Since emby doesn't allow one to specify a password on user creation, we:
// Create the account
// Immediately disable it
// Set password
// Re-enable it
func embyNewUser(emby *MediaBrowser, username, password string) (User, error) {
	url := fmt.Sprintf("%s/Users/New", emby.Server)
	data := map[string]interface{}{
		"Name": username,
	}
	response, status, err := emby.post(url, data, true)
	if customErr := emby.genericErr(status, ""); customErr != nil {
		err = customErr
	}
	if err != nil {
		return User{}, err
	}
	var recv User
	json.Unmarshal([]byte(response), &recv)
	// Step 2: Set password
	id := recv.ID
	url = fmt.Sprintf("%s/Users/%s/Password", emby.Server, id)
	data = map[string]interface{}{
		"Id":        id,
		"CurrentPw": "",
		"NewPw":     password,
	}
	var resp string
	resp, status, err = emby.post(url, data, true)
	if customErr := emby.genericErr(status, resp); customErr != nil {
		err = customErr
	}
	// Step 3: If setting password errored, try to delete the account
	if err != nil {
		err = emby.DeleteUser(id)
	}
	return recv, err
}
