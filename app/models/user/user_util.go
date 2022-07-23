// 存放用户模型相关的数据库操作
// 可以直接使用 user. 调用的都存在此文件中
package user

import "gohub/pkg/database"

// 判断 Email 是否被注册
func IsEmailExist(email string) bool {
	var count int64
	database.DB.Model(User{}).Where("email = ?", email).Count(&count)
	return count > 0
}

// 判断手机号是否被注册
func IsPhoneExist(phone string) bool {
	var count int64
	database.DB.Model(User{}).Where("phone = ?", phone).Count(&count)
	return count > 0
}
