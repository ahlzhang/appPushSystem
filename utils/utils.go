package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"time"
	"strings"
	"path/filepath"
)

//MD5加密
func MyMd5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
	}
	dirctory := strings.Replace(dir, "\\", "/", 0)
	return dirctory
}

//时间换成字符串。格式：2016-01-02 15:04:05
func TimeToString(timeby time.Time) string {

	return timeby.Format("2006-01-02 15:04:05")
}