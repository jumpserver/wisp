// @generated
// This file is @generated by prost-build.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct User {
    #[prost(string, tag="1")]
    pub id: ::prost::alloc::string::String,
    #[prost(string, tag="2")]
    pub name: ::prost::alloc::string::String,
    #[prost(string, tag="3")]
    pub username: ::prost::alloc::string::String,
    #[prost(string, tag="4")]
    pub role: ::prost::alloc::string::String,
    #[prost(bool, tag="5")]
    pub is_valid: bool,
    #[prost(bool, tag="6")]
    pub is_active: bool,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Account {
    #[prost(string, tag="1")]
    pub id: ::prost::alloc::string::String,
    #[prost(string, tag="2")]
    pub name: ::prost::alloc::string::String,
    #[prost(string, tag="4")]
    pub username: ::prost::alloc::string::String,
    #[prost(string, tag="5")]
    pub secret: ::prost::alloc::string::String,
    #[prost(message, optional, tag="6")]
    pub secret_type: ::core::option::Option<LabelValue>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct LabelValue {
    #[prost(string, tag="1")]
    pub label: ::prost::alloc::string::String,
    #[prost(string, tag="2")]
    pub value: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Asset {
    #[prost(string, tag="1")]
    pub id: ::prost::alloc::string::String,
    #[prost(string, tag="2")]
    pub name: ::prost::alloc::string::String,
    #[prost(string, tag="3")]
    pub address: ::prost::alloc::string::String,
    #[prost(string, tag="4")]
    pub org_id: ::prost::alloc::string::String,
    #[prost(string, tag="5")]
    pub org_name: ::prost::alloc::string::String,
    #[prost(message, repeated, tag="6")]
    pub protocols: ::prost::alloc::vec::Vec<Protocol>,
    #[prost(message, optional, tag="7")]
    pub specific: ::core::option::Option<asset::Specific>,
}
/// Nested message and enum types in `Asset`.
pub mod asset {
    #[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Specific {
        #[prost(string, tag="1")]
        pub db_name: ::prost::alloc::string::String,
        #[prost(bool, tag="2")]
        pub use_ssl: bool,
        #[prost(string, tag="3")]
        pub ca_cert: ::prost::alloc::string::String,
        #[prost(string, tag="4")]
        pub client_cert: ::prost::alloc::string::String,
        #[prost(string, tag="5")]
        pub client_key: ::prost::alloc::string::String,
        #[prost(bool, tag="6")]
        pub allow_invalid_cert: bool,
        #[prost(string, tag="7")]
        pub auto_fill: ::prost::alloc::string::String,
        #[prost(string, tag="8")]
        pub username_selector: ::prost::alloc::string::String,
        #[prost(string, tag="9")]
        pub password_selector: ::prost::alloc::string::String,
        #[prost(string, tag="10")]
        pub submit_selector: ::prost::alloc::string::String,
        #[prost(string, tag="11")]
        pub script: ::prost::alloc::string::String,
        #[prost(string, tag="12")]
        pub http_proxy: ::prost::alloc::string::String,
        #[prost(string, tag="13")]
        pub pg_ssl_mode: ::prost::alloc::string::String,
    }
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Protocol {
    #[prost(string, tag="2")]
    pub name: ::prost::alloc::string::String,
    #[prost(int32, tag="1")]
    pub id: i32,
    #[prost(int32, tag="3")]
    pub port: i32,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Gateway {
    #[prost(string, tag="1")]
    pub id: ::prost::alloc::string::String,
    #[prost(string, tag="2")]
    pub name: ::prost::alloc::string::String,
    #[prost(string, tag="3")]
    pub ip: ::prost::alloc::string::String,
    #[prost(int32, tag="4")]
    pub port: i32,
    #[prost(string, tag="5")]
    pub protocol: ::prost::alloc::string::String,
    #[prost(string, tag="6")]
    pub username: ::prost::alloc::string::String,
    #[prost(string, tag="7")]
    pub password: ::prost::alloc::string::String,
    #[prost(string, tag="8")]
    pub private_key: ::prost::alloc::string::String,
}
#[derive(Clone, Copy, PartialEq, ::prost::Message)]
pub struct Permission {
    #[prost(bool, tag="1")]
    pub enable_connect: bool,
    #[prost(bool, tag="2")]
    pub enable_download: bool,
    #[prost(bool, tag="3")]
    pub enable_upload: bool,
    #[prost(bool, tag="4")]
    pub enable_copy: bool,
    #[prost(bool, tag="5")]
    pub enable_paste: bool,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct CommandAcl {
    #[prost(string, tag="1")]
    pub id: ::prost::alloc::string::String,
    #[prost(string, tag="2")]
    pub name: ::prost::alloc::string::String,
    #[prost(int32, tag="3")]
    pub priority: i32,
    #[prost(enumeration="command_acl::Action", tag="5")]
    pub action: i32,
    #[prost(bool, tag="6")]
    pub is_active: bool,
    #[prost(message, repeated, tag="7")]
    pub command_groups: ::prost::alloc::vec::Vec<CommandGroup>,
}
/// Nested message and enum types in `CommandACL`.
pub mod command_acl {
    #[derive(Clone, Copy, Debug, PartialEq, Eq, Hash, PartialOrd, Ord, ::prost::Enumeration)]
    #[repr(i32)]
    pub enum Action {
        Reject = 0,
        Accept = 1,
        Review = 2,
        Warning = 3,
        NotifyWarning = 4,
        Unknown = 5,
    }
    impl Action {
        /// String value of the enum field names used in the ProtoBuf definition.
        ///
        /// The values are not transformed in any way and thus are considered stable
        /// (if the ProtoBuf definition does not change) and safe for programmatic use.
        pub fn as_str_name(&self) -> &'static str {
            match self {
                Self::Reject => "Reject",
                Self::Accept => "Accept",
                Self::Review => "Review",
                Self::Warning => "Warning",
                Self::NotifyWarning => "NotifyWarning",
                Self::Unknown => "Unknown",
            }
        }
        /// Creates an enum from field names used in the ProtoBuf definition.
        pub fn from_str_name(value: &str) -> ::core::option::Option<Self> {
            match value {
                "Reject" => Some(Self::Reject),
                "Accept" => Some(Self::Accept),
                "Review" => Some(Self::Review),
                "Warning" => Some(Self::Warning),
                "NotifyWarning" => Some(Self::NotifyWarning),
                "Unknown" => Some(Self::Unknown),
                _ => None,
            }
        }
    }
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct CommandGroup {
    #[prost(string, tag="1")]
    pub id: ::prost::alloc::string::String,
    #[prost(string, tag="2")]
    pub name: ::prost::alloc::string::String,
    #[prost(string, tag="3")]
    pub content: ::prost::alloc::string::String,
    #[prost(string, tag="4")]
    pub r#type: ::prost::alloc::string::String,
    #[prost(string, tag="5")]
    pub pattern: ::prost::alloc::string::String,
    #[prost(bool, tag="6")]
    pub ignore_case: bool,
}
#[derive(Clone, Copy, PartialEq, ::prost::Message)]
pub struct ExpireInfo {
    #[prost(int64, tag="1")]
    pub expire_at: i64,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Session {
    #[prost(string, tag="1")]
    pub id: ::prost::alloc::string::String,
    #[prost(string, tag="2")]
    pub user: ::prost::alloc::string::String,
    #[prost(string, tag="3")]
    pub asset: ::prost::alloc::string::String,
    #[prost(string, tag="4")]
    pub account: ::prost::alloc::string::String,
    #[prost(enumeration="session::LoginFrom", tag="5")]
    pub login_from: i32,
    #[prost(string, tag="6")]
    pub remote_addr: ::prost::alloc::string::String,
    #[prost(string, tag="7")]
    pub protocol: ::prost::alloc::string::String,
    #[prost(int64, tag="8")]
    pub date_start: i64,
    #[prost(string, tag="9")]
    pub org_id: ::prost::alloc::string::String,
    #[prost(string, tag="10")]
    pub user_id: ::prost::alloc::string::String,
    #[prost(string, tag="11")]
    pub asset_id: ::prost::alloc::string::String,
    #[prost(string, tag="12")]
    pub account_id: ::prost::alloc::string::String,
    #[prost(string, tag="13")]
    pub token_id: ::prost::alloc::string::String,
}
/// Nested message and enum types in `Session`.
pub mod session {
    #[derive(Clone, Copy, Debug, PartialEq, Eq, Hash, PartialOrd, Ord, ::prost::Enumeration)]
    #[repr(i32)]
    pub enum LoginFrom {
        Wt = 0,
        St = 1,
        Rt = 2,
        Dt = 3,
    }
    impl LoginFrom {
        /// String value of the enum field names used in the ProtoBuf definition.
        ///
        /// The values are not transformed in any way and thus are considered stable
        /// (if the ProtoBuf definition does not change) and safe for programmatic use.
        pub fn as_str_name(&self) -> &'static str {
            match self {
                Self::Wt => "WT",
                Self::St => "ST",
                Self::Rt => "RT",
                Self::Dt => "DT",
            }
        }
        /// Creates an enum from field names used in the ProtoBuf definition.
        pub fn from_str_name(value: &str) -> ::core::option::Option<Self> {
            match value {
                "WT" => Some(Self::Wt),
                "ST" => Some(Self::St),
                "RT" => Some(Self::Rt),
                "DT" => Some(Self::Dt),
                _ => None,
            }
        }
    }
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct TokenStatus {
    #[prost(string, tag="1")]
    pub code: ::prost::alloc::string::String,
    #[prost(string, tag="2")]
    pub detail: ::prost::alloc::string::String,
    #[prost(bool, tag="3")]
    pub is_expired: bool,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct TerminalTask {
    #[prost(string, tag="1")]
    pub id: ::prost::alloc::string::String,
    #[prost(enumeration="TaskAction", tag="2")]
    pub action: i32,
    #[prost(string, tag="3")]
    pub session_id: ::prost::alloc::string::String,
    #[prost(string, tag="4")]
    pub terminated_by: ::prost::alloc::string::String,
    #[prost(string, tag="5")]
    pub created_by: ::prost::alloc::string::String,
    #[prost(message, optional, tag="6")]
    pub token_status: ::core::option::Option<TokenStatus>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct TokenAuthInfo {
    #[prost(string, tag="1")]
    pub key_id: ::prost::alloc::string::String,
    #[prost(string, tag="2")]
    pub secrete_id: ::prost::alloc::string::String,
    #[prost(message, optional, tag="3")]
    pub asset: ::core::option::Option<Asset>,
    #[prost(message, optional, tag="4")]
    pub user: ::core::option::Option<User>,
    #[prost(message, optional, tag="5")]
    pub account: ::core::option::Option<Account>,
    #[prost(message, optional, tag="6")]
    pub permission: ::core::option::Option<Permission>,
    #[prost(message, optional, tag="7")]
    pub expire_info: ::core::option::Option<ExpireInfo>,
    #[prost(message, repeated, tag="8")]
    pub filter_rules: ::prost::alloc::vec::Vec<CommandAcl>,
    #[prost(message, repeated, tag="9")]
    pub gateways: ::prost::alloc::vec::Vec<Gateway>,
    #[prost(message, optional, tag="10")]
    pub setting: ::core::option::Option<ComponentSetting>,
    #[prost(message, optional, tag="11")]
    pub platform: ::core::option::Option<Platform>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Platform {
    #[prost(int32, tag="1")]
    pub id: i32,
    #[prost(string, tag="2")]
    pub name: ::prost::alloc::string::String,
    #[prost(string, tag="3")]
    pub category: ::prost::alloc::string::String,
    #[prost(string, tag="4")]
    pub charset: ::prost::alloc::string::String,
    #[prost(string, tag="5")]
    pub r#type: ::prost::alloc::string::String,
    #[prost(message, repeated, tag="6")]
    pub protocols: ::prost::alloc::vec::Vec<PlatformProtocol>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct PlatformProtocol {
    #[prost(int32, tag="1")]
    pub id: i32,
    #[prost(string, tag="2")]
    pub name: ::prost::alloc::string::String,
    #[prost(int32, tag="3")]
    pub port: i32,
    #[prost(map="string, string", tag="4")]
    pub settings: ::std::collections::HashMap<::prost::alloc::string::String, ::prost::alloc::string::String>,
}
#[derive(Clone, Copy, PartialEq, ::prost::Message)]
pub struct ComponentSetting {
    #[prost(int32, tag="1")]
    pub max_idle_time: i32,
    #[prost(int32, tag="2")]
    pub max_session_time: i32,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Forward {
    #[prost(string, tag="1")]
    pub id: ::prost::alloc::string::String,
    #[prost(string, tag="2")]
    pub host: ::prost::alloc::string::String,
    #[prost(int32, tag="3")]
    pub port: i32,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct PublicSetting {
    #[prost(bool, tag="1")]
    pub xpack_enabled: bool,
    #[prost(bool, tag="2")]
    pub valid_license: bool,
    #[prost(string, tag="3")]
    pub gpt_base_url: ::prost::alloc::string::String,
    #[prost(string, tag="4")]
    pub gpt_api_key: ::prost::alloc::string::String,
    #[prost(string, tag="5")]
    pub gpt_proxy: ::prost::alloc::string::String,
    #[prost(string, tag="6")]
    pub gpt_model: ::prost::alloc::string::String,
    #[prost(string, tag="7")]
    pub license_content: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Cookie {
    #[prost(string, tag="1")]
    pub name: ::prost::alloc::string::String,
    #[prost(string, tag="2")]
    pub value: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct LifecycleLogData {
    #[prost(enumeration="lifecycle_log_data::EventType", tag="1")]
    pub event: i32,
    #[prost(string, tag="2")]
    pub reason: ::prost::alloc::string::String,
    #[prost(string, tag="3")]
    pub user: ::prost::alloc::string::String,
}
/// Nested message and enum types in `LifecycleLogData`.
pub mod lifecycle_log_data {
    #[derive(Clone, Copy, Debug, PartialEq, Eq, Hash, PartialOrd, Ord, ::prost::Enumeration)]
    #[repr(i32)]
    pub enum EventType {
        AssetConnectSuccess = 0,
        AssetConnectFinished = 1,
        CreateShareLink = 2,
        UserJoinSession = 3,
        UserLeaveSession = 4,
        AdminJoinMonitor = 5,
        AdminExitMonitor = 6,
        ReplayConvertStart = 7,
        ReplayConvertSuccess = 8,
        ReplayConvertFailure = 9,
        ReplayUploadStart = 10,
        ReplayUploadSuccess = 11,
        ReplayUploadFailure = 12,
    }
    impl EventType {
        /// String value of the enum field names used in the ProtoBuf definition.
        ///
        /// The values are not transformed in any way and thus are considered stable
        /// (if the ProtoBuf definition does not change) and safe for programmatic use.
        pub fn as_str_name(&self) -> &'static str {
            match self {
                Self::AssetConnectSuccess => "AssetConnectSuccess",
                Self::AssetConnectFinished => "AssetConnectFinished",
                Self::CreateShareLink => "CreateShareLink",
                Self::UserJoinSession => "UserJoinSession",
                Self::UserLeaveSession => "UserLeaveSession",
                Self::AdminJoinMonitor => "AdminJoinMonitor",
                Self::AdminExitMonitor => "AdminExitMonitor",
                Self::ReplayConvertStart => "ReplayConvertStart",
                Self::ReplayConvertSuccess => "ReplayConvertSuccess",
                Self::ReplayConvertFailure => "ReplayConvertFailure",
                Self::ReplayUploadStart => "ReplayUploadStart",
                Self::ReplayUploadSuccess => "ReplayUploadSuccess",
                Self::ReplayUploadFailure => "ReplayUploadFailure",
            }
        }
        /// Creates an enum from field names used in the ProtoBuf definition.
        pub fn from_str_name(value: &str) -> ::core::option::Option<Self> {
            match value {
                "AssetConnectSuccess" => Some(Self::AssetConnectSuccess),
                "AssetConnectFinished" => Some(Self::AssetConnectFinished),
                "CreateShareLink" => Some(Self::CreateShareLink),
                "UserJoinSession" => Some(Self::UserJoinSession),
                "UserLeaveSession" => Some(Self::UserLeaveSession),
                "AdminJoinMonitor" => Some(Self::AdminJoinMonitor),
                "AdminExitMonitor" => Some(Self::AdminExitMonitor),
                "ReplayConvertStart" => Some(Self::ReplayConvertStart),
                "ReplayConvertSuccess" => Some(Self::ReplayConvertSuccess),
                "ReplayConvertFailure" => Some(Self::ReplayConvertFailure),
                "ReplayUploadStart" => Some(Self::ReplayUploadStart),
                "ReplayUploadSuccess" => Some(Self::ReplayUploadSuccess),
                "ReplayUploadFailure" => Some(Self::ReplayUploadFailure),
                _ => None,
            }
        }
    }
}
#[derive(Clone, Copy, Debug, PartialEq, Eq, Hash, PartialOrd, Ord, ::prost::Enumeration)]
#[repr(i32)]
pub enum TaskAction {
    KillSession = 0,
    LockSession = 1,
    UnlockSession = 2,
    TokenPermExpired = 3,
    TokenPermValid = 4,
}
impl TaskAction {
    /// String value of the enum field names used in the ProtoBuf definition.
    ///
    /// The values are not transformed in any way and thus are considered stable
    /// (if the ProtoBuf definition does not change) and safe for programmatic use.
    pub fn as_str_name(&self) -> &'static str {
        match self {
            Self::KillSession => "KillSession",
            Self::LockSession => "LockSession",
            Self::UnlockSession => "UnlockSession",
            Self::TokenPermExpired => "TokenPermExpired",
            Self::TokenPermValid => "TokenPermValid",
        }
    }
    /// Creates an enum from field names used in the ProtoBuf definition.
    pub fn from_str_name(value: &str) -> ::core::option::Option<Self> {
        match value {
            "KillSession" => Some(Self::KillSession),
            "LockSession" => Some(Self::LockSession),
            "UnlockSession" => Some(Self::UnlockSession),
            "TokenPermExpired" => Some(Self::TokenPermExpired),
            "TokenPermValid" => Some(Self::TokenPermValid),
            _ => None,
        }
    }
}
#[derive(Clone, Copy, Debug, PartialEq, Eq, Hash, PartialOrd, Ord, ::prost::Enumeration)]
#[repr(i32)]
pub enum RiskLevel {
    Normal = 0,
    Warning = 1,
    Reject = 2,
    ReviewReject = 3,
    ReviewAccept = 4,
    ReviewCancel = 5,
}
impl RiskLevel {
    /// String value of the enum field names used in the ProtoBuf definition.
    ///
    /// The values are not transformed in any way and thus are considered stable
    /// (if the ProtoBuf definition does not change) and safe for programmatic use.
    pub fn as_str_name(&self) -> &'static str {
        match self {
            Self::Normal => "Normal",
            Self::Warning => "Warning",
            Self::Reject => "Reject",
            Self::ReviewReject => "ReviewReject",
            Self::ReviewAccept => "ReviewAccept",
            Self::ReviewCancel => "ReviewCancel",
        }
    }
    /// Creates an enum from field names used in the ProtoBuf definition.
    pub fn from_str_name(value: &str) -> ::core::option::Option<Self> {
        match value {
            "Normal" => Some(Self::Normal),
            "Warning" => Some(Self::Warning),
            "Reject" => Some(Self::Reject),
            "ReviewReject" => Some(Self::ReviewReject),
            "ReviewAccept" => Some(Self::ReviewAccept),
            "ReviewCancel" => Some(Self::ReviewCancel),
            _ => None,
        }
    }
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct FaceRecognitionCallbackRequest {
    #[prost(string, tag="1")]
    pub token: ::prost::alloc::string::String,
    #[prost(bool, tag="2")]
    pub success: bool,
    #[prost(string, tag="3")]
    pub error_message: ::prost::alloc::string::String,
    #[prost(string, tag="4")]
    pub face_code: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct FaceRecognitionCallbackResponse {
    #[prost(message, optional, tag="1")]
    pub status: ::core::option::Option<Status>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct AssetLoginTicketRequest {
    #[prost(string, tag="1")]
    pub user_id: ::prost::alloc::string::String,
    #[prost(string, tag="2")]
    pub asset_id: ::prost::alloc::string::String,
    #[prost(string, tag="4")]
    pub account_username: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct AssetLoginTicketResponse {
    #[prost(message, optional, tag="1")]
    pub status: ::core::option::Option<Status>,
    #[prost(message, optional, tag="2")]
    pub ticket_info: ::core::option::Option<TicketInfo>,
    #[prost(bool, tag="3")]
    pub need_confirm: bool,
    #[prost(string, tag="4")]
    pub ticket_id: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Status {
    #[prost(bool, tag="1")]
    pub ok: bool,
    #[prost(string, tag="2")]
    pub err: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct TokenRequest {
    #[prost(string, tag="1")]
    pub token: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct TokenResponse {
    #[prost(message, optional, tag="1")]
    pub status: ::core::option::Option<Status>,
    #[prost(message, optional, tag="2")]
    pub data: ::core::option::Option<TokenAuthInfo>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct SessionCreateRequest {
    #[prost(message, optional, tag="1")]
    pub data: ::core::option::Option<Session>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct SessionCreateResponse {
    #[prost(message, optional, tag="1")]
    pub status: ::core::option::Option<Status>,
    #[prost(message, optional, tag="2")]
    pub data: ::core::option::Option<Session>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct SessionFinishRequest {
    #[prost(string, tag="1")]
    pub id: ::prost::alloc::string::String,
    #[prost(bool, tag="2")]
    pub success: bool,
    #[prost(int64, tag="3")]
    pub date_end: i64,
    #[prost(string, tag="4")]
    pub err: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct SessionFinishResp {
    #[prost(message, optional, tag="1")]
    pub status: ::core::option::Option<Status>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ReplayRequest {
    #[prost(string, tag="1")]
    pub session_id: ::prost::alloc::string::String,
    #[prost(string, tag="2")]
    pub replay_file_path: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ReplayResponse {
    #[prost(message, optional, tag="1")]
    pub status: ::core::option::Option<Status>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct CommandRequest {
    #[prost(string, tag="1")]
    pub sid: ::prost::alloc::string::String,
    #[prost(string, tag="2")]
    pub org_id: ::prost::alloc::string::String,
    #[prost(string, tag="3")]
    pub input: ::prost::alloc::string::String,
    #[prost(string, tag="4")]
    pub output: ::prost::alloc::string::String,
    #[prost(string, tag="5")]
    pub user: ::prost::alloc::string::String,
    #[prost(string, tag="6")]
    pub asset: ::prost::alloc::string::String,
    #[prost(string, tag="7")]
    pub account: ::prost::alloc::string::String,
    #[prost(int64, tag="8")]
    pub timestamp: i64,
    #[prost(enumeration="RiskLevel", tag="9")]
    pub risk_level: i32,
    #[prost(string, tag="10")]
    pub cmd_acl_id: ::prost::alloc::string::String,
    #[prost(string, tag="11")]
    pub cmd_group_id: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct CommandResponse {
    #[prost(message, optional, tag="1")]
    pub status: ::core::option::Option<Status>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct FinishedTaskRequest {
    #[prost(string, tag="1")]
    pub task_id: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct TaskResponse {
    #[prost(message, optional, tag="1")]
    pub task: ::core::option::Option<TerminalTask>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct RemainReplayRequest {
    #[prost(string, tag="1")]
    pub replay_dir: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct RemainReplayResponse {
    #[prost(message, optional, tag="1")]
    pub status: ::core::option::Option<Status>,
    #[prost(string, repeated, tag="2")]
    pub success_files: ::prost::alloc::vec::Vec<::prost::alloc::string::String>,
    #[prost(string, repeated, tag="3")]
    pub failure_files: ::prost::alloc::vec::Vec<::prost::alloc::string::String>,
    #[prost(string, repeated, tag="4")]
    pub failure_errs: ::prost::alloc::vec::Vec<::prost::alloc::string::String>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct StatusResponse {
    #[prost(message, optional, tag="1")]
    pub status: ::core::option::Option<Status>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct CommandConfirmRequest {
    #[prost(string, tag="1")]
    pub session_id: ::prost::alloc::string::String,
    #[prost(string, tag="2")]
    pub cmd_acl_id: ::prost::alloc::string::String,
    #[prost(string, tag="3")]
    pub cmd: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ReqInfo {
    #[prost(string, tag="1")]
    pub method: ::prost::alloc::string::String,
    #[prost(string, tag="2")]
    pub url: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct CommandConfirmResponse {
    #[prost(message, optional, tag="1")]
    pub status: ::core::option::Option<Status>,
    #[prost(message, optional, tag="2")]
    pub info: ::core::option::Option<TicketInfo>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct TicketInfo {
    #[prost(message, optional, tag="1")]
    pub check_req: ::core::option::Option<ReqInfo>,
    #[prost(message, optional, tag="2")]
    pub cancel_req: ::core::option::Option<ReqInfo>,
    #[prost(string, tag="3")]
    pub ticket_detail_url: ::prost::alloc::string::String,
    #[prost(string, repeated, tag="4")]
    pub reviewers: ::prost::alloc::vec::Vec<::prost::alloc::string::String>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct TicketRequest {
    #[prost(message, optional, tag="1")]
    pub req: ::core::option::Option<ReqInfo>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct TicketStateResponse {
    #[prost(message, optional, tag="1")]
    pub data: ::core::option::Option<TicketState>,
    #[prost(message, optional, tag="2")]
    pub status: ::core::option::Option<Status>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct TicketState {
    #[prost(enumeration="ticket_state::State", tag="1")]
    pub state: i32,
    #[prost(string, tag="2")]
    pub processor: ::prost::alloc::string::String,
}
/// Nested message and enum types in `TicketState`.
pub mod ticket_state {
    #[derive(Clone, Copy, Debug, PartialEq, Eq, Hash, PartialOrd, Ord, ::prost::Enumeration)]
    #[repr(i32)]
    pub enum State {
        Open = 0,
        Approved = 1,
        Rejected = 2,
        Closed = 3,
    }
    impl State {
        /// String value of the enum field names used in the ProtoBuf definition.
        ///
        /// The values are not transformed in any way and thus are considered stable
        /// (if the ProtoBuf definition does not change) and safe for programmatic use.
        pub fn as_str_name(&self) -> &'static str {
            match self {
                Self::Open => "Open",
                Self::Approved => "Approved",
                Self::Rejected => "Rejected",
                Self::Closed => "Closed",
            }
        }
        /// Creates an enum from field names used in the ProtoBuf definition.
        pub fn from_str_name(value: &str) -> ::core::option::Option<Self> {
            match value {
                "Open" => Some(Self::Open),
                "Approved" => Some(Self::Approved),
                "Rejected" => Some(Self::Rejected),
                "Closed" => Some(Self::Closed),
                _ => None,
            }
        }
    }
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ForwardRequest {
    #[prost(string, tag="1")]
    pub host: ::prost::alloc::string::String,
    #[prost(int32, tag="2")]
    pub port: i32,
    #[prost(message, repeated, tag="3")]
    pub gateways: ::prost::alloc::vec::Vec<Gateway>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ForwardDeleteRequest {
    #[prost(string, tag="1")]
    pub id: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ForwardResponse {
    #[prost(message, optional, tag="1")]
    pub status: ::core::option::Option<Status>,
    #[prost(string, tag="2")]
    pub id: ::prost::alloc::string::String,
    #[prost(string, tag="3")]
    pub host: ::prost::alloc::string::String,
    #[prost(int32, tag="4")]
    pub port: i32,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct PublicSettingResponse {
    #[prost(message, optional, tag="1")]
    pub status: ::core::option::Option<Status>,
    #[prost(message, optional, tag="2")]
    pub data: ::core::option::Option<PublicSetting>,
}
#[derive(Clone, Copy, PartialEq, ::prost::Message)]
pub struct Empty {
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct ListenPortResponse {
    #[prost(message, optional, tag="1")]
    pub status: ::core::option::Option<Status>,
    #[prost(int32, repeated, tag="2")]
    pub ports: ::prost::alloc::vec::Vec<i32>,
}
#[derive(Clone, Copy, PartialEq, ::prost::Message)]
pub struct PortInfoRequest {
    #[prost(int32, tag="1")]
    pub port: i32,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct PortInfoResponse {
    #[prost(message, optional, tag="1")]
    pub status: ::core::option::Option<Status>,
    #[prost(message, optional, tag="2")]
    pub data: ::core::option::Option<PortInfo>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct PortInfo {
    #[prost(message, optional, tag="1")]
    pub asset: ::core::option::Option<Asset>,
    #[prost(message, repeated, tag="2")]
    pub gateways: ::prost::alloc::vec::Vec<Gateway>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct PortFailure {
    #[prost(int32, tag="1")]
    pub port: i32,
    #[prost(string, tag="2")]
    pub reason: ::prost::alloc::string::String,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct PortFailureRequest {
    #[prost(message, repeated, tag="1")]
    pub data: ::prost::alloc::vec::Vec<PortFailure>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct CookiesRequest {
    #[prost(message, repeated, tag="1")]
    pub cookies: ::prost::alloc::vec::Vec<Cookie>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct UserResponse {
    #[prost(message, optional, tag="1")]
    pub status: ::core::option::Option<Status>,
    #[prost(message, optional, tag="2")]
    pub data: ::core::option::Option<User>,
}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct SessionLifecycleLogRequest {
    #[prost(string, tag="1")]
    pub session_id: ::prost::alloc::string::String,
    #[prost(enumeration="session_lifecycle_log_request::EventType", tag="2")]
    pub event: i32,
    #[prost(string, tag="3")]
    pub reason: ::prost::alloc::string::String,
    #[prost(string, tag="4")]
    pub user: ::prost::alloc::string::String,
}
/// Nested message and enum types in `SessionLifecycleLogRequest`.
pub mod session_lifecycle_log_request {
    #[derive(Clone, Copy, Debug, PartialEq, Eq, Hash, PartialOrd, Ord, ::prost::Enumeration)]
    #[repr(i32)]
    pub enum EventType {
        AssetConnectSuccess = 0,
        AssetConnectFinished = 1,
        CreateShareLink = 2,
        UserJoinSession = 3,
        UserLeaveSession = 4,
        AdminJoinMonitor = 5,
        AdminExitMonitor = 6,
        ReplayConvertStart = 7,
        ReplayConvertSuccess = 8,
        ReplayConvertFailure = 9,
        ReplayUploadStart = 10,
        ReplayUploadSuccess = 11,
        ReplayUploadFailure = 12,
    }
    impl EventType {
        /// String value of the enum field names used in the ProtoBuf definition.
        ///
        /// The values are not transformed in any way and thus are considered stable
        /// (if the ProtoBuf definition does not change) and safe for programmatic use.
        pub fn as_str_name(&self) -> &'static str {
            match self {
                Self::AssetConnectSuccess => "AssetConnectSuccess",
                Self::AssetConnectFinished => "AssetConnectFinished",
                Self::CreateShareLink => "CreateShareLink",
                Self::UserJoinSession => "UserJoinSession",
                Self::UserLeaveSession => "UserLeaveSession",
                Self::AdminJoinMonitor => "AdminJoinMonitor",
                Self::AdminExitMonitor => "AdminExitMonitor",
                Self::ReplayConvertStart => "ReplayConvertStart",
                Self::ReplayConvertSuccess => "ReplayConvertSuccess",
                Self::ReplayConvertFailure => "ReplayConvertFailure",
                Self::ReplayUploadStart => "ReplayUploadStart",
                Self::ReplayUploadSuccess => "ReplayUploadSuccess",
                Self::ReplayUploadFailure => "ReplayUploadFailure",
            }
        }
        /// Creates an enum from field names used in the ProtoBuf definition.
        pub fn from_str_name(value: &str) -> ::core::option::Option<Self> {
            match value {
                "AssetConnectSuccess" => Some(Self::AssetConnectSuccess),
                "AssetConnectFinished" => Some(Self::AssetConnectFinished),
                "CreateShareLink" => Some(Self::CreateShareLink),
                "UserJoinSession" => Some(Self::UserJoinSession),
                "UserLeaveSession" => Some(Self::UserLeaveSession),
                "AdminJoinMonitor" => Some(Self::AdminJoinMonitor),
                "AdminExitMonitor" => Some(Self::AdminExitMonitor),
                "ReplayConvertStart" => Some(Self::ReplayConvertStart),
                "ReplayConvertSuccess" => Some(Self::ReplayConvertSuccess),
                "ReplayConvertFailure" => Some(Self::ReplayConvertFailure),
                "ReplayUploadStart" => Some(Self::ReplayUploadStart),
                "ReplayUploadSuccess" => Some(Self::ReplayUploadSuccess),
                "ReplayUploadFailure" => Some(Self::ReplayUploadFailure),
                _ => None,
            }
        }
    }
}
include!("message.tonic.rs");
// @@protoc_insertion_point(module)
