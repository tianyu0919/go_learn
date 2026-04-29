// Package calculator 提供基本的数学运算
package calculator

import "fmt"

// Add 返回两个整数的和
func Add(a, b int) int {
	return a + b
}

// Subtract 返回两个整数的差
func Subtract(a, b int) int {
	return a - b
}

// Multiply 返回两个整数的积
func Multiply(a, b int) int {
	return a * b
}

// Divide 返回两个浮点数的商，除数为零时返回错误
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("cannot divide by zero")
	}
	return a / b, nil
}

// 未导出函数（小写开头），只能在包内使用
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
