package model

import "time"

// Contact 好友和群都存在这个表里面
// 可根据具体业务做拆分
type Contact struct {
	Id       int64     `xorm:"pk autoincr bigint(20)" form:"id" json:"id"`
	Ownerid  int64     `xorm:"bigint(20)" form:"ownerid" json:"ownerid"` // 记录是谁的
	Dstobj   int64     `xorm:"bigint(20)" form:"dstobj" json:"dstobj"`   // 对端信息（用户id，群id // ）
	Cate     int       `xorm:"int(11)" form:"cate" json:"cate"`          // 什么类型
	Memo     string    `xorm:"varchar(120)" form:"memo" json:"memo"`     // 备注
	Createat time.Time `xorm:"datetime" form:"createat" json:"createat"` // 创建时间
}

const (
	ContactCateUser      = 0x01 //个人
	ContactCateCommunity = 0x02 //群
)
