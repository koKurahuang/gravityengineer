package log

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"
)

const (
	defaultFormat          = "[%time%] [%lvl%] [%module%] [%file%] → %msg%"
	defaultTimestampFormat = time.RFC3339

	nocolor = 0
	red     = 31
	green   = 32
	yellow  = 33
	blue    = 36
	gray    = 37
)

//==============================================================================
// log 格式 接口实现
//==============================================================================
type CustomFormatter struct {
	TimestampFormat string
	LogFormat       string
	EnableColors    bool
	sync.Once
}

//==============================================================================
// Format log 格式 接口实现
//==============================================================================
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	b = &bytes.Buffer{}
	output := f.LogFormat
	if output == "" {
		output = defaultFormat
	}
	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}
	//日志输出格式设置
	output = strings.Replace(output, "%time%", entry.Time.Format(timestampFormat), 1)
	output = strings.Replace(output, "%msg%", entry.Message, 1)
	level := strings.ToUpper(entry.Level.String())[0:4]
	output = strings.Replace(output, "%lvl%", level, 1)
	//output = strings.Replace(output, "%file%", FileInfo(5), 1)
	for k, v := range entry.Data {
		if s, ok := v.(string); ok {
			output = strings.Replace(output, "%"+k+"%", s, 1)
		}
	}
	if b.Len() > 0 {
		b.WriteByte(' ')
	}

	if f.EnableColors {
		f.printColored(b, entry, output)
	} else {
		b.WriteString(output)
	}
	b.WriteByte('\n')
	return b.Bytes(), nil
}

//==============================================================================
// printColored 终端日志输出颜色设定
// -- 终端日志输出颜色设定
//==============================================================================
func (f *CustomFormatter) printColored(b *bytes.Buffer, entry *logrus.Entry, str string) {
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel:
		levelColor = gray
	case logrus.WarnLevel:
		levelColor = yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = red
	default:
		levelColor = blue
	}

	fmt.Fprintf(b, " \x1b[%dm%s\x1b[0m", levelColor, str)

}
