package mediabrowser

type User struct {
	Name                      string        `json:"Name"`
	ServerID                  string        `json:"ServerId"`
	ID                        string        `json:"Id"`
	HasPassword               bool          `json:"HasPassword"`
	HasConfiguredPassword     bool          `json:"HasConfiguredPassword"`
	HasConfiguredEasyPassword bool          `json:"HasConfiguredEasyPassword"`
	EnableAutoLogin           bool          `json:"EnableAutoLogin"`
	LastLoginDate             Time          `json:"LastLoginDate"`
	LastActivityDate          Time          `json:"LastActivityDate"`
	Configuration             Configuration `json:"Configuration"`
	// Policy stores the user's permissions.
	Policy Policy `json:"Policy"`
}

type SessionInfo struct {
	RemoteEndpoint string `json:"RemoteEndPoint"`
	UserID         string `json:"UserId"`
}

type AuthenticationResult struct {
	User        User        `json:"User"`
	AccessToken string      `json:"AccessToken"`
	ServerID    string      `json:"ServerId"`
	SessionInfo SessionInfo `json:"SessionInfo"`
}

type Configuration struct {
	PlayDefaultAudioTrack      bool          `json:"PlayDefaultAudioTrack"`
	SubtitleLanguagePreference string        `json:"SubtitleLanguagePreference"`
	DisplayMissingEpisodes     bool          `json:"DisplayMissingEpisodes"`
	GroupedFolders             []interface{} `json:"GroupedFolders"`
	SubtitleMode               string        `json:"SubtitleMode"`
	DisplayCollectionsView     bool          `json:"DisplayCollectionsView"`
	EnableLocalPassword        bool          `json:"EnableLocalPassword"`
	OrderedViews               []interface{} `json:"OrderedViews"`
	LatestItemsExcludes        []interface{} `json:"LatestItemsExcludes"`
	MyMediaExcludes            []interface{} `json:"MyMediaExcludes"`
	HidePlayedInLatest         bool          `json:"HidePlayedInLatest"`
	RememberAudioSelections    bool          `json:"RememberAudioSelections"`
	RememberSubtitleSelections bool          `json:"RememberSubtitleSelections"`
	EnableNextEpisodeAutoPlay  bool          `json:"EnableNextEpisodeAutoPlay"`
}

// Policy stores a users permissions.
type Policy struct {
	IsAdministrator                  bool          `json:"IsAdministrator"`
	IsHidden                         bool          `json:"IsHidden"`
	IsDisabled                       bool          `json:"IsDisabled"`
	BlockedTags                      []interface{} `json:"BlockedTags"`
	EnableUserPreferenceAccess       bool          `json:"EnableUserPreferenceAccess"`
	AccessSchedules                  []interface{} `json:"AccessSchedules"`
	BlockUnratedItems                []interface{} `json:"BlockUnratedItems"`
	EnableRemoteControlOfOtherUsers  bool          `json:"EnableRemoteControlOfOtherUsers"`
	EnableSharedDeviceControl        bool          `json:"EnableSharedDeviceControl"`
	EnableRemoteAccess               bool          `json:"EnableRemoteAccess"`
	EnableLiveTvManagement           bool          `json:"EnableLiveTvManagement"`
	EnableLiveTvAccess               bool          `json:"EnableLiveTvAccess"`
	EnableMediaPlayback              bool          `json:"EnableMediaPlayback"`
	EnableAudioPlaybackTranscoding   bool          `json:"EnableAudioPlaybackTranscoding"`
	EnableVideoPlaybackTranscoding   bool          `json:"EnableVideoPlaybackTranscoding"`
	EnablePlaybackRemuxing           bool          `json:"EnablePlaybackRemuxing"`
	EnableContentDeletion            bool          `json:"EnableContentDeletion"`
	EnableContentDeletionFromFolders []interface{} `json:"EnableContentDeletionFromFolders"`
	EnableContentDownloading         bool          `json:"EnableContentDownloading"`
	EnableSyncTranscoding            bool          `json:"EnableSyncTranscoding"`
	EnableMediaConversion            bool          `json:"EnableMediaConversion"`
	EnabledDevices                   []interface{} `json:"EnabledDevices"`
	EnableAllDevices                 bool          `json:"EnableAllDevices"`
	EnabledChannels                  []interface{} `json:"EnabledChannels"`
	EnableAllChannels                bool          `json:"EnableAllChannels"`
	EnabledFolders                   []string      `json:"EnabledFolders"`
	EnableAllFolders                 bool          `json:"EnableAllFolders"`
	InvalidLoginAttemptCount         int           `json:"InvalidLoginAttemptCount"`
	EnablePublicSharing              bool          `json:"EnablePublicSharing"`
	RemoteClientBitrateLimit         int           `json:"RemoteClientBitrateLimit"`
	AuthenticationProviderID         string        `json:"AuthenticationProviderId"`
	// Jellyfin Only
	ForceRemoteSourceTranscoding bool          `json:"ForceRemoteSourceTranscoding"`
	LoginAttemptsBeforeLockout   int           `json:"LoginAttemptsBeforeLockout"`
	MaxActiveSessions            int           `json:"MaxActiveSessions"`
	BlockedMediaFolders          []interface{} `json:"BlockedMediaFolders"`
	BlockedChannels              []interface{} `json:"BlockedChannels"`
	PasswordResetProviderID      string        `json:"PasswordResetProviderId"`
	SyncPlayAccess               string        `json:"SyncPlayAccess"`
	// Emby Only
	IsHiddenRemotely           bool          `json:"IsHiddenRemotely"`
	IsTagBlockingModeInclusive bool          `json:"IsTagBlockingModeInclusive"`
	EnableSubtitleDownloading  bool          `json:"EnableSubtitleDownloading"`
	EnableSubtitleManagement   bool          `json:"EnableSubtitleManagement"`
	ExcludedSubFolders         []interface{} `json:"ExcludedSubFolders"`
	SimultaneousStreamLimit    int           `json:"SimultaneousStreamLimit"`
}

