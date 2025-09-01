package util

import (
	"regexp"
	"strings"
)

var charReplacements = map[string]string{
	":":  "：",
	"<":  "＜",
	">":  "＞",
	"“":  "\"",
	"”":  "\"",
	"|":  "｜",
	"?":  "？",
	"*":  "＊",
	"\\": "＼",
	"/":  "／",
}

func SanitizeFileName(fn string) string {
	for k, v := range charReplacements {
		fn = strings.Replace(fn, k, v, -1)
	}
	return fn
}

func ExtractContentWithinTag(content string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(content, "")
}
