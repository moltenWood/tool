package goToolsUtils

import (
	"bufio"
	"fmt"
	"os"
)

// 判断文件是否存在
func IsExist(fileAddr string) bool {
	// 读取文件信息，判断文件是否存在
	_, err := os.Stat(fileAddr)
	if err != nil {
		if os.IsExist(err) { // 根据错误类型进行判断
			return true
		}
		return false
	}
	return true
}

func WriteLog(filename string, content string) {
	//contentJson, err := json.Marshal(content)
	//if err != nil {
	//	fmt.Println("json.Marshal")
	//	return
	//}
	var fp *os.File
	var err1 error

	defer fp.Close()

	if IsExist(filename) != true {
		fp, err1 = os.Create(filename)
	} else {
		fp, err1 = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0666)
	}
	if err1 != nil {
		fmt.Println("文件操作失败", err1)
	}

	writer := bufio.NewWriter(fp)
	_, err := writer.WriteString(GetCurrentTime() +" -> "+ content + "\n")
	fmt.Println(GetCurrentTime() +" -> "+ content + "\n")
	if err != nil {
		fmt.Println("write error:", err)
	} else {
		//fmt.Println("write success:", n)
	}
	writer.Flush()
}
