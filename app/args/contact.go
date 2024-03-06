package args

// AddNewMember 添加新的成员
type AddNewMember struct {
	PageArg
	Userid  int64  `json:"userid" form:"userid"`
	DstName string `json:"dstname" form:"dstname"`
}
