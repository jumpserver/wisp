package service

// 与Core交互的API
const (
	UserProfileURL       = "/api/v1/users/profile/"                   // 获取当前用户的基本信息
	TerminalRegisterURL  = "/api/v1/terminal/terminal-registrations/" // 注册
	TerminalConfigURL    = "/api/v1/terminal/terminals/config/"       // 获取配置
	TerminalHeartBeatURL = "/api/v1/terminal/terminals/status/"
)

// 用户登陆认证使用的API
const (
	TokenAssetURL      = "/api/v1/authentication/connection-token/?token=%s" // Token name
	UserTokenAuthURL   = "/api/v1/authentication/tokens/"                    // 用户登录验证
	UserConfirmAuthURL = "/api/v1/authentication/login-confirm-ticket/status/"
	AuthMFASelectURL   = "/api/v1/authentication/mfa/select/" // 选择 MFA

)

// Session相关API
const (
	SessionListURL      = "/api/v1/terminal/sessions/"           //上传创建的资产会话session id
	SessionDetailURL    = "/api/v1/terminal/sessions/%s/"        // finish session的时候发送
	SessionReplayURL    = "/api/v1/terminal/sessions/%s/replay/" //上传录像
	SessionCommandURL   = "/api/v1/terminal/commands/"           //上传批量命令
	FinishTaskURL       = "/api/v1/terminal/tasks/%s/"
	JoinRoomValidateURL = "/api/v1/terminal/sessions/join/validate/"
	FTPLogListURL       = "/api/v1/audits/ftp-logs/" // 上传 ftp日志

	SessionLifecycleLogURL = "/api/v1/terminal/sessions/%s/lifecycle_log/"
)

// 授权相关API
const (
	UserPermsNodesListURL            = "/api/v1/perms/users/%s/nodes/"
	UserPermsNodeAssetsListURL       = "/api/v1/perms/users/%s/nodes/%s/assets/"
	UserPermsNodeTreeWithAssetURL    = "/api/v1/perms/users/%s/nodes/children-with-assets/tree/" // 资产树
	ValidateUserAssetPermissionURL   = "/api/v1/perms/asset-permissions/user/validate/"
	ValidateApplicationPermissionURL = "/api/v1/perms/application-permissions/user/validate/"
)

// 各资源详情相关API
const (
	UserDetailURL    = "/api/v1/users/users/%s/"
	AssetDetailURL   = "/api/v1/assets/assets/%s/"
	AssetPlatFormURL = "/api/v1/assets/assets/%s/platform/"

	DomainDetailWithGateways = "/api/v1/assets/domains/%s/?gateway=1"
)

const (
	NotificationCommandURL = "/api/v1/terminal/commands/insecure-command/"
)

const (
	PermissionURL = "/api/v1/perms/asset-permissions/user/actions/"
)

const (
	ShareCreateURL        = "/api/v1/terminal/session-sharings/"
	ShareSessionJoinURL   = "/api/v1/terminal/session-join-records/"
	ShareSessionFinishURL = "/api/v1/terminal/session-join-records/%s/finished/"
)

const (
	PublicSettingURL = "/api/v1/settings/public/"

	TicketSessionURL = "/api/v1/tickets/ticket-session-relation/"
)

// 数据库端口映射相关API
const (
	DBListenPortsURL = "/api/v1/terminal/db-listen-ports/"
	DBPortInfoURL    = "/api/v1/terminal/db-listen-ports/db-info/"

	SuperConnectTokenSecretURL = "/api/v1/authentication/super-connection-token/secret/"
	SuperConnectTokenInfoURL   = "/api/v1/authentication/super-connection-token/"
	SuperTokenRenewalURL       = "/api/v1/authentication/super-connection-token/renewal/"
	SuperConnectTokenCheckURL  = "/api/v1/authentication/super-connection-token/%s/check/"

	UserPermsAssetsURL = "/api/v1/perms/users/%s/assets/"

	AssetLoginConfirmURL = "/api/v1/acls/login-asset/check/"
	AclCommandReviewURL  = "/api/v1/acls/command-filter-acls/command-review/"
)
