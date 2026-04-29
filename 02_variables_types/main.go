// ========================================
// Lesson 02: Variables & Types
// ========================================
// Go 是静态类型语言，每个变量都有确定的类型

package main

import "fmt"

// 包级变量声明（在函数外部）
var globalVar string = "I am global"

func main() {
	// ---- 变量声明方式 ----

	// 方式1: var 关键字 + 类型
	var age int = 25
	var name string = "Alice"
	var isStudent bool = true

	// 方式2: 类型推断（省略类型，编译器自动推断）
	var city = "Beijing"

	// 方式3: 短变量声明（最常用，只能在函数内使用）
	score := 95.5 // float64
	count := 10   // int

	fmt.Println("--- Basic Variables ---")
	fmt.Printf("name: %s, age: %d, student: %t\n", name, age, isStudent)
	fmt.Printf("city: %s, score: %.1f, count: %d\n", city, score, count)

	// ---- 基本类型 ----
	fmt.Println("\n--- Basic Types ---")

	// 整数类型
	var i8 int8 = 127        // -128 to 127
	var i16 int16 = 32767    // -32768 to 32767
	var i32 int32 = 100      // also called rune
	var i64 int64 = 99999999 // 64位整数
	var u8 uint8 = 255       // 0 to 255, also called byte

	fmt.Printf("int8: %d, int16: %d, int32: %d, int64: %d, uint8: %d\n",
		i8, i16, i32, i64, u8)

	// 浮点类型
	var f32 float32 = 3.14
	var f64 float64 = 3.141592653589793
	fmt.Printf("float32: %f, float64: %.15f\n", f32, f64)

	// 字符串和字符
	var ch byte = 'A'          // byte 是 uint8 的别名
	var r rune = '中'           // rune 是 int32 的别名，用于 Unicode
	var str string = "Hello 世界" // 字符串是 UTF-8 编码
	fmt.Printf("char: %c, rune: %c, string: %s\n", ch, r, str)

	// ---- 零值 (Zero Values) ----
	fmt.Println("\n--- Zero Values ---")
	// Go 中未初始化的变量会被赋予零值
	var zeroInt int
	var zeroFloat float64
	var zeroString string
	var zeroBool bool
	fmt.Printf("int: %d, float: %f, string: '%s', bool: %t\n",
		zeroInt, zeroFloat, zeroString, zeroBool)

	// ---- 类型转换 ----
	fmt.Println("\n--- Type Conversion ---")
	// Go 不允许隐式类型转换，必须显式转换
	var intVal int = 42
	var floatVal float64 = float64(intVal)    // int -> float64
	var smallInt int32 = int32(intVal)         // int -> int32
	var strFromInt string = fmt.Sprintf("%d", intVal) // int -> string
	fmt.Printf("int: %d -> float64: %f, int32: %d, string: '%s'\n",
		intVal, floatVal, smallInt, strFromInt)

	// ---- 常量 ----
	fmt.Println("\n--- Constants ---")
	const pi = 3.14159
	const (
		statusOK    = 200
		statusError = 500
	)
	fmt.Printf("Pi: %f, OK: %d, Error: %d\n", pi, statusOK, statusError)

	// iota: 常量生成器，从 0 开始自增
	const (
		Sunday    = iota // 0
		Monday           // 1
		Tuesday          // 2
		Wednesday        // 3
	)
	fmt.Printf("Sunday=%d, Monday=%d, Tuesday=%d, Wednesday=%d\n",
		Sunday, Monday, Tuesday, Wednesday)
}

// ========================================
// 练习:
// 1. 声明一个表示温度的变量，将摄氏度转换为华氏度 (F = C*9/5 + 32)
// 2. 使用 iota 定义四季（Spring, Summer, Autumn, Winter）
// 3. 探索 string 和 []byte 之间的转换
// ========================================
