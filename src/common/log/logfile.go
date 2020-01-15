package log

import (
	"common/config"
	"time"
)

const (
	defaultMaxAge   = 14 //keep log for 2 weeks
	defaultRotation = 12 //every 12 hours
)

//==============================================================================
// GetWriter 日志输出文件
// -- 参数1： 文件路径
// -- 设置日志文件分割条件、输出名称。
//==============================================================================
func GetWriter(filepath string) *rotatelogs.RotateLogs {
	maxAge := config.GetInt("common.logging.dumpFile.maxage")
	if maxAge == 0 {
		maxAge = defaultMaxAge
	}
	rotation := config.GetInt("common.logging.dumpFile.rotation")
	if rotation == 0 {
		rotation = defaultRotation
	}
	writer, err := rotatelogs.New(
		filepath+".%Y%m%d%H%M",
		rotatelogs.WithMaxAge(time.Hour*24*time.Duration(maxAge)),
		rotatelogs.WithRotationTime(time.Hour*time.Duration(rotation)),
		//rotatelogs.WithLinkName()
	)
	if err != nil {
		return nil
	}
	return writer
}
