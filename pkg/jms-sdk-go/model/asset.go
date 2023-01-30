package model

import (
	"fmt"
	"strings"
)

type Specific struct {
	// database
	DBName string `json:"db_name"`

	UseSSL           bool   `json:"use_ssl"`
	CaCert           string `json:"ca_cert"`
	ClientCert       string `json:"client_cert"`
	CertKey          string `json:"client_key"`
	AllowInvalidCert bool   `json:"allow_invalid_cert"`

	// web
	AutoFill         string `json:"autofill"`
	UsernameSelector string `json:"username_selector"`
	PasswordSelector string `json:"password_selector"`
	SubmitSelector   string `json:"submit_selector"`
}

type Asset struct {
	ID        string     `json:"id"`
	Address   string     `json:"address"`
	Name      string     `json:"name"`
	OrgID     string     `json:"org_id"`
	Protocols []Protocol `json:"protocols"`
	Specific  Specific   `json:"specific"`
	Platform  NameIntID  `json:"platform"`
	Domain    *NameStrID `json:"domain"`

	Comment  string `json:"comment"`
	OrgName  string `json:"org_name"`
	IsActive bool   `json:"is_active"` // 判断资产是否禁用

	Gateway *Gateway `json:"gateway,omitempty"`
}

type NameStrID struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type NameIntID struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (a *Asset) String() string {
	return fmt.Sprintf("%s(%s)", a.Name, a.Address)
}

func (a *Asset) ProtocolPort(protocol string) int {
	for _, item := range a.Protocols {
		protocolName := strings.ToLower(item.Name)
		if protocolName == strings.ToLower(protocol) {
			return item.Port
		}
	}
	return 0
}

func (a *Asset) SupportProtocols() []string {
	protocols := make([]string, 0, len(a.Protocols))
	for _, item := range a.Protocols {
		protocols = append(protocols, item.Name)
	}
	return protocols
}

func (a *Asset) IsSupportProtocol(protocol string) bool {
	for _, item := range a.Protocols {
		protocolName := strings.ToLower(item.Name)
		if protocolName == strings.ToLower(protocol) {
			return true
		}
	}
	return false
}

type Gateway struct {
	ID        string    `json:"id"`
	Name      string    `json:"Name"`
	Address   string    `json:"address"`
	Protocols Protocols `json:"protocols"`
	Account   Account   `json:"account"`
}

type Protocols []Protocol

func (p Protocols) GetProtocolPort(protocol string) int {
	for i := range p {
		if strings.ToLower(p[i].Name) == strings.ToLower(protocol) {
			return p[i].Port
		}
	}
	return 0
}
func (p Protocols) IsSupportProtocol(protocol string) bool {
	for _, item := range p {
		protocolName := strings.ToLower(item.Name)
		if protocolName == strings.ToLower(protocol) {
			return true
		}
	}
	return false
}

type Domain struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Gateways []Gateway `json:"gateways"`
}

const (
	ProtocolSSH    = "ssh"
	ProtocolTelnet = "telnet"
	ProtocolK8S    = "k8s"
	ProtocolMysql  = "mysql"
)
