package strings_utils

import (
	"github.com/mozillazg/go-pinyin"
)

// GetInitialAscii 获取首字母Ascii码
func GetInitialAscii(s string) int {
	ascii := []rune(s)
	if len(ascii) == 0 {
		return 0
	}
	return int(ascii[0])
}

// GetInitialPinyin 获取拼音首字母
func GetInitialPinyin(s string) string {
	py := pinyin.NewArgs()
	py.Style = pinyin.FirstLetter
	py.Fallback = func(r rune, a pinyin.Args) []string {
		return []string{string(r)}
	}

	p := pinyin.Pinyin(s, py) // 获取拼音
	if len(p) == 0 || len(p[0]) == 0 {
		return ""
	}
	return p[0][0]
}

