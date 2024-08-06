package mediabrowser

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func jfDeleteUser(jf *MediaBrowser, userID string) error {
	url := fmt.Sprintf("%s/Users/%s", jf.Server, userID)
	req, _ := http.NewRequest("DELETE", url, nil)
	for name, value := range jf.header {
		req.Header.Add(name, value)
	}
	resp, err := jf.httpClient.Do(req)
	defer jf.timeoutHandler()
	defer resp.Body.Close()
	data := ""
	if jf.Verbose {
		data = bodyToString(resp)
	}
	// Should be 404 but sometimes isn't
	if resp.StatusCode == 404 || resp.StatusCode == 500 {
		err = ErrUserNotFound{id: userID}
		if jf.Verbose {
			json.Unmarshal([]byte(data), &err)
		}
	} else if customErr := jf.genericErr(resp.StatusCode, data); customErr != nil {
		err = customErr
	}
	return err
}

func jfGetUsers(jf *MediaBrowser, public bool) ([]User, error) {
	var result []User
	var data string
	var status int
	var err error
	if time.Now().After(jf.CacheExpiry) {
		if public {
			url := fmt.Sprintf("%s/users/public", jf.Server)
			data, status, err = jf.get(url, nil)
		} else {
			url := fmt.Sprintf("%s/users", jf.Server)
			data, status, err = jf.get(url, jf.loginParams)
		}
		if customErr := jf.genericErr(status, data); customErr != nil {
			err = customErr
		}
		if err != nil {
			return nil, err
		}
		err := json.Unmarshal([]byte(data), &result)
		if err != nil {
			return nil, err
		}
		jf.userCache = result
		jf.CacheExpiry = time.Now().Add(time.Minute * time.Duration(jf.cacheLength))
		if result[0].ID[8] == '-' {
			jf.Hyphens = true
		}
		return result, err
	}
	return jf.userCache, nil
}

func jfUserByID(jf *MediaBrowser, userID string, public bool) (User, error) {
	if jf.CacheExpiry.After(time.Now()) {
		for _, user := range jf.userCache {
			if user.ID == userID {
				return user, nil
			}
		}
		// If the user isn't found in the cache then we update it
	}
	if public {
		users, err := jf.GetUsers(public)
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
	url := fmt.Sprintf("%s/users/%s", jf.Server, userID)
	data, status, err = jf.get(url, jf.loginParams)
	if status == 404 || status == 400 {
		err = ErrUserNotFound{id: userID}
		if jf.Verbose {
			json.Unmarshal([]byte(data), &err)
		}
	} else if customErr := jf.genericErr(status, data); customErr != nil {
		err = customErr
	}
	if err != nil {
		return User{}, err
	}
	json.Unmarshal([]byte(data), &result)
	return result, nil
}

func jfNewUser(jf *MediaBrowser, username, password string) (User, error) {
	url := fmt.Sprintf("%s/Users/New", jf.Server)
	stringData := map[string]string{
		"Name":     username,
		"Password": password,
	}
	data := map[string]interface{}{}
	for key, value := range stringData {
		data[key] = value
	}
	resp, status, err := jf.post(url, data, true)
	if customErr := jf.genericErr(status, resp); customErr != nil {
		err = customErr
	}
	if err != nil {
		return User{}, err
	}
	var recv User
	json.Unmarshal([]byte(resp), &recv)
	return recv, nil
}

func jfResetPassword(jf *MediaBrowser, pin string) (PasswordResetResponse, error) {
	url := fmt.Sprintf("%s/Users/ForgotPassword/Pin", jf.Server)
	resp, status, err := jf.post(url, map[string]string{
		"Pin": pin,
	}, true)
	if customErr := jf.genericErr(status, resp); customErr != nil {
		err = customErr
	}
	recv := PasswordResetResponse{}
	if err != nil {
		return recv, err
	}
	json.Unmarshal([]byte(resp), &recv)
	return recv, err
}
