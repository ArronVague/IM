package controller

import (
	"IM/app/args"
	"IM/app/service"
	"IM/app/util"
	"net/http"
)

var contactService service.ContactService

// AddFriend 添加朋友
func AddFriend(writer http.ResponseWriter, request *http.Request) {
	var arg args.AddNewMember
	util.Bind(request, &arg)
	friend := contactService.SearchFriendByName(arg.DstName)
	if friend.Id == 0 {
		util.RespFail(writer, "您要添加的好友不存在")
	} else {
		err := contactService.AddFriend(arg.Userid, friend.Id)
		if err != nil {
			util.RespFail(writer, err.Error())
		} else {
			util.RespOk(writer, nil, "好友添加成功")
		}
	}
}
