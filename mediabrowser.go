// Package mediabrowser provides user-related bindings to the Jellyfin & Emby APIs.
// Some data aren't bound to structs as jfa-go doesn't need to interact with them, for example DisplayPreferences.
// See Jellyfin/Emby swagger docs for more info on them.
package mediabrowser

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// TimeoutHandler should recover from an http timeout or panic.
type TimeoutHandler func()

// NewNamedTimeoutHandler returns a new Timeout handler that logs the error.
// name is the name of the server to use in the log (e.g Jellyfin/Emby)
// addr is the address of the server being accessed
// if noFail is false, the program will exit on a timeout.
func NewNamedTimeoutHandler(name, addr string, noFail bool) TimeoutHandler {
	return func() {
		if r := recover(); r != nil {
			out := fmt.Sprintf("Failed to authenticate with %s @ %s: %v", name, addr, r)
			if noFail {
				log.Print(out)
			} else {
				log.Fatalf(out)
			}
		}
	}
}

type serverType int

const (
	JellyfinServer serverType = iota
	EmbyServer
)

// ServerInfo stores info about the server.
type ServerInfo struct {
	LocalAddress string `json:"LocalAddress"`
	Name         string `json:"ServerName"`
	Version      string `json:"Version"`
	OS           string `json:"OperatingSystem"`
	ID           string `json:"Id"`
}

// MediaBrowser is an api instance of Jellyfin/Emby.
type MediaBrowser struct {
	Server                          string
	client                          string
	version                         string
	device                          string
	deviceID                        string
	useragent                       string
	auth                            string
	header                          map[string]string
	ServerInfo                      ServerInfo
	Username                        string
	password                        string
	Authenticated                   bool
	AccessToken                     string
	userID                          string
	httpClient                      *http.Client
	loginParams                     map[string]string
	userCache                       []User
	libraryCache                    []VirtualFolder
	CacheExpiry, LibraryCacheExpiry time.Time // first is UserCacheExpiry, keeping name for compatability
	cacheLength                     int
	noFail                          bool
	Hyphens                         bool
	serverType                      serverType
	timeoutHandler                  TimeoutHandler
	Verbose                         bool // Jellyfin only, errors will include more info when true
}

// NewServer returns a new Mediabrowser object.
func NewServer(st serverType, server, client, version, device, deviceID string, timeoutHandler TimeoutHandler, cacheTimeout int) (*MediaBrowser, error) {
	mb := &MediaBrowser{}
	mb.serverType = st
	mb.Server = server
	mb.client = client
	mb.version = version
	mb.device = device
	mb.deviceID = deviceID
	mb.useragent = fmt.Sprintf("%s/%s", client, version)
	mb.timeoutHandler = timeoutHandler
	mb.auth = fmt.Sprintf("MediaBrowser Client=\"%s\", Device=\"%s\", DeviceId=\"%s\", Version=\"%s\"", client, device, deviceID, version)
	mb.header = map[string]string{
		"Accept":               "application/json",
		"Content-type":         "application/json; charset=UTF-8",
		"X-Application":        mb.useragent,
		"Accept-Charset":       "UTF-8,*",
		"Accept-Encoding":      "gzip",
		"User-Agent":           mb.useragent,
		"X-Emby-Authorization": mb.auth,
	}
	mb.httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}
	infoURL := fmt.Sprintf("%s/System/Info/Public", server)
	req, err := http.NewRequest("GET", infoURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := mb.httpClient.Do(req)
	defer mb.timeoutHandler()
	if err == nil {
		defer resp.Body.Close()
		data, _ := io.ReadAll(resp.Body)
		json.Unmarshal(data, &mb.ServerInfo)
	}
	mb.cacheLength = cacheTimeout
	mb.CacheExpiry, mb.LibraryCacheExpiry = time.Now(), time.Now()
	return mb, nil
}

// SetTransport sets the HTTP transport to be used for all requests. Can be used to set a proxy.
func (mb *MediaBrowser) SetTransport(t *http.Transport) {
	mb.httpClient.Transport = t
}

func bodyToString(resp *http.Response) string {
	var data io.Reader
	encoding := resp.Header.Get("Content-Encoding")
	switch encoding {
	case "gzip":
		data, _ = gzip.NewReader(resp.Body)
	default:
		data = resp.Body
	}
	buf := new(strings.Builder)
	io.Copy(buf, data)
	return buf.String()
}

func (mb *MediaBrowser) get(url string, params map[string]string) (string, int, error) {
	var req *http.Request
	if params != nil {
		jsonParams, _ := json.Marshal(params)
		req, _ = http.NewRequest("GET", url, bytes.NewBuffer(jsonParams))
	} else {
		req, _ = http.NewRequest("GET", url, nil)
	}
	for name, value := range mb.header {
		req.Header.Add(name, value)
	}
	resp, err := mb.httpClient.Do(req)
	defer mb.timeoutHandler()
	defer resp.Body.Close()
	if err != nil || resp.StatusCode != 200 {
		if resp.StatusCode == 401 && mb.Authenticated {
			mb.Authenticated = false
			_, authErr := mb.Authenticate(mb.Username, mb.password)
			if authErr == nil {
				v1, v2, v3 := mb.get(url, params)
				return v1, v2, v3
			}
		}
		return "", resp.StatusCode, err
	}
	//var respData map[string]interface{}
	//json.NewDecoder(data).Decode(&respData)
	return bodyToString(resp), resp.StatusCode, nil
}

