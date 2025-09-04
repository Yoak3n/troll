package main

import (
	"fmt"
	"math/rand"
)

// CreateLink 创建带超链接的终端文本
func CreateLink(text, url string) string {
	return fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", url, text)
}

// CreateStyledLink 创建带样式的超链接
func CreateStyledLink(text, url string, colorCode int) string {
	return fmt.Sprintf("\033[%dm\033]8;;%s\033\\%s\033]8;;\033\\\033[0m", colorCode, url, text)
}

const alpha = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateTextId() string {
	b := make([]byte, 8)
	for i := range b {
		b[i] = alpha[rand.Intn(len(alpha))]
	}
	return fmt.Sprintf("%s%d", b, rand.Intn(1000000000))
}
