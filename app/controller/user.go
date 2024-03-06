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

func UserLogin(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()

	mobile := request.PostForm.Get("mobile")
	plainPwd := request.PostForm.Get("passwd")

	if len(mobile) == 0 || len(plainPwd) == 0 {
		util.RespFail(writer, "用户名或密码不正确")
	}

	loginUser, err := UserService.Login(mobile, plainPwd)
	if err != nil {
		util.RespFail(writer, err.Error())
	} else {
		util.RespOk(writer, loginUser, "")
	}
}
