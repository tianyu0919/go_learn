package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

const dbDir = "db"
const dbFile = "db/names.db"

// NameRecord 是对应数据库 names 表的模型
type NameRecord struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex"`
}

// RequestPayload 用于解析前端传来的 JSON 数据
type RequestPayload struct {
	Names string `json:"names"`
}

// ResponsePayload 用于给前端返回处理结果
type ResponsePayload struct {
	Added int    `json:"added"`
	Total int64  `json:"total"`
	Error string `json:"error,omitempty"`
}

func initDB() {
	// 确保 db 目录存在
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Fatalf("创建数据库目录失败: %v\n", err)
	}

	var err error
	// 连接 SQLite 数据库，并配置 GORM 不输出普通的 SQL 执行日志
	db, err = gorm.Open(sqlite.Open(dbFile), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("连接数据库失败: %v\n", err)
	}

	// GORM 自动迁移功能：根据 NameRecord 结构体自动创建表或修改表结构
	if err := db.AutoMigrate(&NameRecord{}); err != nil {
		log.Fatalf("自动迁移数据表失败: %v\n", err)
	}

	// 检查是否需要从 names.json 迁移数据
	migrateDataFromJSON()
}

func migrateDataFromJSON() {
	const oldFileName = "names.json"
	data, err := os.ReadFile(oldFileName)
	if err != nil {
		return // 文件不存在或无法读取，跳过迁移
	}

	var names []string
	if err := json.Unmarshal(data, &names); err != nil {
		return // JSON 格式不正确，跳过迁移
	}

	if len(names) == 0 {
		return
	}

	// 检查数据库中是否已有数据
	var count int64
	db.Model(&NameRecord{}).Count(&count)
	if count > 0 {
		return // 数据库已有数据，不再自动迁移
	}

	log.Printf("正在从 %s 迁移 %d 条数据到 SQLite...\n", oldFileName, len(names))

	// 批量插入并忽略冲突
	records := make([]NameRecord, 0, len(names))
	for _, name := range names {
		records = append(records, NameRecord{Name: name})
	}

	// 使用 GORM 的 clause.OnConflict 实现 INSERT OR IGNORE
	result := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&records)
	if result.Error != nil {
		log.Printf("迁移失败: %v\n", result.Error)
		return
	}

	log.Printf("迁移完成，成功导入 %d 条数据。\n", result.RowsAffected)
}

func main() {
	initDB()

	// 静态文件服务：直接响应 index.html
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	// API 接口：处理提交的人名
	http.HandleFunc("/api/names", handleNames)

	fmt.Println("服务已启动，请在浏览器访问 http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("启动服务失败:", err)
	}
}

func handleNames(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendJSON(w, ResponsePayload{Error: "只支持 POST 请求"}, http.StatusMethodNotAllowed)
		return
	}

	var payload RequestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		sendJSON(w, ResponsePayload{Error: "无效的 JSON 数据"}, http.StatusBadRequest)
		return
	}

	lines := strings.Split(payload.Names, "\n")
	var records []NameRecord

	// 清洗数据并封装到 struct
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			records = append(records, NameRecord{Name: line})
		}
	}

	addedCount := int64(0)
	if len(records) > 0 {
		// 使用 GORM 的 Create 批量插入，并遇到冲突自动忽略
		result := db.Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(&records)

		if result.Error != nil {
			sendJSON(w, ResponsePayload{Error: "执行插入失败"}, http.StatusInternalServerError)
			return
		}
		addedCount = result.RowsAffected
	}

	// 查询当前总数
	var total int64
	if err := db.Model(&NameRecord{}).Count(&total).Error; err != nil {
		sendJSON(w, ResponsePayload{Error: "查询总数失败"}, http.StatusInternalServerError)
		return
	}

	// 成功响应，返回当前新增数和总数
	sendJSON(w, ResponsePayload{
		Added: int(addedCount),
		Total: total,
	}, http.StatusOK)
}

func sendJSON(w http.ResponseWriter, resp ResponsePayload, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}
