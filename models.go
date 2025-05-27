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
	AudioLanguagePreference    string        `json:"AudioLanguagePreference"`
	PlayDefaultAudioTrack      bool          `json:"PlayDefaultAudioTrack"`
	SubtitleLanguagePreference string        `json:"SubtitleLanguagePreference"`
	DisplayMissingEpisodes     bool          `json:"DisplayMissingEpisodes"`
	GroupedFolders             []interface{} `json:"GroupedFolders,omitempty"`
	SubtitleMode               string        `json:"SubtitleMode"`
	DisplayCollectionsView     bool          `json:"DisplayCollectionsView"`
	EnableLocalPassword        bool          `json:"EnableLocalPassword"`
	OrderedViews               []interface{} `json:"OrderedViews,omitempty"`
	LatestItemsExcludes        []interface{} `json:"LatestItemsExcludes,omitempty"`
	MyMediaExcludes            []interface{} `json:"MyMediaExcludes,omitempty"`
	HidePlayedInLatest         bool          `json:"HidePlayedInLatest"`
	RememberAudioSelections    bool          `json:"RememberAudioSelections"`
	RememberSubtitleSelections bool          `json:"RememberSubtitleSelections"`
	EnableNextEpisodeAutoPlay  bool          `json:"EnableNextEpisodeAutoPlay"`
	CastReceiverID             string        `json:"CastReceiverId"`
}

// DeNullConfiguration ensures there are no "null" fields in the given Configuration.
// Jellyfin isn't a fan of null.
func DeNullConfiguration(c *Configuration) {
	if c.GroupedFolders == nil {
		c.GroupedFolders = []interface{}{}
	}
	if c.OrderedViews == nil {
		c.OrderedViews = []interface{}{}
	}
	if c.LatestItemsExcludes == nil {
		c.LatestItemsExcludes = []interface{}{}
	}
	if c.MyMediaExcludes == nil {
		c.MyMediaExcludes = []interface{}{}
	}
}

// Policy stores a users permissions.
type Policy struct {
	IsAdministrator                  bool          `json:"IsAdministrator"`
	IsHidden                         bool          `json:"IsHidden"`
	IsDisabled                       bool          `json:"IsDisabled"`
	BlockedTags                      []interface{} `json:"BlockedTags,omitempty"`
	AllowedTags                      []interface{} `json:"AllowedTags"`
	EnableUserPreferenceAccess       bool          `json:"EnableUserPreferenceAccess"`
	AccessSchedules                  []interface{} `json:"AccessSchedules,omitempty"`
	BlockUnratedItems                []interface{} `json:"BlockUnratedItems,omitempty"`
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
	EnableContentDeletionFromFolders []interface{} `json:"EnableContentDeletionFromFolders,omitempty"`
	EnableContentDownloading         bool          `json:"EnableContentDownloading"`
	EnableSyncTranscoding            bool          `json:"EnableSyncTranscoding"`
	EnableMediaConversion            bool          `json:"EnableMediaConversion"`
	EnabledDevices                   []interface{} `json:"EnabledDevices,omitempty"`
	EnableAllDevices                 bool          `json:"EnableAllDevices"`
	EnabledChannels                  []interface{} `json:"EnabledChannels,omitempty"`
	EnableAllChannels                bool          `json:"EnableAllChannels"`
	EnabledFolders                   []string      `json:"EnabledFolders"`
	EnableAllFolders                 bool          `json:"EnableAllFolders"`
	InvalidLoginAttemptCount         int           `json:"InvalidLoginAttemptCount"`
	EnablePublicSharing              bool          `json:"EnablePublicSharing"`
	RemoteClientBitrateLimit         int           `json:"RemoteClientBitrateLimit"`
	AuthenticationProviderID         string        `json:"AuthenticationProviderId"`

	EnableCollectionManagement bool `json:"EnableCollectionManagement"`
	EnableSubtitleManagement   bool `json:"EnableSubtitleManagement"`
	EnableLyricManagement      bool `json:"EnableLyricManagement"`

	// Jellyfin Only
	ForceRemoteSourceTranscoding bool          `json:"ForceRemoteSourceTranscoding"`
	LoginAttemptsBeforeLockout   int           `json:"LoginAttemptsBeforeLockout"`
	MaxActiveSessions            int           `json:"MaxActiveSessions"`
	MaxParentalRating            *int          `json:"MaxParentalRating,omitempty"`
	BlockedMediaFolders          []interface{} `json:"BlockedMediaFolders,omitempty"`
	BlockedChannels              []interface{} `json:"BlockedChannels,omitempty"`
	PasswordResetProviderID      string        `json:"PasswordResetProviderId"`
	SyncPlayAccess               string        `json:"SyncPlayAccess"`
	// Emby Only
	IsHiddenRemotely           bool          `json:"IsHiddenRemotely"`
	IsHiddenFromUnusedDevices  bool          `json:"IsHiddenFromUnusedDevices"`
	IsTagBlockingModeInclusive bool          `json:"IsTagBlockingModeInclusive"`
	EnableSubtitleDownloading  bool          `json:"EnableSubtitleDownloading"`
	ExcludedSubFolders         []interface{} `json:"ExcludedSubFolders,omitempty"`
	SimultaneousStreamLimit    int           `json:"SimultaneousStreamLimit"`
}

