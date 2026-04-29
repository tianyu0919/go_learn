// ========================================
// Lesson 03: Control Flow
// ========================================
// Go 的控制流：if, for, switch

package main

import "fmt"

func main() {
	// ---- if/else ----
	fmt.Println("--- if/else ---")
	score := 85

	if score >= 90 {
		fmt.Println("Grade: A")
	} else if score >= 80 {
		fmt.Println("Grade: B")
	} else if score >= 70 {
		fmt.Println("Grade: C")
	} else {
		fmt.Println("Grade: D")
	}

	// if 支持初始化语句（变量作用域仅在 if 块内）
	if num := 42; num%2 == 0 {
		fmt.Printf("%d is even\n", num)
	} else {
		fmt.Printf("%d is odd\n", num)
	}
	// num 在这里不可访问

	// ---- for 循环 ----
	fmt.Println("\n--- for loop ---")

	// 标准 for 循环（类似 C/Java）
	for i := 0; i < 5; i++ {
		fmt.Printf("%d ", i)
	}
	fmt.Println()

	// while 风格（Go 没有 while 关键字，用 for 代替）
	count := 0
	for count < 3 {
		fmt.Printf("count = %d\n", count)
		count++
	}

	// 无限循环（需要 break 退出）
	sum := 0
	for {
		sum++
		if sum >= 10 {
			break
		}
	}
	fmt.Printf("sum = %d\n", sum)

	// for range 遍历（后面学集合时会大量使用）
	message := "Hello"
	for i, ch := range message {
		fmt.Printf("index %d: %c\n", i, ch)
	}

	// continue 跳过当前迭代
	fmt.Print("odd numbers: ")
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			continue
		}
		fmt.Printf("%d ", i)
	}
	fmt.Println()

	// ---- switch ----
	fmt.Println("\n--- switch ---")

	day := "Tuesday"
	switch day {
	case "Monday":
		fmt.Println("Start of work week")
	case "Tuesday", "Wednesday", "Thursday": // 多值匹配
		fmt.Println("Midweek")
	case "Friday":
		fmt.Println("TGIF!")
	default:
		fmt.Println("Weekend!")
	}

	// 无条件 switch（替代 if-else 链）
	temperature := 35
	switch {
	case temperature >= 40:
		fmt.Println("Extremely hot!")
	case temperature >= 30:
		fmt.Println("Hot")
	case temperature >= 20:
		fmt.Println("Comfortable")
	case temperature >= 10:
		fmt.Println("Cool")
	default:
		fmt.Println("Cold")
	}

	// 注意：Go 的 switch 默认不会 fall through
	// 如果需要 fall through，使用 fallthrough 关键字
	num := 1
	switch num {
	case 1:
		fmt.Println("One")
		fallthrough // 会继续执行下一个 case
	case 2:
		fmt.Println("Two (via fallthrough)")
	case 3:
		fmt.Println("Three")
	}

	// ---- 标签和 goto（了解即可，一般不推荐使用）----
	fmt.Println("\n--- Labels ---")
outer:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if i == 1 && j == 1 {
				fmt.Println("Breaking outer loop at i=1, j=1")
				break outer // 跳出外层循环
			}
			fmt.Printf("(%d,%d) ", i, j)
		}
	}
	fmt.Println()
}

// ========================================
// 练习:
// 1. 写一个 FizzBuzz: 打印 1-100，3 的倍数打印 Fizz，5 的倍数打印 Buzz，
//    同时是 3 和 5 的倍数打印 FizzBuzz
// 2. 用 for 循环打印九九乘法表
// 3. 用 switch 实现一个简单的计算器（+, -, *, /）
// ========================================
