package mediabrowser

import (
	"encoding/json"
	"fmt"
	"time"
)

// GetActivityLog returns a page of Jellyfin/Emby activity log entries.
// Pass -1 to unset integer parameters, time.Time{} to unset since.
// hasUserID makes the server return only entries which contain a User ID, may not work on Emby.
func (mb *MediaBrowser) GetActivityLog(skip, limit int, since time.Time, hasUserID bool) (ActivityLog, error) {
	result := ActivityLog{}
	if !mb.Authenticated {
		_, err := mb.Authenticate(mb.Username, mb.password)
		if err != nil {
			return result, err
		}
	}
	url := fmt.Sprintf("%s/System/ActivityLog/Entries", mb.Server)
	qp := map[string]interface{}{}
	if skip != -1 {
		qp["startIndex"] = skip
	}
	if limit != -1 {
		qp["limit"] = limit
	}
	if !since.IsZero() {
		qp["minDate"] = toMediabrowserTime(since)
	}
	if hasUserID {
		qp["hasUserId"] = true
	}
	data, status, err := mb.get(url, nil, qp)
	if customErr := mb.genericErr(status, ""); customErr != nil {
		err = customErr
	}
	if err != nil || status != 200 {
		return result, err
	}
	err = json.Unmarshal([]byte(data), &result)
	return result, err
}
