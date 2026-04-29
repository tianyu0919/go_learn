// ========================================
// Lesson 12: File I/O
// ========================================
// Go 标准库提供了丰富的文件操作功能

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// ---- 写文件 ----
	fmt.Println("--- Writing Files ---")
	writeSimpleFile()
	writeWithBufio()

	// ---- 读文件 ----
	fmt.Println("\n--- Reading Files ---")
	readEntireFile()
	readLineByLine()

	// ---- 文件操作 ----
	fmt.Println("\n--- File Operations ---")
	fileOperations()

	// ---- 目录操作 ----
	fmt.Println("\n--- Directory Operations ---")
	directoryOperations()

	// 清理
	cleanup()
}

func writeSimpleFile() {
	// 方式1: os.WriteFile（最简单）
	content := []byte("Hello, File I/O!\nThis is line 2.\nThis is line 3.\n")
	err := os.WriteFile("example.txt", content, 0644)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Written example.txt")

	// 方式2: os.Create + Write
	file, err := os.Create("example2.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close() // 重要！确保文件被关闭

	file.WriteString("Written with os.Create\n")
	file.WriteString("Second line\n")
	fmt.Println("Written example2.txt")
}

func writeWithBufio() {
	file, err := os.Create("buffered.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	// bufio.Writer 提供缓冲写入，减少系统调用
	writer := bufio.NewWriter(file)
	for i := 1; i <= 5; i++ {
		fmt.Fprintf(writer, "Line %d: Hello from bufio!\n", i)
	}
	writer.Flush() // 确保缓冲区内容写入文件
	fmt.Println("Written buffered.txt")
}

func readEntireFile() {
	// 方式1: os.ReadFile（读取整个文件）
	content, err := os.ReadFile("example.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("File content (%d bytes):\n%s", len(content), string(content))
}

func readLineByLine() {
	// 方式2: 逐行读取（适合大文件）
	file, err := os.Open("buffered.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	fmt.Println("\nReading line by line:")
	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		fmt.Printf("  [%d] %s\n", lineNum, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning:", err)
	}
}

func fileOperations() {
	// 检查文件是否存在
	if _, err := os.Stat("example.txt"); err == nil {
		fmt.Println("example.txt exists")
	} else if os.IsNotExist(err) {
		fmt.Println("example.txt does not exist")
	}

	// 复制文件
	err := copyFile("example.txt", "example_copy.txt")
	if err != nil {
		fmt.Println("Copy error:", err)
	} else {
		fmt.Println("Copied example.txt -> example_copy.txt")
	}

	// 获取文件信息
	info, err := os.Stat("example.txt")
	if err == nil {
		fmt.Printf("Name: %s, Size: %d bytes, ModTime: %s\n",
			info.Name(), info.Size(), info.ModTime().Format("2006-01-02 15:04:05"))
	}

	// 追加内容到文件
	file, err := os.OpenFile("example.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err == nil {
		defer file.Close()
		file.WriteString("Appended line!\n")
		fmt.Println("Appended to example.txt")
	}
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

func directoryOperations() {
	// 创建目录
	err := os.MkdirAll("testdir/subdir", 0755)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Created testdir/subdir")

	// 在目录中创建文件
	for i := 1; i <= 3; i++ {
		filename := filepath.Join("testdir", fmt.Sprintf("file%d.txt", i))
		os.WriteFile(filename, []byte(fmt.Sprintf("Content of file %d", i)), 0644)
	}

	// 遍历目录
	fmt.Println("Directory contents:")
	entries, err := os.ReadDir("testdir")
	if err == nil {
		for _, entry := range entries {
			kind := "FILE"
			if entry.IsDir() {
				kind = "DIR "
			}
			fmt.Printf("  [%s] %s\n", kind, entry.Name())
		}
	}

	// filepath.Walk 递归遍历
	fmt.Println("\nRecursive walk:")
	filepath.Walk("testdir", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		indent := strings.Repeat("  ", strings.Count(path, string(os.PathSeparator)))
		fmt.Printf("%s%s\n", indent, info.Name())
		return nil
	})
}

func cleanup() {
	os.Remove("example.txt")
	os.Remove("example2.txt")
	os.Remove("example_copy.txt")
	os.Remove("buffered.txt")
	os.RemoveAll("testdir")
	fmt.Println("\nCleaned up temporary files")
}

// ========================================
// 练习:
// 1. 实现一个简单的 CSV 读取器
// 2. 实现文件搜索：在目录中查找包含指定关键词的文件
// 3. 实现一个简单的日志文件轮转（按大小拆分）
// ========================================
