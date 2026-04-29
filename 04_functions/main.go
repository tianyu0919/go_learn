// ========================================
// Lesson 04: Functions
// ========================================
// Go 函数的各种用法

package main

import (
	"fmt"
	"strings"
)

// ---- 基本函数 ----
func greet(name string) {
	fmt.Printf("Hello, %s!\n", name)
}

// 带返回值的函数
func add(a, b int) int {
	return a + b
}

// 多返回值（Go 的特色功能！）
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return a / b, nil
}

// 命名返回值
func swap(a, b string) (first, second string) {
	first = b
	second = a
	return // 裸返回（naked return）
}

// ---- 可变参数 ----
func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// ---- 函数作为值 ----
func applyOperation(a, b int, op func(int, int) int) int {
	return op(a, b)
}

// ---- 闭包 (Closure) ----
func makeCounter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

// ---- defer 关键字 ----
func demoDefer() {
	fmt.Println("defer demo start")

	// defer 会在函数返回前执行，按 LIFO（后进先出）顺序
	defer fmt.Println("deferred: 1st")
	defer fmt.Println("deferred: 2nd")
	defer fmt.Println("deferred: 3rd")

	fmt.Println("defer demo end")
	// 输出顺序: start -> end -> 3rd -> 2nd -> 1st
}

// ---- init 函数（了解） ----
// init 函数在 main 之前自动执行，常用于初始化
func init() {
	fmt.Println("[init] Package initialized")
}

func main() {
	// 基本函数调用
	fmt.Println("--- Basic Functions ---")
	greet("Go Learner")
	fmt.Printf("3 + 5 = %d\n", add(3, 5))

	// 多返回值
	fmt.Println("\n--- Multiple Returns ---")
	result, err := divide(10, 3)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("10 / 3 = %.2f\n", result)
	}

	// 使用 _ 忽略不需要的返回值
	_, err2 := divide(10, 0)
	if err2 != nil {
		fmt.Println("Error:", err2)
	}

	// 命名返回值
	fmt.Println("\n--- Named Returns ---")
	a, b := swap("hello", "world")
	fmt.Printf("swapped: %s, %s\n", a, b)

	// 可变参数
	fmt.Println("\n--- Variadic Functions ---")
	fmt.Printf("sum(1,2,3) = %d\n", sum(1, 2, 3))
	fmt.Printf("sum(1,2,3,4,5) = %d\n", sum(1, 2, 3, 4, 5))

	// 展开切片传递给可变参数
	nums := []int{10, 20, 30}
	fmt.Printf("sum(slice...) = %d\n", sum(nums...))

	// 函数作为值
	fmt.Println("\n--- Functions as Values ---")
	multiply := func(a, b int) int { return a * b }
	fmt.Printf("applyOperation(3, 4, multiply) = %d\n",
		applyOperation(3, 4, multiply))

	// 匿名函数（立即执行）
	func() {
		fmt.Println("I'm an anonymous function!")
	}()

	// 闭包
	fmt.Println("\n--- Closures ---")
	counter := makeCounter()
	fmt.Printf("counter: %d\n", counter()) // 1
	fmt.Printf("counter: %d\n", counter()) // 2
	fmt.Printf("counter: %d\n", counter()) // 3

	// defer
	fmt.Println("\n--- Defer ---")
	demoDefer()

	// 实际用例：defer 常用于资源清理
	fmt.Println("\n--- Practical Example ---")
	processText("  Hello, Go World!  ")
}

// 实际用例：字符串处理函数
func processText(input string) {
	// 链式操作
	result := strings.TrimSpace(input)
	result = strings.ToUpper(result)
	words := strings.Split(result, " ")

	fmt.Printf("Original: '%s'\n", input)
	fmt.Printf("Processed: '%s'\n", result)
	fmt.Printf("Words: %v (count: %d)\n", words, len(words))
}

// ========================================
// 练习:
// 1. 写一个递归函数计算斐波那契数列 fibonacci(n int) int
// 2. 写一个高阶函数 filter(nums []int, predicate func(int) bool) []int
// 3. 用闭包实现一个简单的缓存（memoization）
// ========================================