// DeNullPolicy ensures there are no "null" fields in the given Policy.
// Jellyfin isn't a fan of null.
func DeNullPolicy(p *Policy) {
	if p.BlockedTags == nil {
		p.BlockedTags = []interface{}{}
	}
	if p.AllowedTags == nil {
		p.AllowedTags = []interface{}{}
	}
	if p.AccessSchedules == nil {
		p.AccessSchedules = []interface{}{}
	}
	if p.BlockUnratedItems == nil {
		p.BlockUnratedItems = []interface{}{}
	}
	if p.EnableContentDeletionFromFolders == nil {
		p.EnableContentDeletionFromFolders = []interface{}{}
	}
	if p.EnabledDevices == nil {
		p.EnabledDevices = []interface{}{}
	}
	if p.EnabledChannels == nil {
		p.EnabledChannels = []interface{}{}
	}
	if p.BlockedMediaFolders == nil {
		p.BlockedMediaFolders = []interface{}{}
	}
	if p.BlockedChannels == nil {
		p.BlockedChannels = []interface{}{}
	}
	if p.ExcludedSubFolders == nil {
		p.ExcludedSubFolders = []interface{}{}
	}
	if p.EnabledFolders == nil {
		p.EnabledFolders = []string{}
	}
}

type PasswordResetResponse struct {
	Success    bool     `json:"Success"`
	UsersReset []string `json:"UsersReset"`
}

type setPasswordRequest struct {
	Current       string `json:"CurrentPassword"`
	CurrentPw     string `json:"CurrentPw"`
	New           string `json:"NewPw"`
	ResetPassword bool   `json:"ResetPassword"`
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

// DeNullVirtualFolder ensures there are no "null" fields in the given VirtualFolder.
// Jellyfin isn't a fan of null.
func DeNullVirtualFolder(vf *VirtualFolder) {
	if vf.Locations == nil {
		vf.Locations = []string{}
	}
	lo := vf.LibraryOptions
	DeNullLibraryOptions(&lo)
	vf.LibraryOptions = lo
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

// DeNullLibraryOptions ensures there are no "null" fields in the given LibraryOptions.
// Jellyfin isn't a fan of null.
func DeNullLibraryOptions(lo *LibraryOptions) {
	if lo.PathInfos == nil {
		lo.PathInfos = []PathInfo{}
	}
	if lo.MetadataSavers == nil {
		lo.MetadataSavers = []string{}
	}
	if lo.DisabledLocalMetadataReaders == nil {
		lo.DisabledLocalMetadataReaders = []string{}
	}
	if lo.LocalMetadataReaderOrder == nil {
		lo.LocalMetadataReaderOrder = []string{}
	}
	if lo.DisabledSubtitleFetchers == nil {
		lo.DisabledSubtitleFetchers = []string{}
	}
	if lo.SubtitleFetcherOrder == nil {
		lo.SubtitleFetcherOrder = []string{}
	}
	if lo.SubtitleDownloadLanguages == nil {
		lo.SubtitleDownloadLanguages = []string{}
	}
	if lo.TypeOptions == nil {
		lo.TypeOptions = []TypeOptions{}
	}
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
