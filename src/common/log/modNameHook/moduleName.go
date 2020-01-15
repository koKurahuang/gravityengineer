package modNameHook

import "github.com/sirupsen/logrus"

//==============================================================================
// log Hook接口实现
//==============================================================================
type ModuleHook struct {
	Key   string
	Value string
}

//==============================================================================
// Levels log Hook接口实现
//==============================================================================
func (hook *ModuleHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

//==============================================================================
// Fire log Hook接口实现
//==============================================================================
func (hook *ModuleHook) Fire(entry *logrus.Entry) error {
	entry.Data[hook.Key] = hook.Value
	return nil
}

//==============================================================================
// NewModuleHook 模块名Hook
// -- 参数1： 模块名
// -- 返回： log Hook
// -- 在日志打印信息中添加模块名称
//==============================================================================
func NewModuleHook(module string) *ModuleHook {
	return &ModuleHook{
		"module",
		module,
	}
}
