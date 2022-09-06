package job

import (
	"fmt"
	"github.com/zhiting-tech/smartassistant/modules/config"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	lifeTime   = 3600 * 24 * 7
	timeLayout = "20060102"
	logPrefix  = "smartassistant"
)

// RemoveLogTask 定时删除超过七天的日志
func RemoveLogTask() {
	// 当前运行路径下log文件夹
	path := fmt.Sprintf("%s/log", config.GetConf().SmartAssistant.RuntimePath)
	// 获取文件夹下的文件追加进arr切片
	filepath.Walk(path, removeLog)
}

// 字符串以.号切割
func split(s rune) bool {
	if s == '.' {
		return true
	}
	return false
}

func removeLog(path string, info fs.FileInfo, err error) error {
	str := strings.FieldsFunc(path, split)

	// 前缀名
	str[0] = strings.Replace(str[0], fmt.Sprintf("%s/log/", config.GetConf().SmartAssistant.RuntimePath), "", 1)

	if isExpired(str) {
		os.Remove(path)
	}

	return err
}

func isExpired(str []string) bool {
	lastDay := time.Unix(time.Now().Unix()-(lifeTime), 0).Format(timeLayout)

	// 判断文件前缀名,过期时间,是否带log后缀
	if len(str) >= 3 && str[0] == logPrefix && str[1] <= lastDay && str[2] == "log" {
		return true
	}

	return false
}
