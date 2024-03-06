package controller

import (
	"IM/app/model"
	"IM/app/service"
	"IM/app/util"
	"net/http"
)

var UserService service.UserService

// UserRegister 用户注册
func UserRegister(writer http.ResponseWriter, request *http.Request) {
	var user model.User
	util.Bind(request, &user)
	user, err := UserService.UserRegister(user.Mobile, user.Passwd, user.Nickname, user.Avatar, user.Sex)
	if err != nil {
		util.RespFail(writer, err.Error())
	} else {
		util.RespOk(writer, user, "")
	}
}
