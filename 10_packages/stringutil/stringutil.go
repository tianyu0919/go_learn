// Package stringutil 提供字符串工具函数
package stringutil

import "strings"

// Reverse 反转字符串
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// IsPalindrome 判断字符串是否是回文
func IsPalindrome(s string) bool {
	s = strings.ToLower(s)
	return s == Reverse(s)
}

// Capitalize 将每个单词的首字母大写
func Capitalize(s string) string {
	words := strings.Fields(s)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + word[1:]
		}
	}
	return strings.Join(words, " ")
}
