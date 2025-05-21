package mediabrowser

// Almost identical to jfapi, with the most notable change being the password workaround.

import (
	"encoding/json"
	"fmt"
	"net/http"
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
