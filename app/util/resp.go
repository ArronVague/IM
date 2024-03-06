package util

import (
	"encoding/json"
	"log"
	"net/http"
)

// ResponseData 定义一个结构体
type ResponseData struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	//omitempty 标签表示如果该字段的值为空（零值），则在生成 JSON 时会忽略该字段。
	Data interface{} `json:"data,omitempty"`
}

type H struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data,omitempty"`
	Rows  interface{} `json:"rows,omitempty"`
	Total interface{} `json:"total,omitempty"`
}

// RespFail 失败的返回结果
func RespFail(writer http.ResponseWriter, msg string) {
	Resp(writer, -1, nil, msg)
}

// RespOk 返回成功
func RespOk(writer http.ResponseWriter, data interface{}, msg string) {
	Resp(writer, 0, data, msg)
}

// RespOkList 返回列表数据
func RespOkList(w http.ResponseWriter, lists interface{}, total interface{}) {
	RespList(w, 0, lists, total)
}

func RespList(w http.ResponseWriter, code int, data interface{}, total interface{}) {
	w.Header().Set("Content-Type", "application/json")
	//设置200状态
	w.WriteHeader(http.StatusOK)
	h := H{
		Code:  code,
		Rows:  data,
		Total: total,
	}
	//将结构体转化成JSOn字符串
	ret, err := json.Marshal(h)
	if err != nil {
		log.Println(err.Error())
	}
	//输出
	w.Write(ret)
}

func Resp(writer http.ResponseWriter, code int, data interface{}, msg string) {
	//设置header 为JSON 默认是test/html,所以特别指出返回的数据类型为application/json
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	rep := ResponseData{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	//将结构体转化为json字符串
	ret, err := json.Marshal(rep)
	if err != nil {
		//panic 是 Go 语言的一个内建函数，它会立即停止当前函数的执行，并开始回溯（unwinding）过程。在回溯过程中，运行时会执行所有的 defer 语句，然后返回到当前函数的调用者。这个过程会一直继续，直到回溯到当前的 Go 程（goroutine）的起点，然后程序会退出。比较激进，慎用。
		log.Panicln(err.Error())
	}

	//返回json ok
	writer.Write(ret)
}
