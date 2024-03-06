package service

import "IM/app/model"

type ContactService struct{}

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
