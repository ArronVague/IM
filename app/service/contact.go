package service

import (
	"IM/app/model"
	"errors"
	"time"
)

type ContactService struct{}

// AddFriend 添加好友
func (service *ContactService) AddFriend(userid int64, dstid int64) error {
	if dstid == userid {
		return errors.New("不能添加自己为好友啊")
	}
	//判断是否已经添加了好友
	friend := model.Contact{}
	model.DbEngine.Where("ownerid = ?", userid).And("dstobj = ?", dstid).And("cate = ?", model.ContactCateUser).Get(&friend)
	//如果好友已经存在，则不添加
	if friend.Id > 0 {
		return errors.New("该好友已经添加过了")
	}
	//开启事务
	session := model.DbEngine.NewSession()
	session.Begin()
	//插入两条好友关系数据
	_, s1 := session.InsertOne(model.Contact{
		Ownerid:  userid,
		Dstobj:   dstid,
		Cate:     model.ContactCateUser,
		Createat: time.Now(),
	})
	_, s2 := session.InsertOne(model.Contact{
		Ownerid:  dstid,
		Dstobj:   userid,
		Cate:     model.ContactCateUser,
		Createat: time.Now(),
	})
	if s1 == nil && s2 == nil {
		session.Commit()
		return nil
	} else {
		session.Rollback()
		if s1 != nil {
			return s1
		}
		return s2
	}
}

// SearchFriendByName 根据姓名搜索用户（看样子其实是手机号）
func (service *ContactService) SearchFriendByName(mobile string) model.User {
	user := model.User{}
	model.DbEngine.Where("mobile = ?", mobile).Get(&user)
	return user
}

// SearchCommunityIds 获取用户全部群ID
func (service *ContactService) SearchCommunityIds(userId int64) (comIds []int64) {
	contacts := make([]model.Contact, 0)
	comIds = make([]int64, 0)

	model.DbEngine.Where("ownerid = ? and cate = ?", userId, model.ContactCateCommunity).Find(&contacts)
	for _, v := range contacts {
		comIds = append(comIds, v.Dstobj)
	}
	return comIds
}
