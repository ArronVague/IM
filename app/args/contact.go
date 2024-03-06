package args

type ContactArg struct {
	PageArg
	Userid int64 `json:"userid" form:"userid"`
	Dstid  int64 `json:"dstid" form:"dstid"`
}

// AddNewMember 添加新的成员
type AddNewMember struct {
	PageArg
	Userid  int64  `json:"userid" form:"userid"`
	DstName string `json:"dstname" form:"dstname"`
}