func (mb *MediaBrowser) post(url string, data interface{}, response bool) (string, int, error) {
	params, _ := json.Marshal(data)
	// fmt.Printf("Data: %s\n", string(params))
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(params))
	for name, value := range mb.header {
		req.Header.Add(name, value)
	}
	resp, err := mb.httpClient.Do(req)
	defer mb.timeoutHandler()
	defer resp.Body.Close()
	if err != nil || resp.StatusCode != 200 {
		if resp.StatusCode == 401 && mb.Authenticated {
			mb.Authenticated = false
			_, authErr := mb.Authenticate(mb.Username, mb.password)
			if authErr == nil {
				v1, v2, v3 := mb.post(url, data, response)
				return v1, v2, v3
			}
		}
		return "", resp.StatusCode, err
	}
	if response {
		defer resp.Body.Close()
		return bodyToString(resp), resp.StatusCode, nil
	}
	return "", resp.StatusCode, nil
}

// Authenticate attempts to authenticate using a username & password
func (mb *MediaBrowser) Authenticate(username, password string) (User, error) {
	mb.Username = username
	mb.password = password
	mb.loginParams = map[string]string{
		"Username": username,
		"Pw":       password,
		"Password": password,
	}
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(mb.loginParams)
	if err != nil {
		return User{}, err
	}
	// loginParams, _ := json.Marshal(jf.loginParams)
	url := fmt.Sprintf("%s/Users/authenticatebyname", mb.Server)
	req, err := http.NewRequest("POST", url, buffer)
	defer mb.timeoutHandler()
	if err != nil {
		return User{}, err
	}
	for name, value := range mb.header {
		req.Header.Add(name, value)
	}
	resp, err := mb.httpClient.Do(req)
	if err != nil {
		return User{}, err
	}
	// Jellyfin likes to return 400 for a lot of things, even if the api docs don't say so.
	if resp.StatusCode == 400 {
		err = ErrUnauthorized{}
	} else if customErr := mb.genericErr(resp.StatusCode, ""); customErr != nil {
		err = customErr
	}
	if err != nil {
		return User{}, err
	}
	defer resp.Body.Close()
	var d io.Reader
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		d, _ = gzip.NewReader(resp.Body)
	default:
		d = resp.Body
	}
	data, err := io.ReadAll(d)
	if err != nil {
		return User{}, err
	}
	var respData map[string]interface{}
	json.Unmarshal(data, &respData)
	mb.AccessToken = respData["AccessToken"].(string)
	var user User
	ju, err := json.Marshal(respData["User"])
	if err != nil {
		return User{}, err
	}
	json.Unmarshal(ju, &user)
	mb.userID = user.ID
	mb.auth = fmt.Sprintf("MediaBrowser Client=\"%s\", Device=\"%s\", DeviceId=\"%s\", Version=\"%s\", Token=\"%s\"", mb.client, mb.device, mb.deviceID, mb.version, mb.AccessToken)
	mb.header["X-Emby-Authorization"] = mb.auth
	mb.Authenticated = true
	return user, nil
}

// MustAuthenticateOptions is used to control the behaviour of the MustAuthenticate method.
type MustAuthenticateOptions struct {
	RetryCount  int           // Number of Retries before failure.
	RetryGap    time.Duration // Duration to wait between tries.
	LogFailures bool          // Whether or not to print failures to the log.
}

// MustAuthenticate attempts to authenticate using a username & password, with configurable retries in the event of failure.
func (mb *MediaBrowser) MustAuthenticate(username, password string, opts MustAuthenticateOptions) (user User, err error) {
	for i := 0; i < opts.RetryCount; i++ {
		user, err = mb.Authenticate(username, password)
		if err == nil {
			return
		}
		if opts.LogFailures {
			log.Printf("Failed to authenticate on attempt %d, retrying in %s...\n", i+1, opts.RetryGap)
		}
		time.Sleep(opts.RetryGap)
	}
	return
}

// DeleteUser deletes the user corresponding to the provided ID.
func (mb *MediaBrowser) DeleteUser(userID string) error {
	if mb.serverType == JellyfinServer {
		return jfDeleteUser(mb, userID)
	}
	return embyDeleteUser(mb, userID)
}

// GetUsers returns all (visible) users on the instance. If public, no authentication is needed but hidden users will not be visible.
func (mb *MediaBrowser) GetUsers(public bool) ([]User, error) {
	if mb.serverType == JellyfinServer {
		return jfGetUsers(mb, public)
	}
	return embyGetUsers(mb, public)
}

// UserByID returns the user corresponding to the provided ID.
func (mb *MediaBrowser) UserByID(userID string, public bool) (User, error) {
	if mb.serverType == JellyfinServer {
		return jfUserByID(mb, userID, public)
	}
	return embyUserByID(mb, userID, public)
}

// NewUser creates a new user with the provided username and password.
func (mb *MediaBrowser) NewUser(username, password string) (User, error) {
	if mb.serverType == JellyfinServer {
		return jfNewUser(mb, username, password)
	}
	return embyNewUser(mb, username, password)
}

// ResetPassword resets a user's password by setting it to the given PIN,
// which is generated when a user attempts to reset on the login page.
// Only supported on Jellyfin, will return (PasswordResetResponse, -1, nil) on Emby.
func (mb *MediaBrowser) ResetPassword(pin string) (PasswordResetResponse, error) {
	if mb.serverType == EmbyServer {
		return PasswordResetResponse{}, nil
	}
	return jfResetPassword(mb, pin)
}
