package model

type ConnectToken struct {
	Id       string     `json:"id"`
	User     User       `json:"user"`
	Value    string     `json:"value"`
	Account  Account    `json:"account"`
	Actions  Actions    `json:"actions"`
	Asset    Asset      `json:"asset"`
	Protocol string     `json:"protocol"`
	Domain   *Domain    `json:"domain"`
	Gateway  *Gateway   `json:"gateway"`
	ExpireAt ExpireInfo `json:"expire_at"`
	OrgId    string     `json:"org_id"`
	OrgName  string     `json:"org_name"`
	Platform Platform   `json:"platform"`

	CommandFilterACLs []CommandACL `json:"command_filter_acls"`

	Ticket           *ObjectId   `json:"from_ticket,omitempty"`
	TicketInfo       interface{} `json:"from_ticket_info,omitempty"`
	FaceMonitorToken string      `json:"face_monitor_token,omitempty"`

	Code   string `json:"code"`
	Detail string `json:"detail"`
}

func (c *ConnectToken) Permission() Permission {
	var permission Permission
	permission.Actions = make([]string, 0, len(c.Actions))
	for i := range c.Actions {
		permission.Actions = append(permission.Actions, c.Actions[i].Value)
	}
	return permission
}

type ConnectTokenInfo struct {
	ID          string `json:"id"`
	Value       string `json:"value"`
	ExpireTime  int    `json:"expire_time"`
	AccountName string `json:"account_name"`
	Protocol    string `json:"protocol"`
}

// token 授权和过期状态

type TokenCheckStatus struct {
	Detail  string `json:"detail"`
	Code    string `json:"code"`
	Expired bool   `json:"expired"`
}

const (
	CodePermOk             = "perm_ok"
	CodePermAccountInvalid = "perm_account_invalid"
	CodePermExpired        = "perm_expired"
)
