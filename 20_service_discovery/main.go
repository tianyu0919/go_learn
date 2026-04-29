// ========================================
// Lesson 20: Service Discovery & Config
// ========================================
// 在微服务架构中，服务发现是关键组件
// 本课实现一个简化版的服务注册与发现机制
// 以及配置中心的基本概念

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// ========================================
// 服务注册中心
// ========================================

type ServiceInstance struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Port     int    `json:"port"`
	Metadata map[string]string `json:"metadata,omitempty"`
	Health   bool   `json:"healthy"`
	LastSeen time.Time `json:"last_seen"`
}

func (s *ServiceInstance) Endpoint() string {
	return fmt.Sprintf("%s:%d", s.Address, s.Port)
}

type Registry struct {
	mu       sync.RWMutex
	services map[string][]*ServiceInstance // name -> instances
}

func NewRegistry() *Registry {
	r := &Registry{
		services: make(map[string][]*ServiceInstance),
	}
	// 启动健康检查
	go r.healthCheck()
	return r
}

// Register 注册服务实例
func (r *Registry) Register(instance *ServiceInstance) {
	r.mu.Lock()
	defer r.mu.Unlock()

	instance.Health = true
	instance.LastSeen = time.Now()

	instances := r.services[instance.Name]
	// 检查是否已存在（更新）
	for i, inst := range instances {
		if inst.ID == instance.ID {
			instances[i] = instance
			log.Printf("[Registry] Updated: %s (%s)", instance.Name, instance.Endpoint())
			return
		}
	}

	r.services[instance.Name] = append(instances, instance)
	log.Printf("[Registry] Registered: %s (%s)", instance.Name, instance.Endpoint())
}

// Deregister 注销服务实例
func (r *Registry) Deregister(name, id string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	instances := r.services[name]
	for i, inst := range instances {
		if inst.ID == id {
			r.services[name] = append(instances[:i], instances[i+1:]...)
			log.Printf("[Registry] Deregistered: %s (%s)", name, id)
			return
		}
	}
}

// Discover 发现健康的服务实例
func (r *Registry) Discover(name string) []*ServiceInstance {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var healthy []*ServiceInstance
	for _, inst := range r.services[name] {
		if inst.Health {
			healthy = append(healthy, inst)
		}
	}
	return healthy
}

// Heartbeat 心跳更新
func (r *Registry) Heartbeat(name, id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, inst := range r.services[name] {
		if inst.ID == id {
			inst.LastSeen = time.Now()
			inst.Health = true
			return true
		}
	}
	return false
}

// 健康检查：标记超时实例为不健康
func (r *Registry) healthCheck() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		r.mu.Lock()
		for name, instances := range r.services {
			for _, inst := range instances {
				if time.Since(inst.LastSeen) > 15*time.Second {
					if inst.Health {
						inst.Health = false
						log.Printf("[Registry] Unhealthy: %s/%s (no heartbeat)", name, inst.ID)
					}
				}
			}
		}
		r.mu.Unlock()
	}
}

// ========================================
// 负载均衡器
// ========================================

type LoadBalancer interface {
	Pick(instances []*ServiceInstance) *ServiceInstance
}

// 轮询
type RoundRobinLB struct {
	mu      sync.Mutex
	counter int
}

func (lb *RoundRobinLB) Pick(instances []*ServiceInstance) *ServiceInstance {
	if len(instances) == 0 {
		return nil
	}
	lb.mu.Lock()
	defer lb.mu.Unlock()

	inst := instances[lb.counter%len(instances)]
	lb.counter++
	return inst
}

// 随机
type RandomLB struct{}

func (lb *RandomLB) Pick(instances []*ServiceInstance) *ServiceInstance {
	if len(instances) == 0 {
		return nil
	}
	return instances[rand.Intn(len(instances))]
}

// ========================================
// 服务客户端（带服务发现 + 负载均衡）
// ========================================

type ServiceClient struct {
	registry *Registry
	lb       LoadBalancer
}

func NewServiceClient(registry *Registry, lb LoadBalancer) *ServiceClient {
	return &ServiceClient{registry: registry, lb: lb}
}

func (c *ServiceClient) Call(ctx context.Context, serviceName string) (*ServiceInstance, error) {
	instances := c.registry.Discover(serviceName)
	if len(instances) == 0 {
		return nil, fmt.Errorf("no healthy instances for service: %s", serviceName)
	}

	instance := c.lb.Pick(instances)
	log.Printf("[Client] Calling %s at %s", serviceName, instance.Endpoint())
	return instance, nil
}

// ========================================
// 配置中心
// ========================================

type ConfigCenter struct {
	mu        sync.RWMutex
	configs   map[string]map[string]string // service -> key -> value
	watchers  map[string][]chan ConfigChange
}

type ConfigChange struct {
	Key      string
	OldValue string
	NewValue string
}

func NewConfigCenter() *ConfigCenter {
	return &ConfigCenter{
		configs:  make(map[string]map[string]string),
		watchers: make(map[string][]chan ConfigChange),
	}
}

func (c *ConfigCenter) Set(service, key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.configs[service] == nil {
		c.configs[service] = make(map[string]string)
	}

	oldValue := c.configs[service][key]
	c.configs[service][key] = value

	// 通知 watchers
	if oldValue != value {
		change := ConfigChange{Key: key, OldValue: oldValue, NewValue: value}
		for _, ch := range c.watchers[service] {
			select {
			case ch <- change:
			default: // 不阻塞
			}
		}
	}
}

func (c *ConfigCenter) Get(service, key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if cfg, ok := c.configs[service]; ok {
		val, exists := cfg[key]
		return val, exists
	}
	return "", false
}

