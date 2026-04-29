// ========================================
// Lesson 09: Concurrency (Goroutines & Channels)
// ========================================
// 并发是 Go 最强大的特性之一
// Goroutine = 轻量级线程, Channel = goroutine 之间的通信管道

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// ---- 基本 Goroutine ----
func sayHello(name string) {
	for i := 0; i < 3; i++ {
		fmt.Printf("[%s] Hello #%d\n", name, i+1)
		time.Sleep(100 * time.Millisecond)
	}
}

// ---- Channel 基础 ----
func producer(ch chan<- int) { // chan<- 表示只写 channel
	for i := 0; i < 5; i++ {
		fmt.Printf("  Producing: %d\n", i)
		ch <- i // 发送数据到 channel
		time.Sleep(100 * time.Millisecond)
	}
	close(ch) // 关闭 channel，通知消费者没有更多数据
}

func consumer(ch <-chan int) { // <-chan 表示只读 channel
	// range 会在 channel 关闭后自动退出
	for val := range ch {
		fmt.Printf("  Consumed: %d\n", val)
	}
}

// ---- 带缓冲的 Channel ----
func bufferedDemo() {
	// 无缓冲 channel：发送方会阻塞直到接收方准备好
	// 有缓冲 channel：可以存储指定数量的值
	ch := make(chan string, 3) // 缓冲区大小为 3

	ch <- "first"
	ch <- "second"
	ch <- "third"
	// ch <- "fourth" // 这会阻塞！因为缓冲区满了

	fmt.Println(<-ch) // first
	fmt.Println(<-ch) // second
	fmt.Println(<-ch) // third
}

// ---- Select 语句 ----
func selectDemo() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(100 * time.Millisecond)
		ch1 <- "from channel 1"
	}()

	go func() {
		time.Sleep(200 * time.Millisecond)
		ch2 <- "from channel 2"
	}()

	// select 等待多个 channel，哪个先就绪就执行哪个
	for i := 0; i < 2; i++ {
		select {
		case msg := <-ch1:
			fmt.Println("  Received:", msg)
		case msg := <-ch2:
			fmt.Println("  Received:", msg)
		}
	}
}

// ---- 超时控制 ----
func timeoutDemo() {
	ch := make(chan string)

	go func() {
		time.Sleep(2 * time.Second) // 模拟慢操作
		ch <- "result"
	}()

	select {
	case result := <-ch:
		fmt.Println("  Got:", result)
	case <-time.After(500 * time.Millisecond):
		fmt.Println("  Timeout! Operation took too long")
	}
}

// ---- WaitGroup：等待一组 goroutine 完成 ----
func waitGroupDemo() {
	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1) // 计数器 +1

		go func(id int) {
			defer wg.Done() // 计数器 -1

			duration := time.Duration(rand.Intn(500)) * time.Millisecond
			time.Sleep(duration)
			fmt.Printf("  Worker %d done (took %v)\n", id, duration)
		}(i)
	}

	wg.Wait() // 等待所有 goroutine 完成
	fmt.Println("  All workers finished!")
}

// ---- Mutex：互斥锁 ----
type SafeCounter struct {
	mu    sync.Mutex
	count int
}

func (c *SafeCounter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

// ---- 实际案例：并发 Web 爬虫模拟 ----
func fetchURL(url string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	// 模拟网络请求
	duration := time.Duration(100+rand.Intn(400)) * time.Millisecond
	time.Sleep(duration)
	results <- fmt.Sprintf("%s (took %v)", url, duration)
}

func concurrentCrawler() {
	urls := []string{
		"https://example.com/page1",
		"https://example.com/page2",
		"https://example.com/page3",
		"https://example.com/page4",
		"https://example.com/page5",
	}

	results := make(chan string, len(urls))
	var wg sync.WaitGroup

	start := time.Now()
	for _, url := range urls {
		wg.Add(1)
		go fetchURL(url, results, &wg)
	}

	// 在另一个 goroutine 中等待并关闭 channel
	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Printf("  Fetched: %s\n", result)
	}
	fmt.Printf("  Total time: %v (concurrent!)\n", time.Since(start))
}

func main() {
	// ---- Goroutines ----
	fmt.Println("--- Goroutines ---")
	go sayHello("goroutine-1") // go 关键字启动 goroutine
	go sayHello("goroutine-2")
	sayHello("main") // main 也在执行
	time.Sleep(500 * time.Millisecond)

	// ---- Channels ----
	fmt.Println("\n--- Channels ---")
	ch := make(chan int) // 创建无缓冲 channel
	go producer(ch)
	consumer(ch)

	// ---- Buffered Channel ----
	fmt.Println("\n--- Buffered Channel ---")
	bufferedDemo()

	// ---- Select ----
	fmt.Println("\n--- Select ---")
	selectDemo()

	// ---- Timeout ----
	fmt.Println("\n--- Timeout ---")
	timeoutDemo()

	// ---- WaitGroup ----
	fmt.Println("\n--- WaitGroup ---")
	waitGroupDemo()

	// ---- Mutex ----
	fmt.Println("\n--- Mutex ---")
	counter := &SafeCounter{}
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}
	wg.Wait()
	fmt.Printf("  Counter value: %d (should be 1000)\n", counter.Value())

	// ---- 并发爬虫 ----
	fmt.Println("\n--- Concurrent Crawler ---")
	concurrentCrawler()
}

// ========================================
// 练习:
// 1. 实现一个 worker pool：N 个 worker 处理 M 个任务
// 2. 用 channel 实现 Fan-in 模式（多个 channel 合并为一个）
// 3. 实现一个并发安全的缓存（支持过期时间）
// ========================================
