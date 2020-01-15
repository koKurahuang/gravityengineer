package log

import (
	"common/config"
	"common/log/funcNameHook"
	"common/log/modNameHook"
	"github.com/mattn/go-colorable"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"runtime"
)

// logging configs
type logConfig struct {
	module     string
	level      logrus.Level
	formatter  string
	colored    bool
	moduleList map[string]interface{}
	isDump     bool
	dumpPath   string
	dumpFormat string
	timestamp  string
}

//default filed
const defaultLevel = logrus.DebugLevel // todo debug for testing, info for UAT/PROD
const defaultFormatter = "[%time%] [%lvl%] [%module%] [%file%] → %msg%"
const defaultTimestamp = "2006-01-02 15:04:05.000 Z07:00"

var (
	modules *SyncMap // logger instance map
)

func init() {
	modules = NewSyncmap()
}

//==============================================================================
// InitLog 日志模块初始化
// -- 参数1： 无
// -- 外部调用 (peer node start) 读取配置文件 初始化各个模块日志实例
//==============================================================================
func InitLog() {
	conf := getLogConfig()
	for k, v := range conf.moduleList {
		conf.module = k
		conf.level, _ = logrus.ParseLevel(v.(string))
		resetLogInstance(conf)
	}
}

//==============================================================================
// resetLogInstance 各模块日志实例充值
// -- 参数1： 日志配置参数
// -- peer start时再按照yaml文件配置内容初始化log实例
//==============================================================================
func resetLogInstance(c *logConfig) {
	var loggerIns *logrus.Logger
	if modules.Check(c.module) {
		loggerIns = modules.Get(c.module).(*logrus.Logger)
	} else {
		loggerIns = newLogInstance(c.module)
	}
	updateLogInstance(loggerIns, c)
}

//==============================================================================
// updateLogInstance 各模块日志实例更新
// -- 参数1： 日志配置参数
// -- 各文件调用logger时，按照默认值初始化，在peer start时再按照yaml文件配置内容
// -- 更新个实例配置
// reset logger module instance with yaml file
//==============================================================================
// update logger instance
func updateLogInstance(ins *logrus.Logger, c *logConfig) {
	//ins := insEntry.Logger
	//console configs
	//ins.Out = os.Stdout
	var nullIO *nullWriter
	consoleFlg := config.GetBool("common.logging.enableconsole")
	if !consoleFlg {
		ins.Out = nullIO
	}
	ins.Formatter = c.getTextFormatter(c.colored)
	ins.Level = c.level

	// Hooks for dump log to file
	if c.isDump {
		// Hooks for all level message
		allWriter := GetWriter(path.Join(c.dumpPath, "all.log"))
		ins.Hooks.Add(lfshook.NewHook(lfshook.WriterMap{
			logrus.DebugLevel: allWriter,
			logrus.InfoLevel:  allWriter,
			logrus.WarnLevel:  allWriter,
			logrus.ErrorLevel: allWriter,
			logrus.FatalLevel: allWriter,
			logrus.PanicLevel: allWriter,
		}, c.getDumpFileFormatter().(logrus.Formatter)))
		// Hooks for only Error level message
		errWriter := GetWriter(path.Join(c.dumpPath, "error.log"))
		ins.Hooks.Add(lfshook.NewHook(lfshook.WriterMap{
			logrus.ErrorLevel: errWriter,
		}, c.getDumpFileFormatter().(logrus.Formatter)))
	}
	modules.Set(c.module, ins)
}

//==============================================================================
// getLogConfig 获取日志配置信息
// -- 参数1： 无
// -- 返回： log配置
// -- 从配置文件中读取log配置信息
//==============================================================================
func getLogConfig() *logConfig {
	//config.InitConfig()
	return &logConfig{
		formatter:  config.GetString("common.logging.format"),
		isDump:     config.GetBool("common.logging.dumpFile.enable"),
		dumpPath:   config.GetPath("common.logging.dumpFile.path"),
		dumpFormat: config.GetString("common.logging.dumpFile.format"),
		timestamp:  config.GetString("common.logging.timestamp"),
		moduleList: config.GetStringMap("common.logging.level"),
		colored:    config.GetBool("common.logging.enablecolors"),
	}
}

//==============================================================================
// defaultConfig 日志初始配置
// -- 参数1： 无
// -- 返回： log配置
// -- log模块还未初始化，其他包引入的时候初始配置
//==============================================================================
func defaultConfig() *logConfig {
	return &logConfig{
		module:    "",
		level:     defaultLevel,
		formatter: defaultFormatter,
		timestamp: defaultTimestamp,
		isDump:    false,
		colored:   true,
	}
}

//==============================================================================
// GetLogger 获取log实例
// -- 参数1： 需要使用log的模块名称
// -- 返回： log实例
// -- 各模块使用log功能时，需调用次方法获得log实例
//==============================================================================
// get logger instance with module name
func GetLogger(module string) *logrus.Logger {
	instance := modules.Get(module)
	if instance == nil {
		//return setModuleName(newLogInstance(module), module)
		//return newLogInstance(module).
		return newLogInstance(module)
	} else {
		//return setModuleName(instance.(*logrus.Logger), module)
		return instance.(*logrus.Logger)
	}
}

//==============================================================================
// newLogInstance 创建log实例
// -- 参数1： 需要使用log的模块名称
// -- 返回： log实例
// -- 未初始化时，根据模块名创建默认配置的log实例
//==============================================================================
// create logger instance with module name
func newLogInstance(module string) *logrus.Logger {
	c := defaultConfig()
	ins := logrus.New()
	//console configs
	//ins.Out = os.Stdout
	if runtime.GOOS == "windows" {
		ins.Out = colorable.NewColorableStdout()
	} else {
		ins.Out = os.Stdout
	}
	ins.Formatter = c.getTextFormatter(c.colored)
	ins.Level = c.level
	// Hooks for get file name/function name/line number
	ins.AddHook(funcNameHook.NewHook())
	ins.AddHook(modNameHook.NewModuleHook(module))
	modules.Set(module, ins)
	return ins
}

func ChangeLogLevel(logger *logrus.Logger, level string) {

	l, err := logrus.ParseLevel(level)
	if err != nil {
		logger.Error("change log level failed")
	}
	logger.SetLevel(l)
}

// set logger instance module name
func setModuleName(ins *logrus.Logger, module string) *logrus.Entry {
	return ins.WithFields(logrus.Fields{"module": module})
}

//==============================================================================
// getDumpFileFormatter 获取log文件输出格式
// -- 参数1： log配置--dump格式 (text/json)
// -- 设置log输出方式，有text和json两种可选配置
//==============================================================================
// get logger output file formatter
// support type text / json
func (c *logConfig) getDumpFileFormatter() interface{} {
	if c.dumpFormat == "text" {
		return c.getTextFormatter(false)
	} else {
		return &logrus.JSONFormatter{}
	}
}

//==============================================================================
// getTextFormatter 获取log控制台输出格式
// -- 参数1： log配置--输出格式 (text/json)
// -- 设置log输出方式，有text和json两种可选配置
//==============================================================================
// get logger text output formatter
func (c *logConfig) getTextFormatter(colorFlg bool) *CustomFormatter {
	return &CustomFormatter{
		TimestampFormat: c.timestamp,
		LogFormat:       c.formatter,
		EnableColors:    colorFlg,
	}
}

type nullWriter struct {
}

func (nw *nullWriter) Write(p []byte) (n int, err error) {
	return 0, nil
}