func (c *ConfigCenter) Watch(service string) <-chan ConfigChange {
	c.mu.Lock()
	defer c.mu.Unlock()

	ch := make(chan ConfigChange, 10)
	c.watchers[service] = append(c.watchers[service], ch)
	return ch
}

// ========================================
// HTTP API for Registry
// ========================================

func startRegistryAPI(registry *Registry, addr string) {
	mux := http.NewServeMux()

	// POST /register
	mux.HandleFunc("POST /register", func(w http.ResponseWriter, r *http.Request) {
		var inst ServiceInstance
		if err := json.NewDecoder(r.Body).Decode(&inst); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		registry.Register(&inst)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "registered"})
	})

	// GET /discover?service=xxx
	mux.HandleFunc("GET /discover", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("service")
		instances := registry.Discover(name)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(instances)
	})

	// POST /heartbeat
	mux.HandleFunc("POST /heartbeat", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Name string `json:"name"`
			ID   string `json:"id"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		ok := registry.Heartbeat(req.Name, req.ID)
		json.NewEncoder(w).Encode(map[string]bool{"ok": ok})
	})

	// GET /services
	mux.HandleFunc("GET /services", func(w http.ResponseWriter, r *http.Request) {
		registry.mu.RLock()
		defer registry.mu.RUnlock()

		result := make(map[string]int)
		for name, instances := range registry.services {
			healthy := 0
			for _, inst := range instances {
				if inst.Health {
					healthy++
				}
			}
			result[name] = healthy
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	log.Printf("[Registry API] Listening on %s", addr)
	http.ListenAndServe(addr, mux)
}

func main() {
	fmt.Println("============================================")
	fmt.Println("  Service Discovery & Config Demo")
	fmt.Println("============================================")

	// 创建注册中心
	registry := NewRegistry()

	// 注册一些服务实例
	services := []*ServiceInstance{
		{ID: "user-1", Name: "user-service", Address: "10.0.0.1", Port: 8001,
			Metadata: map[string]string{"version": "1.0", "region": "cn-north"}},
		{ID: "user-2", Name: "user-service", Address: "10.0.0.2", Port: 8001,
			Metadata: map[string]string{"version": "1.0", "region": "cn-south"}},
		{ID: "order-1", Name: "order-service", Address: "10.0.0.3", Port: 8002,
			Metadata: map[string]string{"version": "2.0"}},
		{ID: "order-2", Name: "order-service", Address: "10.0.0.4", Port: 8002,
			Metadata: map[string]string{"version": "2.0"}},
		{ID: "order-3", Name: "order-service", Address: "10.0.0.5", Port: 8002,
			Metadata: map[string]string{"version": "2.1"}},
	}

	for _, svc := range services {
		registry.Register(svc)
	}

	// ---- 服务发现 ----
	fmt.Println("\n--- Service Discovery ---")
	userInstances := registry.Discover("user-service")
	fmt.Printf("Found %d healthy user-service instances:\n", len(userInstances))
	for _, inst := range userInstances {
		fmt.Printf("  %s -> %s (version: %s)\n",
			inst.ID, inst.Endpoint(), inst.Metadata["version"])
	}

	// ---- 负载均衡 ----
	fmt.Println("\n--- Load Balancing ---")
	client := NewServiceClient(registry, &RoundRobinLB{})

	for i := 0; i < 6; i++ {
		inst, err := client.Call(context.Background(), "order-service")
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		fmt.Printf("  Request %d -> %s (%s)\n", i+1, inst.ID, inst.Endpoint())
	}

	// ---- 配置中心 ----
	fmt.Println("\n--- Config Center ---")
	configCenter := NewConfigCenter()

	// 设置配置
	configCenter.Set("user-service", "db.host", "localhost")
	configCenter.Set("user-service", "db.port", "5432")
	configCenter.Set("user-service", "cache.ttl", "300")

	// 读取配置
	if host, ok := configCenter.Get("user-service", "db.host"); ok {
		fmt.Printf("  user-service/db.host = %s\n", host)
	}

	// 监听配置变更
	changes := configCenter.Watch("user-service")
	go func() {
		for change := range changes {
			fmt.Printf("  [Config Change] %s: '%s' -> '%s'\n",
				change.Key, change.OldValue, change.NewValue)
		}
	}()

	// 模拟配置变更
	configCenter.Set("user-service", "cache.ttl", "600")
	configCenter.Set("user-service", "db.host", "db.production.internal")
	time.Sleep(100 * time.Millisecond)

	// ---- 模拟心跳 ----
	fmt.Println("\n--- Heartbeat Simulation ---")
	// 模拟 user-1 停止心跳
	fmt.Println("  user-1 stopped sending heartbeat...")
	// user-2 继续发送心跳
	registry.Heartbeat("user-service", "user-2")
	fmt.Println("  user-2 heartbeat sent")

	// 启动注册中心 HTTP API
	fmt.Println("\n============================================")
	fmt.Println("  Registry API starting on :8082")
	fmt.Println("============================================")
	fmt.Println("  Endpoints:")
	fmt.Println("    GET  /services            - List all services")
	fmt.Println("    GET  /discover?service=xxx - Discover instances")
	fmt.Println("    POST /register            - Register instance")
	fmt.Println("    POST /heartbeat           - Send heartbeat")
	fmt.Println()
	fmt.Println("  Test:")
	fmt.Println(`    curl localhost:8082/services`)
	fmt.Println(`    curl "localhost:8082/discover?service=user-service"`)
	fmt.Println("============================================")

	startRegistryAPI(registry, ":8082")
}

// ========================================
// 练习:
// 1. 集成 etcd 或 Consul 作为真正的服务发现后端
// 2. 实现加权轮询负载均衡
// 3. 添加熔断器（Circuit Breaker）模式
// 4. 实现配置热更新（watch + reload）
// ========================================
