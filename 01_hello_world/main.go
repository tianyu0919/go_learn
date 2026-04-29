// ========================================
// Lesson 01: Hello World
// ========================================
// Go 程序的基本结构
// 每个 Go 文件都以 package 声明开头
// main 包是程序的入口包

package main

// import 用于导入其他包
// fmt 是 Go 标准库中的格式化 I/O 包
import "fmt"

// main 函数是程序的入口点，程序从这里开始执行
func main() {
	// Println 会打印一行文本并换行
	fmt.Println("Hello, World!")

	// Printf 支持格式化输出，类似 C 语言
	name := "Go Learner" // := 是短变量声明，自动推断类型
	fmt.Printf("Welcome, %s!\n", name)
	greet(name)

	// 多行打印
	fmt.Println("Let's start learning Go!")
	fmt.Println("Go version: 1.25")
}

func greet(name string) {
	fmt.Printf("Hello, %s!\n", name)
}

// ========================================
// 练习:
// 1. 修改 name 变量为你自己的名字
// 2. 使用 fmt.Printf 打印你的年龄: fmt.Printf("I am %d years old\n", age)
// 3. 尝试导入 "time" 包并打印当前时间: time.Now()
// ========================================
