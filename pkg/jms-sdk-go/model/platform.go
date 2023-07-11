package model

type Platform struct {
	ID   int    `json:"id"`
	Name string `json:"name"`

	Protocols PlatformProtocols `json:"protocols"`
	Category  LabelValue        `json:"category"`
	Charset   LabelValue        `json:"charset"`
	Type      LabelValue        `json:"type"`
	//SuEnabled     bool              `json:"su_enabled"`
	//SuMethod      string            `json:"su_method"`
	//DomainEnabled bool              `json:"domain_enabled"`
	//Comment       string            `json:"comment"`
}

type PlatformProtocols []PlatformProtocol

type PlatformProtocol struct {
	Protocol
	Setting map[string]interface{} `json:"setting"` // 参考 ProtocolSetting 里的字段
}

// 这个字段会频繁变动，所以不定义结构体，这里只记录场景的需要的字段

type ProtocolSetting struct {
	Security string `json:"security"`

	// chatgpt 专用
	ApiMode string `json:"api_mode"`

	//AutoFill         bool   `json:"auto_fill"`
	//UsernameSelector string `json:"username_selector"`
	//PasswordSelector string `json:"password_selector"`
	//SubmitSelector   string `json:"submit_selector"`
}

type Protocol struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Port int    `json:"port"`
}

type LabelValue struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type Action LabelValue

type SecretType LabelValue
