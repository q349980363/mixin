package main

import (
	"time"

	"github.com/jinzhu/gorm"
)

//用户信息
type UserInfo struct {
	gorm.Model
	UserName string
	PassWord string
	ToKen    string
	Tags     string
}

//基础消息
type Chat struct {
	gorm.Model
	Type     string
	Data     string
	UserName string
	Ip       string
	IpAddr   string
}

//系统消息       [加好用申请|加群申请] 等
type SystemChat struct {
	gorm.Model
	Type     string
	Txt      string //明文显示
	Data     string
	ExInt    int
	ExString string
	UserName string `gorm:"index"` //所属用户名
	Result   string //用户操作结果
	Read     bool   //消息是否已读  已读做到关系上还是做到消息上?

	Operation string //是否可操作以及操作类型[""|"ok"|"ok_no"]
}

//基础消息
type FriendsChat struct {
	gorm.Model
	Type     string //消息类型
	Data     string
	UserName string `gorm:"index"` //所属用户名
	Target   string `gorm:"index"` //目标关联用户名
	Ip       string //发送者IP
	IpAddr   string //发送者IP物理位置定位
	Read     bool   //消息是否已读  已读做到关系上还是做到消息上?
}

//基础消息
type GroupChat struct {
	gorm.Model
	Type     string //消息类型
	Data     string
	UserName string `gorm:"index"` //所属用户名
	Target   string `gorm:"index"` //目标关联用户名
	Ip       string //发送者IP
	IpAddr   string //发送者IP物理位置定位
	Read     bool   //消息是否已读  已读做到关系上还是做到消息上?
}

//好友关系
type Friends struct {
	gorm.Model
	UserName   string    `gorm:"index"` //所属用户名
	Target     string    `gorm:"index"` //目标关联用户名
	Unread     int       //消息是否已读  已读做到关系上还是做到消息上?
	LastChatAt time.Time //消息是否已读  已读做到关系上还是做到消息上?
	LastChat   string
	Path       string //特殊参数
}

//组
type Group struct {
	gorm.Model
	Name     string
	UserName string `gorm:"index"` //所属用户名
}

//组关系
type GroupRelation struct {
	gorm.Model
	UserName string `gorm:"index"` //所属用户名
	Target   string `gorm:"index"` //目标关联群名称

	Unread     int       //消息是否已读  已读做到关系上还是做到消息上?
	LastChatAt time.Time //消息是否已读  已读做到关系上还是做到消息上?
	LastChat   string
}
