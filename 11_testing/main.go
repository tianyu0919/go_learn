// ========================================
// Lesson 11: Testing
// ========================================
// Go 内置了强大的测试框架，不需要第三方库
// 测试文件以 _test.go 结尾
// 运行: go test ./11_testing/

package main

import "fmt"

// ---- 被测试的函数 ----

// Add 返回两个整数的和
func Add(a, b int) int {
	return a + b
}

// Max 返回最大值
func Max(nums ...int) (int, error) {
	if len(nums) == 0 {
		return 0, fmt.Errorf("empty slice")
	}
	max := nums[0]
	for _, n := range nums[1:] {
		if n > max {
			max = n
		}
	}
	return max, nil
}

// FizzBuzz 返回 FizzBuzz 结果
func FizzBuzz(n int) string {
	switch {
	case n%15 == 0:
		return "FizzBuzz"
	case n%3 == 0:
		return "Fizz"
	case n%5 == 0:
		return "Buzz"
	default:
		return fmt.Sprintf("%d", n)
	}
}

// Fibonacci 返回第 n 个斐波那契数
func Fibonacci(n int) int {
	if n <= 0 {
		return 0
	}
	if n == 1 {
		return 1
	}
	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

func main() {
	fmt.Println("Run tests with: go test -v ./11_testing/")
	fmt.Println()
	fmt.Println("Test commands:")
	fmt.Println("  go test ./11_testing/           # Run tests")
	fmt.Println("  go test -v ./11_testing/        # Verbose output")
	fmt.Println("  go test -run TestAdd ./11_testing/  # Run specific test")
	fmt.Println("  go test -cover ./11_testing/    # Show coverage")
	fmt.Println("  go test -bench=. ./11_testing/  # Run benchmarks")
}
