package config

import (
	"gohub/pkg/config"
)

func init() {
	config.AddEnv("app", func() map[string]interface{} {
		return map[string]interface{}{
			// 应用名称
			"name": config.Env("APP_NAME", "gohub"),
			// 当前环境,用以区分多环境, 一般为 local、stage、production、test
			"env": config.Env("APP_ENV", "local"),
			// 是否进入调试模式
			"debug": config.Env("APP_DEBUG", false),
			// 应用服务端口
			"port": config.Env("APP_PORT", "3000"),
			// 加密会话、JWT 加密
			"key": config.Env("APP_KEY", "33446a9dcf9ea060a0a6532b166da32f304af0d"),
			// 用以生成链接
			"uel": config.Env("APP_URL", "http://localhost:3000"),
			// 设置时区, JWT 里会使用,日志记录里也会使用
			"timezone": config.Env("APP_TIMEZONE", "Asia/Shanghai"),
		}
	})
}
