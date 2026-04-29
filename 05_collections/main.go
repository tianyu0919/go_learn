// ========================================
// Lesson 05: Arrays, Slices & Maps
// ========================================
// Go 最常用的数据结构

package main

import (
	"fmt"
	"sort"
)

func main() {
	// ---- 数组 (Array) ----
	// 数组是固定长度的，长度是类型的一部分
	fmt.Println("--- Arrays ---")

	var arr1 [5]int                        // 零值初始化
	arr2 := [3]string{"Go", "is", "fun"}   // 字面量初始化
	arr3 := [...]int{1, 2, 3, 4}           // ... 让编译器计算长度

	fmt.Println("arr1:", arr1)
	fmt.Println("arr2:", arr2)
	fmt.Println("arr3:", arr3, "len:", len(arr3))

	// 数组是值类型（赋值会复制）
	arr4 := arr3
	arr4[0] = 999
	fmt.Println("arr3 (original):", arr3) // 不受影响
	fmt.Println("arr4 (copy):", arr4)

	// ---- 切片 (Slice) ----
	// 切片是动态数组，是 Go 中最常用的集合类型
	fmt.Println("\n--- Slices ---")

	// 创建切片
	s1 := []int{1, 2, 3, 4, 5}           // 字面量创建
	s2 := make([]int, 5)                  // make 创建，长度 5
	s3 := make([]int, 3, 10)              // 长度 3，容量 10

	fmt.Printf("s1: %v, len=%d, cap=%d\n", s1, len(s1), cap(s1))
	fmt.Printf("s2: %v, len=%d, cap=%d\n", s2, len(s2), cap(s2))
	fmt.Printf("s3: %v, len=%d, cap=%d\n", s3, len(s3), cap(s3))

	// append：添加元素（可能导致扩容）
	s3 = append(s3, 1, 2, 3)
	fmt.Printf("after append: %v, len=%d, cap=%d\n", s3, len(s3), cap(s3))

	// 切片操作 [low:high]
	fmt.Println("\n--- Slice Operations ---")
	data := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Println("data:", data)
	fmt.Println("data[2:5]:", data[2:5])   // [2, 3, 4]
	fmt.Println("data[:3]:", data[:3])       // [0, 1, 2]
	fmt.Println("data[7:]:", data[7:])       // [7, 8, 9]

	// 重要：切片是引用类型！切片共享底层数组
	sub := data[2:5]
	sub[0] = 999
	fmt.Println("after modifying sub, data:", data) // data[2] 也变了！

	// 如果需要独立副本，使用 copy
	src := []int{1, 2, 3}
	dst := make([]int, len(src))
	copy(dst, src)
	dst[0] = 999
	fmt.Println("src (unchanged):", src)
	fmt.Println("dst (modified):", dst)

	// 删除元素（Go 没有内置 delete for slice）
	s := []int{1, 2, 3, 4, 5}
	// 删除索引 2 的元素
	s = append(s[:2], s[3:]...)
	fmt.Println("after delete index 2:", s) // [1, 2, 4, 5]

	// 遍历切片
	fmt.Println("\n--- Iterating ---")
	fruits := []string{"apple", "banana", "cherry"}
	for i, fruit := range fruits {
		fmt.Printf("  %d: %s\n", i, fruit)
	}

	// 排序
	nums := []int{5, 3, 8, 1, 9, 2}
	sort.Ints(nums)
	fmt.Println("sorted:", nums)

	// ---- Map ----
	// map 是键值对集合（类似其他语言的 dict/HashMap）
	fmt.Println("\n--- Maps ---")

	// 创建 map
	m1 := map[string]int{
		"alice": 90,
		"bob":   85,
		"carol": 92,
	}
	m2 := make(map[string]string) // 空 map

	fmt.Println("m1:", m1)
	fmt.Println("m2:", m2)

	// 增删改查
	m2["name"] = "Go"       // 添加/修改
	m2["version"] = "1.25"
	fmt.Println("m2 after set:", m2)

	delete(m2, "version")    // 删除
	fmt.Println("m2 after delete:", m2)

	// 查找（重要！使用 comma-ok 模式）
	if score, ok := m1["alice"]; ok {
		fmt.Printf("alice's score: %d\n", score)
	} else {
		fmt.Println("alice not found")
	}

	if _, ok := m1["david"]; !ok {
		fmt.Println("david not found")
	}

	// 遍历 map（顺序不保证！）
	fmt.Println("\n--- Map Iteration ---")
	for name, score := range m1 {
		fmt.Printf("  %s: %d\n", name, score)
	}

	// 统计词频（实际用例）
	fmt.Println("\n--- Word Count Example ---")
	wordCount("the quick brown fox jumps over the lazy dog the fox")

	// ---- nil slice vs empty slice ----
	fmt.Println("\n--- nil vs empty ---")
	var nilSlice []int
	emptySlice := []int{}
	fmt.Printf("nil slice: %v, len=%d, is nil? %t\n",
		nilSlice, len(nilSlice), nilSlice == nil)
	fmt.Printf("empty slice: %v, len=%d, is nil? %t\n",
		emptySlice, len(emptySlice), emptySlice == nil)
	// 两者都可以安全地使用 append
}

func wordCount(text string) {
	counts := make(map[string]int)
	// 简单按空格分割
	words := []string{}
	word := ""
	for _, ch := range text {
		if ch == ' ' {
			if word != "" {
				words = append(words, word)
			}
			word = ""
		} else {
			word += string(ch)
		}
	}
	if word != "" {
		words = append(words, word)
	}

	for _, w := range words {
		counts[w]++
	}
	for w, c := range counts {
		fmt.Printf("  '%s': %d\n", w, c)
	}
}

// ========================================
// 练习:
// 1. 实现一个函数 removeDuplicates(nums []int) []int
// 2. 实现一个函数 invertMap(m map[string]int) map[int]string
// 3. 用 map 实现一个简单的电话簿（增删改查）
// ========================================
