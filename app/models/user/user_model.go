// 存放用户 model 定义、以及对象操作的逻辑代码
package user

import "gohub/app/models"

// 用户模型
// 不希望讲敏感信息输出给用户,所以 Email、Phone、Password 后面设置了 json:"-"
// 表示 JSON 解析器忽略字段
// 后面接口返回用户数据时候,这三个字短都会被隐藏
type User struct {
	models.BaseModel

	Name     string `json:"name,omitempty"`
	Email    string `json:"-"`
	Phone    string `json:"-"`
	Password string `json:"-"`

	models.CommonTimestampsField
}
