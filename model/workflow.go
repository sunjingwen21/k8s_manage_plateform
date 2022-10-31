package model

import "time"

/**
 * @author 王子龙
 * 时间：2022/9/28 14:30
 */
//定义结构体，属性与mysql表字段对齐
type Workflow struct {
	//gorm:"primaryKey"用于声明主键
	ID        uint       `json:"id" gorm:"primaryKey"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`

	Name       string `json:"name"`
	Namespace  string `json:"namespace"`
	Replicas   int32  `json:"replicas"`
	Deployment string `json:"deployment"`
	Service    string `json:"service"`
	Ingress    string `json:"ingress"`
	Type       string `json:"type" gorm:"column:type"`
	//Type: clusterip nodeport ingrress
}

//定义TableName方法，返回mysql表名，以此来定义mysql中的表名
func (*Workflow) TableName() string {
	return "workflow"
}
