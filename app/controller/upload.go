package controller

import (
	"IM/app/util"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func init() {
	os.MkdirAll("./resource", os.ModePerm)
}

func FileUpload(writer http.ResponseWriter, request *http.Request) {
	UploadLocal(writer, request)
}

// UploadLocal 将文件存储在本地/resource目录下
func UploadLocal(writer http.ResponseWriter, request *http.Request) {
	//获得上传源文件
	srcFile, head, err := request.FormFile("file")
	if err != nil {
		util.RespFail(writer, err.Error())
	}
	//创建一个新的文件
	suffix := ".png"
	srcFilename := head.Filename
	splitMsg := strings.Split(srcFilename, ".")
	if len(splitMsg) > 1 {
		suffix = "." + splitMsg[len(splitMsg)-1]
	}
	filetype := request.FormValue("filetype")
	if len(filetype) > 0 {
		suffix = filetype
	}
	filename := fmt.Sprintf("%d%s%s", time.Now().Unix(), util.GenRandomStr(32), suffix)
	//创建文件
	filepath := "./resource/" + filename
	dstfile, err := os.Create(filepath)
	if err != nil {
		util.RespFail(writer, err.Error())
		return
	}
	//将源文件拷贝到新文件。这是因为在处理 HTTP 文件上传时，上传的文件通常是以数据流的形式存在的，而不是一个实体的文件。当你在浏览器中选择一个文件并点击上传时，你的浏览器会将这个文件的内容读取到内存中，然后通过 HTTP 请求将这些内容发送到服务器。在服务器端，这些内容会被封装成一个临时的文件对象（在 Go 语言中，这个对象通常是一个 `*multipart.File` 类型的对象）。
	//
	//因此，当服务器接收到这个上传的文件时，它实际上是在接收一个包含了文件内容的数据流，而不是一个实体的文件。服务器需要做的是将这个数据流保存到一个新的文件中，这就是为什么需要将上传的文件的内容复制到新的文件中。
	//
	//另外，这种方式还有以下几点好处：
	//
	//1. 可以在保存文件时对文件名进行修改，例如添加时间戳或随机字符串，以避免文件名冲突。
	//2. 可以在保存文件时对文件内容进行处理，例如进行压缩或加密。
	//3. 可以控制文件保存的位置，例如保存到特定的目录或者分布式文件系统中。
	_, err = io.Copy(dstfile, srcFile)
	if err != nil {
		util.RespFail(writer, err.Error())
		return
	}

	util.RespOk(writer, filepath, "")
}
