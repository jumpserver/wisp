package common

import (
	"strings"

	modelCommon "github.com/jumpserver-dev/sdk-go/common"
	"github.com/jumpserver-dev/sdk-go/model"
)

/*
koko   文件名为 sid | sid.replay.gz | sid.cast | sid.cast.gz
lion   文件名为 sid | sid.replay.gz
omnidb 文件名为 sid.cast | sid.cast.gz
xrdp   文件名为 sid.guac

如果存在日期目录，targetDate 使用日期目录的
文件路径名称中解析 录像文件信息

*/

var suffixesMap = map[string]model.ReplayVersion{
	model.SuffixGuac:     model.Version2,
	model.SuffixCast:     model.Version3,
	model.SuffixCastGz:   model.Version3,
	model.SuffixReplayGz: model.Version2,
}

type RemainReplay struct {
	Id          string // session id
	TargetDate  string
	AbsFilePath string
	Version     model.ReplayVersion
	IsGzip      bool
}

func (r *RemainReplay) TargetPath() string {
	gzFilename := r.GetGzFilename()
	return strings.Join([]string{r.TargetDate, gzFilename}, "/")
}

func (r *RemainReplay) GetGzFilename() string {
	suffixGz := ".replay.gz"
	switch r.Version {
	case model.Version3:
		suffixGz = ".cast.gz"
	case model.Version2:
		suffixGz = ".replay.gz"
	}
	return r.Id + suffixGz
}

func ParseReplaySessionID(filename string) (string, bool) {
	if len(filename) == 36 && modelCommon.IsUUID(filename) {
		return filename, true
	}
	sid := strings.Split(filename, ".")[0]
	if !modelCommon.IsUUID(sid) {
		return "", false
	}
	return sid, true
}

func ParseReplayVersion(filename string) (model.ReplayVersion, bool) {
	for suffix := range suffixesMap {
		if strings.HasSuffix(filename, suffix) {
			return suffixesMap[suffix], true

		}
	}
	return model.UnKnown, false
}
