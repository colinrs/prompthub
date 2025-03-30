// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"

	"gorm.io/gorm"
)

const TableNameUsersSave = "users_save"

// UsersSave mapped from table <users_save>
type UsersSave struct {
	ID        int32          `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	UserID    int32          `gorm:"column:user_id;not null" json:"user_id"`
	PromptsID int32          `gorm:"column:prompts_id;not null" json:"prompts_id"`
	CreatedAt time.Time      `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

// TableName UsersSave's table name
func (*UsersSave) TableName() string {
	return TableNameUsersSave
}
