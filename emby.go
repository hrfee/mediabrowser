package mediabrowser

// Almost identical to jfapi, with the most notable change being the password workaround.

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func embyDeleteUser(emby *MediaBrowser, userID string) (int, error) {
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
	return resp.StatusCode, err
}

func embyGetUsers(emby *MediaBrowser, public bool) ([]User, int, error) {
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
		if err != nil || status != 200 {
			return nil, status, err
		}
		err := json.Unmarshal([]byte(data), &result)
		if err != nil {
			return nil, status, err
		}
		emby.userCache = result
		emby.CacheExpiry = time.Now().Add(time.Minute * time.Duration(emby.cacheLength))
		if result[0].ID[8] == '-' {
			emby.Hyphens = true
		}
		return result, status, nil
	}
	return emby.userCache, 200, nil
}

func embyUserByName(emby *MediaBrowser, username string, public bool) (User, int, error) {
	var match User
	find := func() (User, int, error) {
		users, status, err := emby.GetUsers(public)
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
		emby.CacheExpiry = time.Now()
		match, status, err = find()
	}
	return match, status, err
}

func embyUserByID(emby *MediaBrowser, userID string, public bool) (User, int, error) {
	if emby.CacheExpiry.After(time.Now()) {
		for _, user := range emby.userCache {
			if user.ID == userID {
				return user, 200, nil
			}
		}
		// If the user isn't found in the cache then we update it
	}
	if public {
		users, status, err := emby.GetUsers(public)
		if err != nil || status != 200 {
			return User{}, status, err
		}
		for _, user := range users {
			if user.ID == userID {
				return user, status, nil
			}
		}
		return User{}, status, ErrUserNotFound{id: userID}
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
	if err != nil || status != 200 {
		return User{}, status, err
	}
	json.Unmarshal([]byte(data), &result)
	return result, status, nil
}

// Since emby doesn't allow one to specify a password on user creation, we:
// Create the account
// Immediately disable it
// Set password
// Re-enable it
func embyNewUser(emby *MediaBrowser, username, password string) (User, int, error) {
	url := fmt.Sprintf("%s/Users/New", emby.Server)
	data := map[string]interface{}{
		"Name": username,
	}
	response, status, err := emby.post(url, data, true)
	if customErr := emby.genericErr(status, ""); customErr != nil {
		err = customErr
	}
	if err != nil || !(status == 200 || status == 204) {
		return User{}, status, err
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
	_, status, err = emby.post(url, data, false)
	// Step 3: If setting password errored, try to delete the account
	if err != nil || !(status == 200 || status == 204) {
		_, err = emby.DeleteUser(id)
	}
	return recv, status, err
}

func embySetPolicy(emby *MediaBrowser, userID string, policy Policy) (int, error) {
	url := fmt.Sprintf("%s/Users/%s/Policy", emby.Server, userID)
	_, status, err := emby.post(url, policy, false)
	if status == 400 {
		err = ErrNoPolicySupplied{}
	} else if customErr := emby.genericErr(status, ""); customErr != nil {
		err = customErr
	}
	if err != nil || status != 200 {
		return status, err
	}
	return status, nil
}

func embySetConfiguration(emby *MediaBrowser, userID string, configuration Configuration) (int, error) {
	url := fmt.Sprintf("%s/Users/%s/Configuration", emby.Server, userID)
	_, status, err := emby.post(url, configuration, false)
	if customErr := emby.genericErr(status, ""); customErr != nil {
		err = customErr
	}
	return status, err
}

func embyGetDisplayPreferences(emby *MediaBrowser, userID string) (map[string]interface{}, int, error) {
	url := fmt.Sprintf("%s/DisplayPreferences/usersettings?userId=%s&client=emby", emby.Server, userID)
	data, status, err := emby.get(url, nil)
	if customErr := emby.genericErr(status, ""); customErr != nil {
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

func embySetDisplayPreferences(emby *MediaBrowser, userID string, displayprefs map[string]interface{}) (int, error) {
	url := fmt.Sprintf("%s/DisplayPreferences/usersettings?userId=%s&client=emby", emby.Server, userID)
	_, status, err := emby.post(url, displayprefs, false)
	if customErr := emby.genericErr(status, ""); customErr != nil {
		err = customErr
	}
	if err != nil || !(status == 204 || status == 200) {
		return status, err
	}
	return status, nil
}

func embyGetLibraries(emby *MediaBrowser) ([]VirtualFolder, int, error) {
	var result []VirtualFolder
	var data string
	var status int
	var err error
	if time.Now().After(emby.LibraryCacheExpiry) {
		url := fmt.Sprintf("%s/Library/VirtualFolders", emby.Server)
		data, status, err = emby.get(url, nil)
		if customErr := emby.genericErr(status, ""); customErr != nil {
			err = customErr
		}
		if err != nil || status != 200 {
			return nil, status, err
		}
		err := json.Unmarshal([]byte(data), &result)
		if err != nil {
			return nil, status, err
		}
		emby.libraryCache = result
		emby.LibraryCacheExpiry = time.Now().Add(time.Minute * time.Duration(emby.cacheLength))
		return result, status, nil
	}
	return emby.libraryCache, 200, nil
}

func embyAddLibrary(emby *MediaBrowser, name string, collectionType string, paths []string, refreshLibrary bool, LibraryOptions LibraryOptions) (int, error) {
	pathQuery := ""
	for _, path := range paths {
		// element is the element from someSlice for where we are
		pathQuery += fmt.Sprintf("&paths[]=%s", path)
	}
	url := fmt.Sprintf("%s/Library/VirtualFolders?client=emby&name=%s&collectiontype=%s&refreshLibrary=%t%s", emby.Server, name, collectionType, refreshLibrary, pathQuery)
	_, status, err := emby.post(url, LibraryOptions, false)
	if customErr := emby.genericErr(status, ""); customErr != nil {
		err = customErr
	}
	return status, err
}

func embyDeleteLibrary(emby *MediaBrowser, name string) (int, error) {
	url := fmt.Sprintf("%s/Library/VirtualFolders?name=%s", emby.Server, name)
	req, _ := http.NewRequest("DELETE", url, nil)
	for name, value := range emby.header {
		req.Header.Add(name, value)
	}
	resp, err := emby.httpClient.Do(req)
	defer emby.timeoutHandler()
	defer resp.Body.Close()
	if customErr := emby.genericErr(resp.StatusCode, ""); customErr != nil {
		err = customErr
	}
	return resp.StatusCode, err
}

func embyAddFolder(emby *MediaBrowser, refreshLibrary bool, AddMedia AddMedia) (int, error) {
	url := fmt.Sprintf("%s/Library/VirtualFolders/Paths?client=emby&refreshLibrary=%t", emby.Server, refreshLibrary)
	_, status, err := emby.post(url, AddMedia, false)
	if customErr := emby.genericErr(status, ""); customErr != nil {
		err = customErr
	}
	return status, err
}

func embyDeleteFolder(emby *MediaBrowser, name string, path string, refreshLibrary bool) (int, error) {
	url := fmt.Sprintf("%s/Library/VirtualFolders/Paths?name=%s&path=%s&refreshLibrary=%t", emby.Server, name, path, refreshLibrary)
	req, _ := http.NewRequest("DELETE", url, nil)
	for name, value := range emby.header {
		req.Header.Add(name, value)
	}
	resp, err := emby.httpClient.Do(req)
	defer emby.timeoutHandler()
	defer resp.Body.Close()
	if customErr := emby.genericErr(resp.StatusCode, ""); customErr != nil {
		err = customErr
	}
	return resp.StatusCode, err
}

func embyScanLibs(emby *MediaBrowser) (int, error) {
	url := fmt.Sprintf("%s/Library/Refresh?client=emby", emby.Server)
	_, status, err := emby.post(url, nil, false)
	if customErr := emby.genericErr(status, ""); customErr != nil {
		err = customErr
	}
	return status, err
}
