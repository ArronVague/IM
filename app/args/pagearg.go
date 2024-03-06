package args

import "time"

type PageArg struct {
	Pagefrom int       `json:"pagefrom" form:"pagefrom"` //从哪页开始
	Pagesize int       `json:"pagesize" form:"pagesize"` //每页大小
	Kword    string    `json:"kword" form:"kword"`       //关键词
	Asc      string    `json:"asc" form:"asc"`
	Desc     string    `json:"desc" form:"desc"`
	Name     string    `json:"name" form:"name"`
	Userid   int64     `json:"userid" form:"userid"`
	Dstid    int64     `json:"dstid" form:"dstid"`       //dstid
	Datefrom time.Time `json:"datafrom" form:"datafrom"` //时间点1
	Dateto   time.Time `json:"dateto" form:"dateto"`     //时间点2
	Total    int64     `json:"total" form:"total"`
}
