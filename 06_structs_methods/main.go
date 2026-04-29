// ========================================
// Lesson 06: Structs & Methods
// ========================================
// 结构体是 Go 中组织数据的主要方式（Go 没有 class）

package main

import (
	"fmt"
	"math"
)

// ---- 定义结构体 ----
type Person struct {
	Name string
	Age  int
	City string
}

// ---- 方法（Methods）----
// 方法是绑定到类型上的函数
// (p Person) 是接收者（receiver），类似其他语言的 this/self

// 值接收者：不会修改原始数据
func (p Person) Greet() string {
	return fmt.Sprintf("Hi, I'm %s, %d years old from %s", p.Name, p.Age, p.City)
}

// 指针接收者：可以修改原始数据
func (p *Person) SetAge(age int) {
	p.Age = age // 修改的是原始对象
}

// ---- 结构体嵌套（组合） ----
type Address struct {
	Street  string
	ZipCode string
}

type Employee struct {
	Person  // 嵌入（匿名字段），继承 Person 的字段和方法
	Address // 嵌入 Address
	Title   string
	Salary  float64
}

// Employee 自己的方法
func (e Employee) Summary() string {
	return fmt.Sprintf("%s - %s (%.0f/month)", e.Name, e.Title, e.Salary)
}

// ---- 几何图形示例 ----
type Point struct {
	X, Y float64
}

type Circle struct {
	Center Point
	Radius float64
}

type Rectangle struct {
	TopLeft     Point
	BottomRight Point
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

func (r Rectangle) Area() float64 {
	width := math.Abs(r.BottomRight.X - r.TopLeft.X)
	height := math.Abs(r.BottomRight.Y - r.TopLeft.Y)
	return width * height
}

func (r Rectangle) Perimeter() float64 {
	width := math.Abs(r.BottomRight.X - r.TopLeft.X)
	height := math.Abs(r.BottomRight.Y - r.TopLeft.Y)
	return 2 * (width + height)
}

func main() {
	// ---- 创建结构体 ----
	fmt.Println("--- Creating Structs ---")

	// 方式1: 按字段名
	p1 := Person{Name: "Alice", Age: 30, City: "Beijing"}

	// 方式2: 按顺序（不推荐，容易出错）
	p2 := Person{"Bob", 25, "Shanghai"}

	// 方式3: 指针
	p3 := &Person{Name: "Carol", Age: 28, City: "Shenzhen"}

	// 方式4: new（返回指针，零值初始化）
	p4 := new(Person)
	p4.Name = "David"
	p4.Age = 35

	fmt.Println(p1.Greet())
	fmt.Println(p2.Greet())
	fmt.Println(p3.Greet()) // 指针也能直接调用方法
	fmt.Println(p4.Greet())

	// ---- 修改结构体 ----
	fmt.Println("\n--- Modifying Structs ---")
	fmt.Printf("Before: %s, age %d\n", p1.Name, p1.Age)
	p1.SetAge(31) // 指针接收者方法
	fmt.Printf("After: %s, age %d\n", p1.Name, p1.Age)

	// ---- 结构体嵌套 ----
	fmt.Println("\n--- Embedded Structs ---")
	emp := Employee{
		Person:  Person{Name: "Eve", Age: 32, City: "Hangzhou"},
		Address: Address{Street: "Tech Road 100", ZipCode: "310000"},
		Title:   "Senior Engineer",
		Salary:  30000,
	}
	// 可以直接访问嵌入字段
	fmt.Println(emp.Name)       // 等同于 emp.Person.Name
	fmt.Println(emp.Street)     // 等同于 emp.Address.Street
	fmt.Println(emp.Greet())    // 继承 Person 的方法
	fmt.Println(emp.Summary())

	// ---- 几何图形 ----
	fmt.Println("\n--- Geometry ---")
	c := Circle{
		Center: Point{0, 0},
		Radius: 5,
	}
	fmt.Printf("Circle: area=%.2f, perimeter=%.2f\n", c.Area(), c.Perimeter())

	r := Rectangle{
		TopLeft:     Point{0, 0},
		BottomRight: Point{4, 3},
	}
	fmt.Printf("Rectangle: area=%.2f, perimeter=%.2f\n", r.Area(), r.Perimeter())

	// ---- 结构体比较 ----
	fmt.Println("\n--- Struct Comparison ---")
	point1 := Point{1, 2}
	point2 := Point{1, 2}
	point3 := Point{3, 4}
	fmt.Printf("point1 == point2: %t\n", point1 == point2) // true
	fmt.Printf("point1 == point3: %t\n", point1 == point3) // false

	// ---- 构造函数模式 ----
	fmt.Println("\n--- Constructor Pattern ---")
	p := NewPerson("Frank", 40, "Guangzhou")
	fmt.Println(p.Greet())
}

// Go 没有构造函数，惯例用 NewXxx 函数
func NewPerson(name string, age int, city string) *Person {
	return &Person{
		Name: name,
		Age:  age,
		City: city,
	}
}

// ========================================
// 练习:
// 1. 定义一个 BankAccount 结构体（Owner, Balance），
//    添加 Deposit, Withdraw, String 方法
// 2. 用结构体嵌套实现 Student（嵌入 Person，添加 Grade, School 字段）
// 3. 实现一个 Stack（栈）数据结构，包含 Push, Pop, Peek, IsEmpty 方法
// ========================================
