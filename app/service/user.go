package service

import (
	"IM/app/model"
	"IM/app/util"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type UserService struct{}

// UserRegister 用户注册
func (s *UserService) UserRegister(mobile, plainPwd, nickname, avatar, sex string) (user model.User, err error) {
	registerUser := model.User{}
	_, err = model.DbEngine.Where("mobile=? ", mobile).Get(&registerUser)
	if err != nil {
		return registerUser, err
	}
	//如果用户已经注册,返回错误信息
	if registerUser.Id > 0 {
		return registerUser, errors.New("该手机号已注册")
	}

	registerUser.Mobile = mobile
	registerUser.Avatar = avatar
	registerUser.Nickname = nickname
	registerUser.Sex = sex
	registerUser.Salt = fmt.Sprintf("%06d", rand.Int31n(10000))
	registerUser.Passwd = util.MakePasswd(plainPwd, registerUser.Salt)
	registerUser.Createat = time.Now()
	//插入用户信息
	_, err = model.DbEngine.InsertOne(&registerUser)

	return registerUser, err
}
