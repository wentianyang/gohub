package logger

import (
	"encoding/json"
	"fmt"
	"gohub/pkg/app"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 全局 Logger 对象
var Logger *zap.Logger

func InitLogger(fileName string, maxSize, maxBackup, maxAge int, compress bool, logType string, level string) {

	// 获取日志写入介质
	writeSyncer := getLogWriter(fileName, maxSize, maxBackup, maxAge, compress, logType)

	// 设置日志等级, 具体详见 config/log.go 文件
	logLevel := new(zapcore.Level)
	if err := logLevel.UnmarshalText([]byte(level)); err != nil {
		fmt.Println("日志初始化错误, 日志级别设置有误. 请修改 config/log.go 文件中的 log.level 配置项")
	}

	// 初始化 Core
	core := zapcore.NewCore(getEncoder(), writeSyncer, logLevel)

	// 初始化 Logger
	Logger = zap.New(core,
		zap.AddCaller(),                   // 调用文件和行号,内部使用 runtime.Caller
		zap.AddCallerSkip(1),              // 封装了一层,调用文件去除一层(runtime.Caller(1))
		zap.AddStacktrace(zap.ErrorLevel), // Error 时才会显示 stacktrace
	)

	// 将自定义的 logger 替换为全局的 logger
	// zap.L().Fatal() 调用时,就会使用我们自定义的 Logger
	zap.ReplaceGlobals(Logger)
}

// 设置日志的存储格式
func getEncoder() zapcore.Encoder {
	// 日志格式规则
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller", // 代码调用, 如 paginator/paginator.go: 158
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,        // 每行日志的结尾添加 "\n"
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // 日志级别名称大写, 如 ERROR、INFO
		EncodeTime:     customTimeEncoder,                // 时间格式,我们自定义为 2006-01-02 15:22:02
		EncodeDuration: zapcore.SecondsDurationEncoder,   // 执行时间,以秒为单位
		EncodeCaller:   zapcore.ShortCallerEncoder,       // Caller 短格式,如: types/converter.go:1, 长合适为绝对路径
	}

	// 本地配置
	if app.IsLocal() {
		// 终端输出的关键词高亮
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		// 本地设置内置的 Console 解码器 (支持 stacktrace 换行)
		return zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 线上环境使用 JSON 解码器
	return zapcore.NewJSONEncoder(encoderConfig)
}

// 自定义友好的时间格式
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func getLogWriter(fileName string, maxSize, maxBackup, maxAge int, compress bool, logType string) zapcore.WriteSyncer {
	// 如果配置了按照日期记录日志文件
	if logType == "daily" {
		logname := time.Now().Format("2006-01-02.log")
		fileName = strings.ReplaceAll(fileName, "logs.log", logname)
	}

	// 滚动日志, 详见 config/log.go
	lumberJackLogger := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
		Compress:   compress,
	}

	// 配置输出介质
	if app.IsLocal() {
		// 本地开发终端打印和记录文件
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
	} else {
		// 生产环境只记录文件
		return zapcore.AddSync(lumberJackLogger)
	}
}

// Dump 调试专用,不会中断程序,在终端打印出 waring 消息
// 第一个参数会使用 json.Marshal 进行渲染,第二个参数消息(可选)
// 		logger.Dump(user.User{Name:"test"})
// 		logger.Dump(user.User{Name:"test"}, "用户信息")
func Dump(value interface{}, msg ...string) {
	valueString := jsonString(value)
	// 判断第二个参数是否传参 msg
	if len(msg) > 0 {
		Logger.Warn("Dump", zap.String(msg[0], valueString))
	} else {
		Logger.Warn("Dump", zap.String("data", valueString))
	}
}

// LogIf 当 err != nil 时记录 error 等级的日志
func LogIf(err error) {
	if err != nil {
		Logger.Error("Error Occurred", zap.Error(err))
	}
}

// LogWarnIf 当 err != nil 时记录 warning 等级的日志
func LogWarnIf(err error) {
	if err != nil {
		Logger.Warn("Error Occurred", zap.Error(err))
	}
}

// LogInfoIf 当 err != nil 时记录 info 等级的日志
func LogInfoIf(err error) {
	if err != nil {
		Logger.Info("Error Occurred", zap.Error(err))
	}
}

// Debug 调试日志,详尽的程序日志
// 调试用例:
// 			Logger.Debug("Database", zap.String("sql", sql))
func Debug(moduleName string, fields ...zap.Field) {
	Logger.Debug(moduleName, fields...)
}

// 告知类日志
func Info(moduleName string, fields ...zap.Field) {
	Logger.Info(moduleName, fields...)
}

// 警告类日志
func Warn(moduleName string, fields ...zap.Field) {
	Logger.Warn(moduleName, fields...)
}

// 错误记录,不应该中断程序,查看日志时重点关注
func Error(moduleName string, fields ...zap.Field) {
	Logger.Error(moduleName, fields...)
}

// 级别同 Error, 写完 log 后调用 os.Exit() 退出程序
func Fatal(moduleName string, fields ...zap.Field) {
	Logger.Fatal(moduleName, fields...)
}

// DebugString 记录一条字符串类型的 debug 日志
// 调试用例:
// 			logger.DebugString("SMS", "短信内容", string(result.RawResponse))
func DebugString(moduleName, name, msg string) {
	Logger.Debug(moduleName, zap.String(name, msg))
}

func InfoString(moduleName, name, msg string) {
	Logger.Info(moduleName, zap.String(name, msg))
}

func WarnString(moduleName, name, msg string) {
	Logger.Warn(moduleName, zap.String(name, msg))
}

func ErrorString(moduleName, name, msg string) {
	Logger.Error(moduleName, zap.String(name, msg))
}

func FatalString(moduleName, name, msg string) {
	Logger.Fatal(moduleName, zap.String(name, msg))
}

// DebugJSON 记录对象类型的 debug 日志,使用 json.Marshal 进行编码
// 调用实例:
// 			logger.DebugJSON("Auth", "读取登录用户", auth..CurrentUser())
func DebugJSON(moduleName, name string, value interface{}) {
	Logger.Debug(moduleName, zap.String(name, jsonString(value)))
}

func InfoJSON(moduleName, name string, value interface{}) {
	Logger.Info(moduleName, zap.String(name, jsonString(value)))
}

func WarnJSON(moduleName, name string, value interface{}) {
	Logger.Warn(moduleName, zap.String(name, jsonString(value)))
}

func ErrorJSON(moduleName, name string, value interface{}) {
	Logger.Error(moduleName, zap.String(name, jsonString(value)))
}

func FatalJSON(moduleName, name string, value interface{}) {
	Logger.Fatal(moduleName, zap.String(name, jsonString(value)))
}

func jsonString(value interface{}) string {
	b, err := json.Marshal(value)
	if err != nil {
		Logger.Error("Logger", zap.String("JSON marshal error", err.Error()))
	}
	return string(b)
}
