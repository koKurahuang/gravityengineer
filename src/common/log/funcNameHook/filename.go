package funcNameHook

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
)

//==============================================================================
// log Hook接口实现
//==============================================================================
type Hook struct {
	Field     string
	Skip      int
	levels    []logrus.Level
	Formatter func(file, function string, line int) string
}

//==============================================================================
// Levels log Hook接口实现
//==============================================================================
func (hook *Hook) Levels() []logrus.Level {
	return hook.levels
}

//==============================================================================
// Fire log Hook接口实现
//==============================================================================
func (hook *Hook) Fire(entry *logrus.Entry) error {
	entry.Data[hook.Field] = hook.Formatter(findCaller(hook.Skip))
	return nil
}

//==============================================================================
// NewHook 文件名和行数Hook
// -- 参数1： 日志级别
// -- 返回： log Hook
// -- 获取打印日志的源文件以及具体行数
//==============================================================================
func NewHook(levels ...logrus.Level) *Hook {
	hook := Hook{
		Field:  "file",
		Skip:   5,
		levels: levels,
		Formatter: func(file, function string, line int) string {
			return fmt.Sprintf("%s:%s:%d", file, function, line)
		},
	}
	if len(hook.levels) == 0 {
		hook.levels = logrus.AllLevels
	}
	return &hook
}

//==============================================================================
// findCaller 获取调用栈信息
// -- 参数1： 调用栈顺序ID
// -- 返回： 文件名/方法/行数
// -- 获取打印日志的源文件以及具体行数
//==============================================================================
func findCaller(skip int) (string, string, int) {
	var (
		pc       uintptr
		file     string
		function string
		line     int
	)
	for i := 0; i < 10; i++ {
		pc, file, line = getCaller(skip + i)
		if !strings.HasPrefix(file, "logrus") {
			break
		}
	}
	if pc != 0 {
		frames := runtime.CallersFrames([]uintptr{pc})
		frame, _ := frames.Next()
		function = frame.Function
		slash := strings.LastIndex(function, "/")
		if slash > 0 {
			function = function[slash+1:]
		}
	}
	return file, function, line
}
func getCaller(skip int) (uintptr, string, int) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return 0, "", 0
	}
	n := 0
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			n += 1
			if n >= 2 {
				file = file[i+1:]
				break
			}
		}
	}
	return pc, file, line
}
