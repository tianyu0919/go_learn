// ========================================
// Lesson 07: Interfaces
// ========================================
// 接口是 Go 中实现多态和抽象的关键
// Go 的接口是隐式实现的（duck typing）—— 如果一个类型有接口要求的所有方法，
// 它就自动实现了该接口，不需要显式声明

package main

import (
	"fmt"
	"math"
	"strings"
)

// ---- 定义接口 ----
type Shape interface {
	Area() float64
	Perimeter() float64
}

// ---- 实现接口的类型 ----
type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

type Triangle struct {
	A, B, C float64 // 三条边
}

func (t Triangle) Area() float64 {
	// 海伦公式
	s := (t.A + t.B + t.C) / 2
	return math.Sqrt(s * (s - t.A) * (s - t.B) * (s - t.C))
}

func (t Triangle) Perimeter() float64 {
	return t.A + t.B + t.C
}

// ---- 使用接口的函数 ----
func printShapeInfo(s Shape) {
	fmt.Printf("  Type: %T\n", s)
	fmt.Printf("  Area: %.2f\n", s.Area())
	fmt.Printf("  Perimeter: %.2f\n", s.Perimeter())
}

func totalArea(shapes []Shape) float64 {
	total := 0.0
	for _, s := range shapes {
		total += s.Area()
	}
	return total
}

// ---- Stringer 接口（类似 toString）----
// fmt 包定义了 Stringer 接口: type Stringer interface { String() string }
func (c Circle) String() string {
	return fmt.Sprintf("Circle(r=%.1f)", c.Radius)
}

func (r Rectangle) String() string {
	return fmt.Sprintf("Rect(%.1f x %.1f)", r.Width, r.Height)
}

// ---- 接口组合 ----
type Reader interface {
	Read(p []byte) (n int, err error)
}

type Writer interface {
	Write(p []byte) (n int, err error)
}

// 接口可以组合
type ReadWriter interface {
	Reader
	Writer
}

// ---- 空接口 interface{} / any ----
func describe(i any) {
	fmt.Printf("  Value: %v, Type: %T\n", i, i)
}

// ---- 类型断言和类型开关 ----
func classifyShape(s Shape) string {
	// 类型断言：尝试将接口转换为具体类型
	switch v := s.(type) {
	case Circle:
		return fmt.Sprintf("It's a circle with radius %.1f", v.Radius)
	case Rectangle:
		return fmt.Sprintf("It's a %s rectangle",
			func() string {
				if v.Width == v.Height {
					return "square"
				}
				return "non-square"
			}())
	case Triangle:
		return fmt.Sprintf("It's a triangle with sides %.1f, %.1f, %.1f", v.A, v.B, v.C)
	default:
		return "Unknown shape"
	}
}

// ---- 实际案例：日志系统 ----
type Logger interface {
	Log(message string)
}

type ConsoleLogger struct {
	Prefix string
}

func (l ConsoleLogger) Log(message string) {
	fmt.Printf("[%s] %s\n", l.Prefix, message)
}

type FilterLogger struct {
	Inner   Logger
	Keyword string
}

func (l FilterLogger) Log(message string) {
	if strings.Contains(message, l.Keyword) {
		l.Inner.Log(message)
	}
}

func main() {
	// ---- 基本接口使用 ----
	fmt.Println("--- Shapes via Interface ---")
	shapes := []Shape{
		Circle{Radius: 5},
		Rectangle{Width: 4, Height: 3},
		Triangle{A: 3, B: 4, C: 5},
	}

	for _, s := range shapes {
		printShapeInfo(s)
		fmt.Println()
	}

	fmt.Printf("Total area: %.2f\n", totalArea(shapes))

	// ---- Stringer ----
	fmt.Println("\n--- Stringer Interface ---")
	c := Circle{Radius: 3}
	r := Rectangle{Width: 5, Height: 2}
	fmt.Println(c) // 自动调用 String()
	fmt.Println(r)

	// ---- 空接口 ----
	fmt.Println("\n--- Empty Interface (any) ---")
	describe(42)
	describe("hello")
	describe(true)
	describe(Circle{Radius: 1})

	// ---- 类型断言 ----
	fmt.Println("\n--- Type Assertion ---")
	var s Shape = Circle{Radius: 7}

	// 安全的类型断言（推荐）
	if circle, ok := s.(Circle); ok {
		fmt.Printf("It's a circle! Radius: %.1f\n", circle.Radius)
	}

	// ---- 类型开关 ----
	fmt.Println("\n--- Type Switch ---")
	for _, shape := range shapes {
		fmt.Println(classifyShape(shape))
	}

	// ---- 实际案例：日志 ----
	fmt.Println("\n--- Logger Example ---")
	console := ConsoleLogger{Prefix: "APP"}
	console.Log("Server started")
	console.Log("Handling request")

	// 装饰器模式：过滤日志
	errorLogger := FilterLogger{
		Inner:   ConsoleLogger{Prefix: "ERROR"},
		Keyword: "error",
	}
	errorLogger.Log("normal message")         // 被过滤
	errorLogger.Log("error: something failed") // 输出

	// ---- 接口值的 nil ----
	fmt.Println("\n--- Interface nil ---")
	var logger Logger // nil 接口
	fmt.Printf("nil interface: %v, is nil: %t\n", logger, logger == nil)
}

// ========================================
// 练习:
// 1. 定义 Animal 接口（Speak() string, Name() string），
//    实现 Dog, Cat, Duck 类型
// 2. 实现 Sorter 接口让自定义类型可以用 sort.Sort 排序
// 3. 实现一个 Storage 接口（Get, Set, Delete），
//    分别用 map 和 slice 实现
// ========================================
