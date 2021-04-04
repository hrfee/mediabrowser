package mediabrowser

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// GetLibraries returns a list of the Libaries (called VirtualFolders) on this node.
func (mb *MediaBrowser) GetLibraries() ([]VirtualFolder, int, error) {
	var result []VirtualFolder
	var data string
	var status int
	var err error
	if time.Now().After(mb.LibraryCacheExpiry) {
		url := fmt.Sprintf("%s/Library/VirtualFolders", mb.Server)
		data, status, err = mb.get(url, nil)
		if customErr := mb.genericErr(status, ""); customErr != nil {
			err = customErr
		}
		if err != nil || status != 200 {
			return nil, status, err
		}
		err := json.Unmarshal([]byte(data), &result)
		if err != nil {
			return nil, status, err
		}
		mb.libraryCache = result
		mb.LibraryCacheExpiry = time.Now().Add(time.Minute * time.Duration(mb.cacheLength))
		return result, status, nil
	}
	return mb.libraryCache, 200, nil
}

// AddLibrary creates a library (VirtualFolder) for this node.
func (mb *MediaBrowser) AddLibrary(name string, collectionType string, paths []string, refreshLibrary bool, LibraryOptions LibraryOptions) (int, error) {
	pathQuery := ""
	for _, path := range paths {
		// element is the element from someSlice for where we are
		pathQuery += fmt.Sprintf("&paths[]=%s", path)
	}
	url := fmt.Sprintf("%s/Library/VirtualFolders?client=emby&name=%s&collectiontype=%s&refreshLibrary=%t%s", mb.Server, name, collectionType, refreshLibrary, pathQuery)
	_, status, err := mb.post(url, LibraryOptions, false)
	if customErr := mb.genericErr(status, ""); customErr != nil {
		err = customErr
	}
	return status, err
}

// DeleteLibrary deletes the library (VirtualFolder) corresponding to the provided name.
func (mb *MediaBrowser) DeleteLibrary(name string) (int, error) {
	url := fmt.Sprintf("%s/Library/VirtualFolders?name=%s", mb.Server, name)
	req, _ := http.NewRequest("DELETE", url, nil)
	for name, value := range mb.header {
		req.Header.Add(name, value)
	}
	resp, err := mb.httpClient.Do(req)
	defer mb.timeoutHandler()
	defer resp.Body.Close()
	if customErr := mb.genericErr(resp.StatusCode, ""); customErr != nil {
		err = customErr
	}
	return resp.StatusCode, err
}

// AddFolder adds a subfolder to a library (VirtualFolder)
func (mb *MediaBrowser) AddFolder(refreshLibrary bool, AddMedia AddMedia) (int, error) {
	url := fmt.Sprintf("%s/Library/VirtualFolders/Paths?client=emby&refreshLibrary=%t", mb.Server, refreshLibrary)
	_, status, err := mb.post(url, AddMedia, false)
	if customErr := mb.genericErr(status, ""); customErr != nil {
		err = customErr
	}
	return status, err
}

// DeleteFolder deletes the library (VirtualFolder) corresponding to the provided name.
func (mb *MediaBrowser) DeleteFolder(name string, path string, refreshLibrary bool) (int, error) {
	url := fmt.Sprintf("%s/Library/VirtualFolders/Paths?name=%s&path=%s&refreshLibrary=%t", mb.Server, name, path, refreshLibrary)
	req, _ := http.NewRequest("DELETE", url, nil)
	for name, value := range mb.header {
		req.Header.Add(name, value)
	}
	resp, err := mb.httpClient.Do(req)
	defer mb.timeoutHandler()
	defer resp.Body.Close()
	if customErr := mb.genericErr(resp.StatusCode, ""); customErr != nil {
		err = customErr
	}
	return resp.StatusCode, err
}

// ScanLibs triggers a scan of all libraries.
func (mb *MediaBrowser) ScanLibs() (int, error) {
	url := fmt.Sprintf("%s/Library/Refresh?client=emby", mb.Server)
	_, status, err := mb.post(url, nil, false)
	if customErr := mb.genericErr(status, ""); customErr != nil {
		err = customErr
	}
	return status, err
}
