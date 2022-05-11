package model

type PublicSetting struct {
	LoginTitle string `json:"LOGIN_TITLE"`
	LogoURLS   struct {
		LogOut  string `json:"logo_logout"`
		Index   string `json:"logo_index"`
		Image   string `json:"login_image"`
		Favicon string `json:"favicon"`
	} `json:"LOGO_URLS"`
	EnableWatermark    bool `json:"SECURITY_WATERMARK_ENABLED"`
	EnableSessionShare bool `json:"SECURITY_SESSION_SHARE"`

	XpackEnabled bool `json:"XPACK_ENABLED"`
	ValidLicense bool `json:"XPACK_LICENSE_IS_VALID"`
}

/*

{
	"WINDOWS_SKIP_ALL_MANUAL_PASSWORD":false,
	"OLD_PASSWORD_HISTORY_LIMIT_COUNT":3,
	"SECURITY_MAX_IDLE_TIME":100,
	"SECURITY_VIEW_AUTH_NEED_MFA":true,
	"SECURITY_MFA_VERIFY_TTL":60,
	"SECURITY_COMMAND_EXECUTION":true,
	"SECURITY_PASSWORD_EXPIRATION_TIME":1000,
	"SECURITY_LUNA_REMEMBER_AUTH":true,
	"PASSWORD_RULE":{
		"SECURITY_PASSWORD_MIN_LENGTH":6,
		"SECURITY_ADMIN_USER_PASSWORD_MIN_LENGTH":8,
		"SECURITY_PASSWORD_UPPER_CASE":true,
		"SECURITY_PASSWORD_LOWER_CASE":true,
		"SECURITY_PASSWORD_NUMBER":true,
		"SECURITY_PASSWORD_SPECIAL_CHAR":true
	},
	"SECURITY_WATERMARK_ENABLED":true,
	"SECURITY_SESSION_SHARE":true,
	"XPACK_ENABLED":true,
	"XPACK_LICENSE_IS_VALID":true,
	"XPACK_LICENSE_INFO":{
		"corporation":"JumpServer"
	},
	"LOGIN_TITLE":"JumpServer Open Source Bastion Host",
	"LOGO_URLS":{
		"logo_logout":"/static/img/logo.png",
		"logo_index":"/static/img/logo_text.png",
		"login_image":"/static/img/login_image.jpg",
		"favicon":"/static/img/facio.ico"
	},
	"HELP_DOCUMENT_URL":"http://docs.jumpserver.org",
	"HELP_SUPPORT_URL":"",
	"AUTH_WECOM":true,
	"AUTH_DINGTALK":true,
	"AUTH_FEISHU":true,
	"XRDP_ENABLED":true,
	"TERMINAL_KOKO_HOST":"",
	"TERMINAL_KOKO_SSH_PORT":2222,
	"TERMINAL_MAGNUS_ENABLED":true,
	"TERMINAL_MAGNUS_HOST":"",
	"TERMINAL_MAGNUS_MYSQL_PORT":33060,
	"TERMINAL_MAGNUS_MARIADB_PORT":33061,
	"TERMINAL_MAGNUS_POSTGRE_PORT":54320,
	"ANNOUNCEMENT_ENABLED":true,
	"ANNOUNCEMENT":{
		"ID":"",
		"SUBJECT":"",
		"CONTENT":"",
		"LINK":""
	}
}
*/
