package model

type TerminalConfig struct {
	AssetListPageSize   string        `json:"TERMINAL_ASSET_LIST_PAGE_SIZE"`
	AssetListSortBy     string        `json:"TERMINAL_ASSET_LIST_SORT_BY"`
	HeaderTitle         string        `json:"TERMINAL_HEADER_TITLE"`
	PasswordAuth        bool          `json:"TERMINAL_PASSWORD_AUTH"`
	PublicKeyAuth       bool          `json:"TERMINAL_PUBLIC_KEY_AUTH"`
	ReplayStorage       ReplayConfig  `json:"TERMINAL_REPLAY_STORAGE"`
	CommandStorage      CommandConfig `json:"TERMINAL_COMMAND_STORAGE"`
	SessionKeepDuration int           `json:"TERMINAL_SESSION_KEEP_DURATION"`
	TelnetRegex         string        `json:"TERMINAL_TELNET_REGEX"`
	MaxIdleTime         int           `json:"SECURITY_MAX_IDLE_TIME"`
	MaxSessionTime      int           `json:"SECURITY_MAX_SESSION_TIME"`
	HeartbeatDuration   int           `json:"TERMINAL_HEARTBEAT_INTERVAL"`
	HostKey             string        `json:"TERMINAL_HOST_KEY"`
	EnableSessionShare  bool          `json:"SECURITY_SESSION_SHARE"`
	LicenseContent      string        `json:"XPACK_LICENSE_CONTENT"`
	LicenseIsValid      bool          `json:"XPACK_LICENSE_IS_VALID"`

	GPTBaseURL string `json:"GPT_BASE_URL"`
	GPTApiKey  string `json:"GPT_API_KEY"`
	GPTProxy   string `json:"GPT_PROXY"`
	GPTModel   string `json:"GPT_MODEL"`
}

type Terminal struct {
	Name           string `json:"name"`
	Comment        string `json:"comment"`
	ServiceAccount struct {
		ID        string    `json:"id"`
		Name      string    `json:"name"`
		AccessKey AccessKey `json:"access_key"`
	} `json:"service_account"`
}

type TerminalTask struct {
	ID     string     `json:"id"`
	Name   string     `json:"name"`
	Args   string     `json:"args"`
	Kwargs TaskKwargs `json:"kwargs"`

	TokenStatus TokenCheckStatus `json:"-"`
}

const (
	TaskKillSession   = "kill_session"
	TaskLockSession   = "lock_session"
	TaskUnlockSession = "unlock_session"

	// TaskPermExpired TaskPermValid 非 api 数据，仅用于内部处理

	TaskPermExpired = "perm_expired"
	TaskPermValid   = "perm_valid"
)

type TaskKwargs struct {
	TerminatedBy  string `json:"terminated_by"`
	CreatedByUser string `json:"created_by"`
}

type ReplayConfig struct {
	TypeName string `json:"TYPE"`

	/*
		obs oss
	*/
	Endpoint  string `json:"ENDPOINT,omitempty"`
	Bucket    string `json:"BUCKET,omitempty"`
	AccessKey string `json:"ACCESS_KEY,omitempty"`
	SecretKey string `json:"SECRET_KEY,omitempty"`

	/*
		s3、 swift 需要
	*/

	Region string `json:"REGION,omitempty"`

	/*
		azure 专属
	*/
	AccountName    string `json:"ACCOUNT_NAME,omitempty"`
	AccountKey     string `json:"ACCOUNT_KEY,omitempty"`
	EndpointSuffix string `json:"ENDPOINT_SUFFIX,omitempty"`
	ContainerName  string `json:"CONTAINER_NAME,omitempty"`
}

type CommandConfig struct {
	TypeName string     `json:"TYPE"`
	Other    *OtherMeta `json:"OTHER,omitempty"`
	Index    string     `json:"INDEX,omitempty"`
	Hosts    []string   `json:"HOSTS,omitempty"`
	DocType  string     `json:"DOC_TYPE,omitempty"`
}

type OtherMeta struct {
	IgnoreVerifyCerts bool `json:"IGNORE_VERIFY_CERTS,omitempty"`
}
