# mediabrowser

a limited API client for Jellyfin/Emby, with a focus on User management as used in [jfa-go](https://github.com/hrfee/jfa-go).

```
go get github.com/hrfee/mediabrowser
```

The below example can also be found in [`example/example.go`](https://github.com/hrfee/mediabrowser/blob/main/example/example.go).

```go
// EXAMPLE: Find and return the last activity date of user "john".
package main

import (
	"log"

	"github.com/hrfee/mediabrowser"
)

const (
	address      = "https://jellyf.in:8097"
	noFail       = true
	cacheTimeout = 30 // Reload cache 30 minutes after most recent load.
	serverType   = mediabrowser.JellyfinServer

	// Authenticates with username/pass as some operations aren't possible with an API key.
	username = "accounts"
	password = "xxxxxx"

	// Below values appear in the Jellyfin/Emby dashboard.
	client   = "mediabrowser-example"
	version  = "v0.0.0"
	device   = "mb-test-device"
	deviceID = "mb-test-device-id"
)

func main() {
	timeoutHandler := mediabrowser.NewNamedTimeoutHandler("Jellyfin", address, noFail)
	mb, _ := mediabrowser.NewServer(
		serverType,
		address,
		client,
		version,
		device,
		deviceID,
		timeoutHandler,
		cacheTimeout,
	)
	_, status, err := mb.Authenticate(username, password)
	if err != nil || status != 200 {
		log.Fatalf("Failed to authenticate: Status %d Error %v", status, err)
	}

	user, status, err := mb.UserByName("john", false)
	if user.Name == "" || err != nil || status != 200 {
		log.Fatalf("Couldn't find user: Status %d Error %v", status, err)
	}

	log.Printf("User %s was last active %v", user.Name, user.LastActivityDate)
}
```
