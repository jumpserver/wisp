package model

import (
	"strings"

	"github.com/jumpserver/wisp/pkg/jms-sdk-go/common"
)

type Session struct {
	ID           string         `json:"id,omitempty"`
	User         string         `json:"user"` // "%s(%s)" Name Username
	Asset        string         `json:"asset"`
	SystemUser   string         `json:"system_user"`
	LoginFrom    string         `json:"login_from"`
	RemoteAddr   string         `json:"remote_addr"`
	Protocol     string         `json:"protocol"`
	DateStart    common.UTCTime `json:"date_start"`
	OrgID        string         `json:"org_id"`
	UserID       string         `json:"user_id"`
	AssetID      string         `json:"asset_id"`
	SystemUserID string         `json:"system_user_id"`
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
	LoginFromWT = "WT"
	LoginFromST = "ST"
	LoginFromRT = "RT"
	LoginFromDT = "DT"
)
