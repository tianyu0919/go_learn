package main

import (
	"fmt"
	"testing"
)

// ---- 基本测试 ----
// 测试函数必须以 Test 开头，接受 *testing.T
func TestAdd(t *testing.T) {
	result := Add(2, 3)
	if result != 5 {
		t.Errorf("Add(2, 3) = %d; want 5", result)
	}
}

// ---- 表驱动测试（Go 推荐的测试模式）----
func TestFizzBuzz(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected string
	}{
		{"regular number", 1, "1"},
		{"divisible by 3", 3, "Fizz"},
		{"divisible by 5", 5, "Buzz"},
		{"divisible by 15", 15, "FizzBuzz"},
		{"another fizz", 9, "Fizz"},
		{"another buzz", 10, "Buzz"},
		{"another fizzbuzz", 30, "FizzBuzz"},
		{"regular 7", 7, "7"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FizzBuzz(tt.input)
			if result != tt.expected {
				t.Errorf("FizzBuzz(%d) = %s; want %s",
					tt.input, result, tt.expected)
			}
		})
	}
}

// ---- 测试错误情况 ----
func TestMax(t *testing.T) {
	t.Run("normal case", func(t *testing.T) {
		result, err := Max(1, 5, 3, 9, 2)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != 9 {
			t.Errorf("Max(1,5,3,9,2) = %d; want 9", result)
		}
	})

	t.Run("single element", func(t *testing.T) {
		result, err := Max(42)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != 42 {
			t.Errorf("Max(42) = %d; want 42", result)
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		_, err := Max()
		if err == nil {
			t.Error("expected error for empty input, got nil")
		}
	})

	t.Run("negative numbers", func(t *testing.T) {
		result, err := Max(-3, -1, -5)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != -1 {
			t.Errorf("Max(-3,-1,-5) = %d; want -1", result)
		}
	})
}

// ---- Fibonacci 测试 ----
func TestFibonacci(t *testing.T) {
	tests := []struct {
		n        int
		expected int
	}{
		{0, 0},
		{1, 1},
		{2, 1},
		{3, 2},
		{4, 3},
		{5, 5},
		{10, 55},
		{20, 6765},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := Fibonacci(tt.n)
			if result != tt.expected {
				t.Errorf("Fibonacci(%d) = %d; want %d",
					tt.n, result, tt.expected)
			}
		})
	}
}

// ---- 基准测试 (Benchmark) ----
func BenchmarkFibonacci(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Fibonacci(20)
	}
}

func BenchmarkFizzBuzz(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FizzBuzz(i % 100)
	}
}

// ---- 示例测试（同时作为文档）----
func ExampleFizzBuzz() {
	fmt.Println(FizzBuzz(3))
	fmt.Println(FizzBuzz(5))
	fmt.Println(FizzBuzz(15))
	fmt.Println(FizzBuzz(7))
	// Output:
	// Fizz
	// Buzz
	// FizzBuzz
	// 7
}
