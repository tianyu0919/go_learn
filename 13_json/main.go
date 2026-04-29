// ========================================
// Lesson 13: JSON Handling
// ========================================
// JSON 是 Web 开发中最常用的数据格式
// Go 通过 encoding/json 包提供原生支持

package main

import (
	"encoding/json"
	"fmt"
)

// ---- 基本结构体与 JSON 标签 ----
type User struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Email     string   `json:"email"`
	Age       int      `json:"age,omitempty"`       // omitempty: 零值时不输出
	Password  string   `json:"-"`                   // -: 永远不序列化
	Tags      []string `json:"tags,omitempty"`
	IsActive  bool     `json:"is_active"`
}

// 嵌套结构体
type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	Country string `json:"country"`
}

type Profile struct {
	User    User    `json:"user"`
	Address Address `json:"address"`
	Score   float64 `json:"score"`
}

func main() {
	// ---- 序列化（Marshal: struct -> JSON）----
	fmt.Println("--- Marshal (Struct -> JSON) ---")

	user := User{
		ID:       1,
		Name:     "Alice",
		Email:    "alice@example.com",
		Age:      30,
		Password: "secret123", // 不会出现在 JSON 中
		Tags:     []string{"admin", "developer"},
		IsActive: true,
	}

	// 普通序列化
	jsonBytes, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("JSON:", string(jsonBytes))

	// 格式化输出（MarshalIndent）
	prettyJSON, _ := json.MarshalIndent(user, "", "  ")
	fmt.Println("\nPretty JSON:")
	fmt.Println(string(prettyJSON))

	// omitempty 效果
	userNoAge := User{ID: 2, Name: "Bob", Email: "bob@example.com"}
	jsonNoAge, _ := json.MarshalIndent(userNoAge, "", "  ")
	fmt.Println("\nWith omitempty (no age, no tags):")
	fmt.Println(string(jsonNoAge))

	// ---- 反序列化（Unmarshal: JSON -> struct）----
	fmt.Println("\n--- Unmarshal (JSON -> Struct) ---")

	jsonStr := `{
		"id": 3,
		"name": "Carol",
		"email": "carol@example.com",
		"age": 28,
		"tags": ["designer", "writer"],
		"is_active": true
	}`

	var decoded User
	err = json.Unmarshal([]byte(jsonStr), &decoded)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Decoded: %+v\n", decoded)

	// ---- 嵌套结构体 ----
	fmt.Println("\n--- Nested Structs ---")
	profile := Profile{
		User: User{
			ID:       4,
			Name:     "David",
			Email:    "david@example.com",
			IsActive: true,
		},
		Address: Address{
			Street:  "123 Main St",
			City:    "Beijing",
			Country: "China",
		},
		Score: 95.5,
	}
	profileJSON, _ := json.MarshalIndent(profile, "", "  ")
	fmt.Println(string(profileJSON))

	// ---- 动态 JSON（map[string]any）----
	fmt.Println("\n--- Dynamic JSON ---")

	// 当你不知道 JSON 结构时
	dynamicJSON := `{"name": "Eve", "scores": [85, 92, 78], "metadata": {"level": "advanced"}}`

	var result map[string]any
	json.Unmarshal([]byte(dynamicJSON), &result)

	fmt.Println("Name:", result["name"])
	fmt.Println("Scores:", result["scores"])
	fmt.Println("Metadata:", result["metadata"])

	// 访问嵌套值需要类型断言
	if metadata, ok := result["metadata"].(map[string]any); ok {
		fmt.Println("Level:", metadata["level"])
	}

	// ---- JSON 数组 ----
	fmt.Println("\n--- JSON Array ---")
	usersJSON := `[
		{"id": 1, "name": "Alice", "email": "alice@example.com", "is_active": true},
		{"id": 2, "name": "Bob", "email": "bob@example.com", "is_active": false},
		{"id": 3, "name": "Carol", "email": "carol@example.com", "is_active": true}
	]`

	var users []User
	json.Unmarshal([]byte(usersJSON), &users)

	for _, u := range users {
		status := "inactive"
		if u.IsActive {
			status = "active"
		}
		fmt.Printf("  %s (%s) - %s\n", u.Name, u.Email, status)
	}

	// ---- 构建动态 JSON ----
	fmt.Println("\n--- Building Dynamic JSON ---")
	response := map[string]any{
		"status":  "success",
		"code":    200,
		"message": "Data retrieved successfully",
		"data": map[string]any{
			"users_count": len(users),
			"page":        1,
			"per_page":    10,
		},
	}
	respJSON, _ := json.MarshalIndent(response, "", "  ")
	fmt.Println(string(respJSON))

	// ---- 自定义 JSON 序列化 ----
	fmt.Println("\n--- Custom Marshaling ---")
	order := Order{
		ID:     "ORD-001",
		Amount: 9999, // 存储为分
		Status: "paid",
	}
	orderJSON, _ := json.MarshalIndent(order, "", "  ")
	fmt.Println(string(orderJSON))
}

// ---- 自定义 JSON 序列化 ----
type Order struct {
	ID     string `json:"id"`
	Amount int    `json:"amount"` // 以分为单位
	Status string `json:"status"`
}

// 自定义 MarshalJSON 方法
func (o Order) MarshalJSON() ([]byte, error) {
	// 输出时将分转换为元
	type Alias Order // 避免递归
	return json.Marshal(&struct {
		Alias
		AmountYuan float64 `json:"amount_yuan"`
	}{
		Alias:      Alias(o),
		AmountYuan: float64(o.Amount) / 100,
	})
}

// ========================================
// 练习:
// 1. 创建一个配置文件系统：读取/写入 JSON 配置文件
// 2. 实现一个 API 响应解析器，处理不同格式的响应
// 3. 实现 JSON 到 CSV 的转换
// ========================================
