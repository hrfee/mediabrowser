package mediabrowser

import (
	"encoding/json"
	"fmt"
	"net/http"
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
