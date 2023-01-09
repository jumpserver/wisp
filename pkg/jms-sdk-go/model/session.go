package model

import (
	"encoding/json"
	"strings"

	"github.com/jumpserver/wisp/pkg/jms-sdk-go/common"
)

var (
	_ json.Unmarshaler = (*LabelFiled)(nil)
	_ json.Marshaler   = (*LabelFiled)(nil)
)

type LabelFiled string

func (s LabelFiled) MarshalJSON() ([]byte, error) {
	return []byte(`"` + s + `"`), nil
}

func (s *LabelFiled) UnmarshalJSON(bytes []byte) error {
	var labelValue struct {
		Label string `json:"label"`
		Value string `json:"value"`
	}
	if err := json.Unmarshal(bytes, &labelValue); err != nil {
		return err
	}
	*s = LabelFiled(labelValue.Value)
	return nil
}

const (
	NORMALType  LabelFiled = "normal"
	TUNNELType  LabelFiled = "tunnel"
	COMMANDType LabelFiled = "command"
)

type Session struct {
	ID         string         `json:"id,omitempty"`
	User       string         `json:"user"` // "%s(%s)" Name Username
	Asset      string         `json:"asset"`
	Account    string         `json:"account"`
	LoginFrom  LabelFiled     `json:"login_from"`
	RemoteAddr string         `json:"remote_addr"`
	Protocol   string         `json:"protocol"`
	DateStart  common.UTCTime `json:"date_start"`
	OrgID      string         `json:"org_id"`
	UserID     string         `json:"user_id"`
	AssetID    string         `json:"asset_id"`
	Type       LabelFiled     `json:"type"`
}

type ReplayVersion string

const (
	UnKnown  ReplayVersion = ""
	Version2 ReplayVersion = "2"
	Version3 ReplayVersion = "3"
)

const (
	SuffixReplayGz = ".replay.gz"
	SuffixCastGz   = ".cast.gz"
	SuffixCast     = ".cast"
	SuffixGuac     = ".guac"
	SuffixGz       = ".gz"
)

var SuffixMap = map[ReplayVersion]string{
	Version2: SuffixReplayGz,
	Version3: SuffixCastGz,
}

func ParseReplayVersion(gzFile string, defaultValue ReplayVersion) ReplayVersion {
	for version, suffix := range SuffixMap {
		if strings.HasSuffix(gzFile, suffix) {
			return version
		}
	}
	return defaultValue
}

const (
	LoginFromWT LabelFiled = "WT"
	LoginFromST LabelFiled = "ST"
	LoginFromRT LabelFiled = "RT"
	LoginFromDT LabelFiled = "DT"
)
