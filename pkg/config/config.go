package config

import (
	"os"

	"github.com/spf13/cast"
	viperlib "github.com/spf13/viper" // 自定义包名，避免与内置 viper 实例冲突
)

// viper 实例
var viper *viperlib.Viper

// 动态加载配置信息
type ConfigFunc func() map[string]interface{}

// 先加载到此数组, loadConfig 再动态加载生成配置信息
var ConfigFuncs map[string]ConfigFunc

func init() {
	// 1. 初始化Viper库
	viper = viperlib.New()
	// 2. 设置配置类型: 支持 "json", "toml", "yaml", "yml", "properties","props", "prop", "env", "dotenv"
	viper.SetConfigType("env")
	// 3. 环境变量配置文件查找路径, 相对于 main.go
	viper.AddConfigPath(".")
	// 4. 设置环境变量的前缀,用于区分 Go 的系统环境变量
	viper.SetEnvPrefix("appenv")
	// 5. 读取环境变量 (支持 flags)
	viper.AutomaticEnv()

	ConfigFuncs = make(map[string]ConfigFunc)
}

// 初始化配置信息,完成对环境变量和config信息的加载
func InitConfig(env string) {
	// 1. 加载环境变量
	loadEnv(env)
	// 2. 注册配置信息
	loadConfig()
}

func loadConfig() {
	for name, fn := range ConfigFuncs {
		viper.Set(name, fn())
	}
}

func loadEnv(envSuffix string) {
	// 1. 默认加载 .env 文件, 如果有传参 --env=name, 则加载 .env.name 文件
	envPath := ".env"
	if len(envSuffix) > 0 {
		filePath := ".env." + envSuffix
		if _, err := os.Stat(filePath); err == nil {
			// 如 .env.testing、.env.stage
			envPath = filePath
		}
	}

	// 2. 加载 env
	viper.SetConfigName(envPath)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	// 3. 监控 .env 文件,变更时重新加载
	viper.WatchConfig()
}

// 读取环境变量,支持默认值
func Env(envName string, defaultValue ...interface{}) interface{} {
	if len(defaultValue) > 0 {
		return internalGet(envName, defaultValue[0])
	}
	return internalGet(envName)
}

func internalGet(envName string, defaultValue ...interface{}) interface{} {
	// config 或者 环境变量不存在的情况下
	if !viper.IsSet(envName) {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}
	return viper.Get(envName)
}

// 新增配置项
func AddEnv(name string, fn ConfigFunc) {
	ConfigFuncs[name] = fn
}

// 获取配置项
// 第一个参数 path 允许使用点获取,如: app.name
// 第二个参数允许传参默认值
func Get(path string, defaultValue ...interface{}) string {
	return GetString(path, defaultValue...)
}

// 获取 string 类型的配置信息
func GetString(path string, defaultValue ...interface{}) string {
	return cast.ToString(internalGet(path, defaultValue...))
}

// 获取 int 类型的配置信息
func GetInt(path string, defaultValue ...interface{}) int {
	return cast.ToInt(internalGet(path, defaultValue...))
}

// 获取 float64 类型的配置信息
func GetFloat64(path string, defaultValue ...interface{}) float64 {
	return cast.ToFloat64(internalGet(path, defaultValue...))
}

// 获取 int64 类型的配置信息
func GetInt64(path string, defaultValue ...interface{}) int64 {
	return cast.ToInt64(internalGet(path, defaultValue...))
}

// 获取 uint 类型的配置信息
func GetUint(path string, defaultValue ...interface{}) uint {
	return cast.ToUint(internalGet(path, defaultValue...))
}

// 获取 bool 类型的配置信息
func GetBool(path string, defaultValue ...interface{}) bool {
	return cast.ToBool(internalGet(path, defaultValue...))
}

// 获取结构类型的配置信息
func GetStringMapString(path string) map[string]string {
	return viper.GetStringMapString(path)
}
