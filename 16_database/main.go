// ========================================
// Lesson 16: Database Operations
// ========================================
// Go 通过 database/sql 标准库操作数据库
// 本课使用 SQLite（无需安装数据库服务器）
//
// 安装驱动: go get github.com/mattn/go-sqlite3
// 注意: go-sqlite3 需要 CGO 支持
// 或使用纯 Go 驱动: go get modernc.org/sqlite

package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite" // 导入驱动（纯 Go，无需 CGO）
)

type User struct {
	ID        int
	Name      string
	Email     string
	CreatedAt time.Time
}

// ---- 数据库初始化 ----
func initDB(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		content TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);
	`
	_, err := db.Exec(query)
	return err
}

// ---- CRUD 操作 ----

// Create
func createUser(db *sql.DB, name, email string) (int64, error) {
	result, err := db.Exec(
		"INSERT INTO users (name, email) VALUES (?, ?)",
		name, email,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// Read - 单条
func getUserByID(db *sql.DB, id int) (*User, error) {
	user := &User{}
	err := db.QueryRow(
		"SELECT id, name, email, created_at FROM users WHERE id = ?", id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Read - 多条
func getAllUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query("SELECT id, name, email, created_at FROM users ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

// Update
func updateUser(db *sql.DB, id int, name, email string) error {
	result, err := db.Exec(
		"UPDATE users SET name = ?, email = ? WHERE id = ?",
		name, email, id,
	)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user %d not found", id)
	}
	return nil
}

// Delete
func deleteUser(db *sql.DB, id int) error {
	result, err := db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user %d not found", id)
	}
	return nil
}

// ---- 事务 (Transaction) ----
func createUserWithPosts(db *sql.DB, name, email string, postTitles []string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	// 如果出错，回滚事务
	defer tx.Rollback()

	// 创建用户
	result, err := tx.Exec("INSERT INTO users (name, email) VALUES (?, ?)", name, email)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	userID, _ := result.LastInsertId()

	// 创建文章
	stmt, err := tx.Prepare("INSERT INTO posts (user_id, title) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	for _, title := range postTitles {
		if _, err := stmt.Exec(userID, title); err != nil {
			return fmt.Errorf("create post '%s': %w", title, err)
		}
	}

	// 提交事务
	return tx.Commit()
}

// ---- 联表查询 ----
type UserPost struct {
	UserName  string
	PostTitle string
	PostDate  time.Time
}

func getUserPosts(db *sql.DB) ([]UserPost, error) {
	query := `
		SELECT u.name, p.title, p.created_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
		ORDER BY p.created_at DESC
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []UserPost
	for rows.Next() {
		var up UserPost
		if err := rows.Scan(&up.UserName, &up.PostTitle, &up.PostDate); err != nil {
			return nil, err
		}
		posts = append(posts, up)
	}
	return posts, rows.Err()
}

func main() {
	// 使用内存数据库（也可以用文件：./test.db）
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 初始化表
	if err := initDB(db); err != nil {
		log.Fatal("init db:", err)
	}
	fmt.Println("Database initialized")

	// ---- Create ----
	fmt.Println("\n--- Create ---")
	id1, _ := createUser(db, "Alice", "alice@example.com")
	id2, _ := createUser(db, "Bob", "bob@example.com")
	fmt.Printf("Created users with IDs: %d, %d\n", id1, id2)

	// ---- Read ----
	fmt.Println("\n--- Read ---")
	user, err := getUserByID(db, int(id1))
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("User: %+v\n", *user)
	}

	// 查询不存在的用户
	_, err = getUserByID(db, 999)
	if err == sql.ErrNoRows {
		fmt.Println("User 999 not found (sql.ErrNoRows)")
	}

	// ---- Update ----
	fmt.Println("\n--- Update ---")
	updateUser(db, int(id1), "Alice Updated", "alice.new@example.com")
	user, _ = getUserByID(db, int(id1))
	fmt.Printf("Updated: %+v\n", *user)

	// ---- List all ----
	fmt.Println("\n--- List All ---")
	users, _ := getAllUsers(db)
	for _, u := range users {
		fmt.Printf("  [%d] %s (%s)\n", u.ID, u.Name, u.Email)
	}

	// ---- Transaction ----
	fmt.Println("\n--- Transaction ---")
	err = createUserWithPosts(db, "Carol", "carol@example.com", []string{
		"My First Post",
		"Learning Go",
		"Database Fun",
	})
	if err != nil {
		fmt.Println("Transaction error:", err)
	} else {
		fmt.Println("Transaction committed successfully")
	}

	// ---- Join Query ----
	fmt.Println("\n--- Join Query ---")
	posts, _ := getUserPosts(db)
	for _, p := range posts {
		fmt.Printf("  %s - '%s'\n", p.UserName, p.PostTitle)
	}

	// ---- Delete ----
	fmt.Println("\n--- Delete ---")
	deleteUser(db, int(id2))
	users, _ = getAllUsers(db)
	fmt.Printf("After delete, %d users remain\n", len(users))
}

// ========================================
// 练习:
// 1. 添加分页查询: getUsers(db, page, pageSize int) ([]User, int, error)
// 2. 实现软删除（添加 deleted_at 字段）
// 3. 将数据库操作封装为 Repository 接口
// ========================================
