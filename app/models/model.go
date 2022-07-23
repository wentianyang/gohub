// 模型通用属性和方法
package models

import "time"

// BaseModel 模型基类
type BaseModel struct {
	ID uint64 `gorm:"column:id;primaryKey;autoIncrement;" json:"id,omitempty"`
}

// 时间戳模型
type CommonTimestampsField struct {
	CreateTime  time.Time `gorm:"column:created_at;index;" json:"created_at,omitempty"`
	UpdatedTime time.Time `gorm:"column:updated_at;index;" json:"updated_at,omitempty"`
}
