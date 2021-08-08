package database

import (
	"time"
)

// create database kim_base default character set utf8mb4 collate utf8mb4_unicode_ci;
// create database kim_message default character set utf8mb4 collate utf8mb4_unicode_ci;

type Model struct {
	ID        int64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type MessageIndex struct {
	ID        int64  `gorm:"primarykey"`
	AccountA  string `gorm:"index;size:60;not null;comment:队列唯一标识"`
	AccountB  string `gorm:"size:60;not null;comment:另一方"`
	Direction byte   `gorm:"default:0;not null;comment:1表示AccountA为发送者"`
	MessageID int64  `gorm:"not null;comment:关联消息内容表中的ID"`
	Group     string `gorm:"size:30;comment:群ID，单聊情况为空"`
	SendTime  int64  `gorm:"index;not null;comment:消息发送时间"`
}

type MessageContent struct {
	ID       int64  `gorm:"primarykey"`
	Type     byte   `gorm:"default:0"`
	Body     string `gorm:"size:5000;not null"`
	Extra    string `gorm:"size:500"`
	SendTime int64  `gorm:"index"`
}

type User struct {
	Model
	App      string `gorm:"size:30"`
	Account  string `gorm:"uniqueIndex;size:60"`
	Password string `gorm:"size:30"`
	Avatar   string `gorm:"size:200"`
	Nickname string `gorm:"size:20"`
}

type Group struct {
	Model
	Group        string `gorm:"uniqueIndex;size:30"`
	App          string `gorm:"size:30"`
	Name         string `gorm:"size:50"`
	Owner        string `gorm:"size:60"`
	Avatar       string `gorm:"size:200"`
	Introduction string `gorm:"size:300"`
}

// GroupMember GroupMember
type GroupMember struct {
	Model
	Account string `gorm:"uniqueIndex:uni_gp_acc;size:60"`
	Group   string `gorm:"uniqueIndex:uni_gp_acc;index;size:30"`
	Alias   string `gorm:"size:30"`
}