type PasswordResetResponse struct {
	Success    bool     `json:"Success"`
	UsersReset []string `json:"UsersReset"`
}

type VirtualFolder struct {
	Name               string         `json:"Name"`
	Locations          []string       `json:"Locations"`
	CollectionType     string         `json:"CollectionType"`
	LibraryOptions     LibraryOptions `json:"LibraryOptions"`
	ItemId             string         `json:"ItemId"`
	PrimaryImageItemId string         `json:"PrimaryImageItemId"`
	RefreshProgress    float64        `json:"RefreshProgress"`
	RefreshStatus      string         `json:"RefreshStatus"`
}

type LibraryOptions struct {
	EnablePhotos                            bool          `json:"EnablePhotos"`
	EnableRealtimeMonitor                   bool          `json:"EnableRealtimeMonitor"`
	EnableChapterImageExtraction            bool          `json:"EnableChapterImageExtraction"`
	ExtractChapterImagesDuringLibraryScan   bool          `json:"ExtractChapterImagesDuringLibraryScan"`
	PathInfos                               []PathInfo    `json:"PathInfos"`
	SaveLocalMetadata                       bool          `json:"SaveLocalMetadata"`
	EnableInternetProviders                 bool          `json:"EnableInternetProviders"`
	EnableAutomaticSeriesGrouping           bool          `json:"EnableAutomaticSeriesGrouping"`
	EnableEmbeddedTitles                    bool          `json:"EnableEmbeddedTitles"`
	EnableEmbeddedEpisodeInfos              bool          `json:"EnableEmbeddedEpisodeInfos"`
	AutomaticRefreshIntervalDays            int           `json:"AutomaticRefreshIntervalDays"`
	PreferredMetadataLanguage               string        `json:"PreferredMetadataLanguage"`
	MetadataCountryCode                     string        `json:"MetadataCountryCode"`
	SeasonZeroDisplayName                   string        `json:"SeasonZeroDisplayName"`
	MetadataSavers                          []string      `json:"MetadataSavers"`
	DisabledLocalMetadataReaders            []string      `json:"DisabledLocalMetadataReaders"`
	LocalMetadataReaderOrder                []string      `json:"LocalMetadataReaderOrder"`
	DisabledSubtitleFetchers                []string      `json:"DisabledSubtitleFetchers"`
	SubtitleFetcherOrder                    []string      `json:"SubtitleFetcherOrder"`
	SkipSubtitlesIfEmbeddedSubtitlesPresent bool          `json:"SkipSubtitlesIfEmbeddedSubtitlesPresent"`
	SkipSubtitlesIfAudioTrackMatches        bool          `json:"SkipSubtitlesIfAudioTrackMatches"`
	SubtitleDownloadLanguages               []string      `json:"SubtitleDownloadLanguages"`
	RequirePerfectSubtitleMatch             bool          `json:"RequirePerfectSubtitleMatch"`
	SaveSubtitlesWithMedia                  bool          `json:"SaveSubtitlesWithMedia"`
	TypeOptions                             []TypeOptions `json:"TypeOptions"`
	CollapseSingleItemFolders               bool          `json:"CollapseSingleItemFolders"`
	MinResumePct                            int           `json:"MinResumePct"`
	MaxResumePct                            int           `json:"MaxResumePct"`
	MinResumeDurationSeconds                int           `json:"MinResumeDurationSeconds"`
	ThumbnailImagesIntervalSeconds          int           `json:"ThumbnailImagesIntervalSeconds"`
}

type PathInfo struct {
	Path        string `json:"Path"`
	NetworkPath string `json:"NetworkPath"`
}

type TypeOptions struct {
	Type                 string         `json:"Type"`
	MetadataFetchers     []string       `json:"MetadataFetchers"`
	MetadataFetcherOrder []string       `json:"MetadataFetcherOrder"`
	ImageFetchers        []string       `json:"ImageFetchers"`
	ImageFetcherOrder    []string       `json:"ImageFetcherOrder"`
	ImageOptions         []ImageOptions `json:"ImageOptions"`
}

type ImageOptions struct {
	Type     string `json:"Type"`
	Limit    int    `json:"Limit"`
	MinWidth int    `json:"MinWidth"`
}

// type MediaFolder struct {
// 	Name       string    `json:"Name"`
// 	Id         string    `json:"Id"`
// 	SubFolders SubFolder `json:"SubFolders"`
// }

type SubFolder struct {
	Name string `json:"Name"`
	Id   string `json:"Id"`
	Path string `json:"Path"`
}

type AddMedia struct {
	Name     string   `json:"Name"`
	Path     string   `json:"Path"`
	PathInfo PathInfo `json:"PathInfo"`
}
