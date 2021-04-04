package mediabrowser

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func jfDeleteUser(jf *MediaBrowser, userID string) (int, error) {
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
	return resp.StatusCode, err
}

func jfGetUsers(jf *MediaBrowser, public bool) ([]User, int, error) {
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
		if err != nil || status != 200 {
			return nil, status, err
		}
		err := json.Unmarshal([]byte(data), &result)
		if err != nil {
			return nil, status, err
		}
		jf.userCache = result
		jf.CacheExpiry = time.Now().Add(time.Minute * time.Duration(jf.cacheLength))
		if result[0].ID[8] == '-' {
			jf.Hyphens = true
		}
		return result, status, err
	}
	return jf.userCache, 200, nil
}

func jfUserByName(jf *MediaBrowser, username string, public bool) (User, int, error) {
	var match User
	find := func() (User, int, error) {
		users, status, err := jf.GetUsers(public)
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
		jf.CacheExpiry = time.Now()
		match, status, err = find()
	}
	return match, status, err
}

func jfUserByID(jf *MediaBrowser, userID string, public bool) (User, int, error) {
	if jf.CacheExpiry.After(time.Now()) {
		for _, user := range jf.userCache {
			if user.ID == userID {
				return user, 200, nil
			}
		}
		// If the user isn't found in the cache then we update it
	}
	if public {
		users, status, err := jf.GetUsers(public)
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
	if err != nil || status != 200 {
		return User{}, status, err
	}
	json.Unmarshal([]byte(data), &result)
	return result, status, nil
}

func jfNewUser(jf *MediaBrowser, username, password string) (User, int, error) {
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
	if err != nil || !(status == 200 || status == 204) {
		return User{}, status, err
	}
	var recv User
	json.Unmarshal([]byte(resp), &recv)
	return recv, status, nil
}

func jfSetPolicy(jf *MediaBrowser, userID string, policy Policy) (int, error) {
	url := fmt.Sprintf("%s/Users/%s/Policy", jf.Server, userID)
	data, status, err := jf.post(url, policy, true)
	if status == 400 {
		err = ErrNoPolicySupplied{}
		if jf.Verbose {
			json.Unmarshal([]byte(data), &err)
		}
	} else if customErr := jf.genericErr(status, data); customErr != nil {
		err = customErr
	}
	return status, err
}

func jfSetConfiguration(jf *MediaBrowser, userID string, configuration Configuration) (int, error) {
	url := fmt.Sprintf("%s/Users/%s/Configuration", jf.Server, userID)
	data, status, err := jf.post(url, configuration, true)
	if customErr := jf.genericErr(status, data); customErr != nil {
		err = customErr
	}
	return status, err
}

func jfGetDisplayPreferences(jf *MediaBrowser, userID string) (map[string]interface{}, int, error) {
	url := fmt.Sprintf("%s/DisplayPreferences/usersettings?userId=%s&client=emby", jf.Server, userID)
	data, status, err := jf.get(url, nil)
	if customErr := jf.genericErr(status, data); customErr != nil {
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

func jfSetDisplayPreferences(jf *MediaBrowser, userID string, displayprefs map[string]interface{}) (int, error) {
	url := fmt.Sprintf("%s/DisplayPreferences/usersettings?userId=%s&client=emby", jf.Server, userID)
	data, status, err := jf.post(url, displayprefs, true)
	if customErr := jf.genericErr(status, data); customErr != nil {
		err = customErr
	}
	if err != nil || !(status == 204 || status == 200) {
		return status, err
	}
	return status, nil
}

func jfResetPassword(jf *MediaBrowser, pin string) (PasswordResetResponse, int, error) {
	url := fmt.Sprintf("%s/Users/ForgotPassword/Pin", jf.Server)
	resp, status, err := jf.post(url, map[string]string{
		"Pin": pin,
	}, true)
	if customErr := jf.genericErr(status, resp); customErr != nil {
		err = customErr
	}
	recv := PasswordResetResponse{}
	if err != nil || status != 200 {
		return recv, status, err
	}
	json.Unmarshal([]byte(resp), &recv)
	return recv, status, err
}

func jfGetLibraries(jf *MediaBrowser) ([]VirtualFolder, int, error) {
	var result []VirtualFolder
	var data string
	var status int
	var err error
	if time.Now().After(jf.LibraryCacheExpiry) {
		url := fmt.Sprintf("%s/Library/VirtualFolders", jf.Server)
		data, status, err = jf.get(url, nil)
		if customErr := jf.genericErr(status, ""); customErr != nil {
			err = customErr
		}
		if err != nil || status != 200 {
			return nil, status, err
		}
		err := json.Unmarshal([]byte(data), &result)
		if err != nil {
			return nil, status, err
		}
		jf.libraryCache = result
		jf.LibraryCacheExpiry = time.Now().Add(time.Minute * time.Duration(jf.cacheLength))
		return result, status, nil
	}
	return jf.libraryCache, 200, nil
}

func jfAddLibrary(jf *MediaBrowser, name string, collectionType string, paths []string, refreshLibrary bool, LibraryOptions LibraryOptions) (int, error) {
	pathQuery := ""
	for _, path := range paths {
		// element is the element from someSlice for where we are
		pathQuery += fmt.Sprintf("&paths[]=%s", path)
	}
	url := fmt.Sprintf("%s/Library/VirtualFolders?client=emby&name=%s&collectiontype=%s&refreshLibrary=%t%s", jf.Server, name, collectionType, refreshLibrary, pathQuery)
	_, status, err := jf.post(url, LibraryOptions, false)
	if customErr := jf.genericErr(status, ""); customErr != nil {
		err = customErr
	}
	return status, err
}

func jfDeleteLibrary(jf *MediaBrowser, name string) (int, error) {
	url := fmt.Sprintf("%s/Library/VirtualFolders?name=%s", jf.Server, name)
	req, _ := http.NewRequest("DELETE", url, nil)
	for name, value := range jf.header {
		req.Header.Add(name, value)
	}
	resp, err := jf.httpClient.Do(req)
	defer jf.timeoutHandler()
	defer resp.Body.Close()
	if customErr := jf.genericErr(resp.StatusCode, ""); customErr != nil {
		err = customErr
	}
	return resp.StatusCode, err
}

func jfAddFolder(jf *MediaBrowser, refreshLibrary bool, AddMedia AddMedia) (int, error) {
	url := fmt.Sprintf("%s/Library/VirtualFolders/Paths?client=emby&refreshLibrary=%t", jf.Server, refreshLibrary)
	_, status, err := jf.post(url, AddMedia, false)
	if customErr := jf.genericErr(status, ""); customErr != nil {
		err = customErr
	}
	return status, err
}

func jfDeleteFolder(jf *MediaBrowser, name string, path string, refreshLibrary bool) (int, error) {
	url := fmt.Sprintf("%s/Library/VirtualFolders/Paths?name=%s&path=%s&refreshLibrary=%t", jf.Server, name, path, refreshLibrary)
	req, _ := http.NewRequest("DELETE", url, nil)
	for name, value := range jf.header {
		req.Header.Add(name, value)
	}
	resp, err := jf.httpClient.Do(req)
	defer jf.timeoutHandler()
	defer resp.Body.Close()
	if customErr := jf.genericErr(resp.StatusCode, ""); customErr != nil {
		err = customErr
	}
	return resp.StatusCode, err
}

func jfScanLibs(jf *MediaBrowser) (int, error) {
	url := fmt.Sprintf("%s/Library/Refresh?client=emby", jf.Server)
	_, status, err := jf.post(url, nil, false)
	if customErr := jf.genericErr(status, ""); customErr != nil {
		err = customErr
	}
	return status, err
}
