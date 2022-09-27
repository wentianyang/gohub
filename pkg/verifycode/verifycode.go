// Package verify 用以发送手机验证码和邮箱验证码
package verifycode

import (
	"gohub/pkg/app"
	"gohub/pkg/config"
	"gohub/pkg/helpers"
	"gohub/pkg/logger"
	"gohub/pkg/redis"
	"gohub/pkg/sms"
	"strings"
	"sync"
)

type VerifyCode struct {
	Store Store
}

var once sync.Once
var internalVerifyCode *VerifyCode

func NewVerifyCode() *VerifyCode {
	once.Do(func() {
		internalVerifyCode = &VerifyCode{
			Store: &RedisStore{
				RedisClient: redis.Redis,
				// 增加前缀保持数据库整洁, 出问题时方便调试
				KeyPrefix: config.GetString("app.name") + ":verifycode",
			},
		}
	})

	return internalVerifyCode
}

// SendSMS 发送短信验证码, 调试实例:
// 		verifycode.NewVerifyCode().SendSMS(request.Phone)
func (vc *VerifyCode) SendSms(phone string) bool {
	// 生成验证码
	code := vc.generateVerifyCode(phone)

	// 方便本地和 API 自动测试
	if !app.IsProduction() && strings.HasPrefix(phone, config.GetString("verifycode.debug_phone_prefix")) {
		return true
	}

	// 发送短信
	return sms.NewSMS().Send(phone, sms.Message{
		Template: config.GetString("sms.aliyun.template_code"),
		Data:     map[string]string{"code": code},
	})
}

// CheckAnswer 检查用户提交的验证码是否正确, key 可以是手机号 或者 email
func (vc *VerifyCode) CheckAnswer(key string, answer string) bool {
	logger.DebugJSON("验证码", "检查验证码", map[string]string{key: answer})

	// 方便开发, 在非生产环境下, 具备特殊前缀的手机号和 email 后缀, 会直接验证成功
	if !app.IsProduction() && (strings.HasSuffix(key, config.GetString("verifycode.debug_email_suffix")) ||
		strings.HasPrefix(key, config.GetString("verifycode.debug_phone_prefix"))) {
		return true
	}
	return vc.Store.Verify(key, answer, false)
}

func (vc *VerifyCode) generateVerifyCode(key string) string {

	// 生成随机验证码
	code := helpers.RandomNumber(config.GetInt("verifycode.code_length"))

	// 为了方便开发
	if app.IsLocal() {
		code = config.GetString("verifycode.debug_code")
	}

	logger.DebugJSON("验证码", "生成验证码", map[string]string{key: code})

	vc.Store.Set(key, code)
	return code
}
