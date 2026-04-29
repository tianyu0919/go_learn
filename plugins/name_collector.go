package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const fileName = "names.json"
const tmpFileName = "names.json.tmp"

// 6A 工作流已激活
// 函数级注释：main 是程序的入口，负责收集终端输入的人名，去重后安全地追加到 JSON 文件中
func main() {
	var names []string
	// 使用 map 进行 O(1) 复杂度的去重判断
	seen := make(map[string]struct{})

	// 尝试读取已存在的 json 文件，实现继续拼接
	data, err := os.ReadFile(fileName)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &names)
		if err != nil {
			fmt.Printf("读取现有文件失败，可能格式不正确: %v\n", err)
		} else {
			// 初始化去重 map，将历史数据加入其中
			for _, name := range names {
				seen[name] = struct{}{}
			}
		}
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("当前已有 %d 条数据。\n", len(names))
	fmt.Println("请输入人名（支持批量粘贴，多个名字用空格或逗号分隔，输入 'exit' 或 'q' 退出）：")

	for {
		fmt.Print("> ")
		// 使用 ReadString 而不是 Scanner，因为 Scanner 默认有 64KB 的行长度限制
		// ReadString 会动态扩容，即使一次性粘贴 1 万条数据在一行内，也能安全读取
		input, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		input = strings.TrimSpace(input)
		if input == "exit" || input == "q" {
			break
		}
		if input == "" {
			continue
		}

		// 处理输入，支持空格和中文、英文逗号分隔
		input = strings.ReplaceAll(input, ",", " ")
		input = strings.ReplaceAll(input, "，", " ")
		parts := strings.Fields(input)

		var addedCount int
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part != "" {
				// 去重逻辑：只有没见过的名字才会写入
				if _, exists := seen[part]; !exists {
					seen[part] = struct{}{}
					names = append(names, part)
					addedCount++
				}
			}
		}

		if addedCount == 0 {
			if len(parts) > 0 {
				fmt.Println("输入的数据已存在（触发去重），跳过写入。")
			}
			continue
		}

		// 写入 json (带缩进)
		// 1 万条字符串序列化大约只有 150KB ~ 200KB，在内存中处理只需几毫秒，性能完全没问题
		fileData, err := json.MarshalIndent(names, "", "  ")
		if err != nil {
			fmt.Println("生成 JSON 失败:", err)
			continue
		}

		// 【关键优化】原子写入策略：先写入临时文件，再重命名。
		// 避免在频繁写入 1 万多条数据时，如果程序突然被强杀或断电，导致原 json 文件被截断损坏。
		err = os.WriteFile(tmpFileName, fileData, 0644)
		if err != nil {
			fmt.Println("写入临时文件失败:", err)
			continue
		}
		err = os.Rename(tmpFileName, fileName)
		if err != nil {
			fmt.Println("保存文件失败:", err)
			continue
		}

		fmt.Printf("本次成功写入 %d 条新记录（已自动去重），当前共 %d 条。\n", addedCount, len(names))
	}
	fmt.Println("程序已退出。")
}
