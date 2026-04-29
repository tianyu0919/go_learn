// ========================================
// Lesson 10: Packages & Modules
// ========================================
// Go 通过 package 组织代码，通过 module 管理依赖
//
// 本课程展示如何组织多包项目
// 包结构:
//   go-learning/
//   ├── go.mod
//   └── 10_packages/
//       ├── main.go          (当前文件)
//       └── calculator/
//           └── calculator.go (子包)

package main

import (
	"fmt"

	// 导入本项目的子包
	"go-learning/10_packages/calculator"
	"go-learning/10_packages/stringutil"
)

func main() {
	fmt.Println("--- Using calculator package ---")

	// 使用导出的函数（首字母大写 = public）
	fmt.Printf("Add(3, 5) = %d\n", calculator.Add(3, 5))
	fmt.Printf("Subtract(10, 4) = %d\n", calculator.Subtract(10, 4))
	fmt.Printf("Multiply(3, 7) = %d\n", calculator.Multiply(3, 7))

	result, err := calculator.Divide(10, 3)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Divide(10, 3) = %.2f\n", result)
	}

	_, err = calculator.Divide(10, 0)
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("\n--- Using stringutil package ---")
	fmt.Println(stringutil.Reverse("Hello, Go!"))
	fmt.Println(stringutil.IsPalindrome("racecar"))
	fmt.Println(stringutil.IsPalindrome("hello"))
	fmt.Println(stringutil.Capitalize("hello world go"))

	fmt.Println("\n--- Key Concepts ---")
	fmt.Println("1. Exported names start with uppercase (Add, Subtract)")
	fmt.Println("2. Unexported names start with lowercase (internal use only)")
	fmt.Println("3. Package name = directory name (convention)")
	fmt.Println("4. One package per directory")
	fmt.Println("5. go.mod defines the module path and dependencies")
}

// ========================================
// 练习:
// 1. 添加一个 Power(base, exp) 函数到 calculator 包
// 2. 创建一个新的 validator 包，包含 IsEmail, IsURL 等函数
// 3. 尝试使用 go get 安装一个第三方包（如 github.com/fatih/color）
// ========================================
