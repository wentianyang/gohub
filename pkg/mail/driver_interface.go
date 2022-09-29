// 邮箱辅助包
package mail

type Driver interface {
	// 检查验证码
	Send(email Email, config map[string]string) bool
}
