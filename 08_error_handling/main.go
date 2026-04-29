// ========================================
// Lesson 08: Error Handling
// ========================================
// Go 使用显式的错误返回值，而不是 try-catch 异常机制
// 错误处理是 Go 代码中非常重要的部分

package main

import (
	"errors"
	"fmt"
	"strconv"
)

// ---- 基本错误处理 ----
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("cannot divide %.1f by zero", a)
	}
	return a / b, nil
}

// ---- 自定义错误类型 ----
type ValidationError struct {
	Field   string
	Message string
}

// 实现 error 接口（只需要 Error() string 方法）
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: field '%s' - %s", e.Field, e.Message)
}

func validateAge(age int) error {
	if age < 0 {
		return &ValidationError{Field: "age", Message: "cannot be negative"}
	}
	if age > 150 {
		return &ValidationError{Field: "age", Message: "unrealistic value"}
	}
	return nil
}

// ---- Sentinel Errors（哨兵错误）----
var (
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
)

type UserStore struct {
	users map[string]string
}

func (s *UserStore) GetUser(id string) (string, error) {
	user, ok := s.users[id]
	if !ok {
		return "", fmt.Errorf("user %s: %w", id, ErrNotFound) // %w 包装错误
	}
	return user, nil
}

// ---- 错误包装和解包 ----
func processUserRequest(store *UserStore, userID string) error {
	user, err := store.GetUser(userID)
	if err != nil {
		// 包装错误，添加上下文
		return fmt.Errorf("processing request for user %s: %w", userID, err)
	}
	fmt.Printf("Processing request for: %s\n", user)
	return nil
}

// ---- panic 和 recover ----
func safeDivide(a, b int) (result int, err error) {
	// defer + recover 捕获 panic
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic: %v", r)
		}
	}()

	// 这会 panic（整数除以零）
	return a / b, nil
}

// ---- 多错误收集 ----
type MultiError struct {
	Errors []error
}

func (m *MultiError) Error() string {
	msgs := make([]string, len(m.Errors))
	for i, err := range m.Errors {
		msgs[i] = err.Error()
	}
	return fmt.Sprintf("%d errors occurred: %v", len(m.Errors), msgs)
}

func (m *MultiError) Add(err error) {
	if err != nil {
		m.Errors = append(m.Errors, err)
	}
}

func (m *MultiError) HasErrors() bool {
	return len(m.Errors) > 0
}

func validateForm(name string, age string, email string) error {
	me := &MultiError{}

	if name == "" {
		me.Add(&ValidationError{Field: "name", Message: "required"})
	}

	if ageInt, err := strconv.Atoi(age); err != nil {
		me.Add(&ValidationError{Field: "age", Message: "must be a number"})
	} else if err := validateAge(ageInt); err != nil {
		me.Add(err)
	}

	if email == "" {
		me.Add(&ValidationError{Field: "email", Message: "required"})
	}

	if me.HasErrors() {
		return me
	}
	return nil
}

func main() {
	// ---- 基本错误处理 ----
	fmt.Println("--- Basic Error Handling ---")
	result, err := divide(10, 3)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("10 / 3 = %.2f\n", result)
	}

	_, err = divide(10, 0)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// ---- 自定义错误 ----
	fmt.Println("\n--- Custom Errors ---")
	if err := validateAge(-5); err != nil {
		fmt.Println(err)

		// 类型断言检查具体错误类型
		var ve *ValidationError
		if errors.As(err, &ve) {
			fmt.Printf("  Field: %s, Message: %s\n", ve.Field, ve.Message)
		}
	}

	// ---- 错误包装和 errors.Is ----
	fmt.Println("\n--- Error Wrapping ---")
	store := &UserStore{
		users: map[string]string{
			"1": "Alice",
			"2": "Bob",
		},
	}

	err = processUserRequest(store, "3")
	if err != nil {
		fmt.Println("Error:", err)

		// errors.Is 可以检查错误链中是否包含特定错误
		if errors.Is(err, ErrNotFound) {
			fmt.Println("  -> User was not found (detected via errors.Is)")
		}
	}

	err = processUserRequest(store, "1")
	fmt.Println()

	// ---- panic 和 recover ----
	fmt.Println("--- Panic & Recover ---")
	result2, err := safeDivide(10, 0)
	if err != nil {
		fmt.Println("Caught:", err)
	} else {
		fmt.Println("Result:", result2)
	}

	result2, err = safeDivide(10, 3)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("10 / 3 =", result2)
	}

	// ---- 表单验证（多错误）----
	fmt.Println("\n--- Form Validation ---")
	if err := validateForm("", "abc", ""); err != nil {
		fmt.Println(err)
	}

	if err := validateForm("Alice", "25", "alice@example.com"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Form is valid!")
	}

	// ---- Go 的错误处理惯例 ----
	fmt.Println("\n--- Best Practices ---")
	fmt.Println("1. Always check errors immediately after function calls")
	fmt.Println("2. Use fmt.Errorf with %w to wrap errors with context")
	fmt.Println("3. Use errors.Is/As to check error types in the chain")
	fmt.Println("4. Don't panic for expected errors; panic only for bugs")
	fmt.Println("5. Custom error types when you need structured error info")
}

// ========================================
// 练习:
// 1. 实现一个 ParseConfig 函数，读取配置并返回结构化的错误
// 2. 创建一个错误重试机制: retry(fn func() error, maxRetries int) error
// 3. 实现 errors.Is 和 errors.As 的完整示例链
// ========================================
