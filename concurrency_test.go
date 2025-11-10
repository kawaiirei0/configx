package configx

import (
	"sync"
	"testing"
	"time"
)

// TestConcurrentGetConfig 测试并发读取配置
func TestConcurrentGetConfig(t *testing.T) {
	type TestConfig struct {
		Value string `mapstructure:"value"`
	}

	manager := NewManager(TestConfig{})
	opts := NewOption()
	opts.Filename.Set("test_config.yaml")
	opts.Filepath.Set("./testdata")
	manager.SetOption(opts)

	if err := manager.LoadConfig(); err != nil {
		t.Skipf("跳过测试，配置文件不存在: %v", err)
		return
	}

	// 并发读取配置
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = manager.GetConfig()
		}()
	}
	wg.Wait()
}

// TestConcurrentSetHook 测试并发设置钩子
func TestConcurrentSetHook(t *testing.T) {
	type TestConfig struct {
		Value string `mapstructure:"value"`
	}

	manager := NewManager(TestConfig{})

	// 并发设置钩子
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			manager.SetHook(Info, func(ctx HookContext) {
				// 钩子处理
			})
		}(i)
	}
	wg.Wait()
}

// TestConcurrentSetOption 测试并发设置选项
func TestConcurrentSetOption(t *testing.T) {
	type TestConfig struct {
		Value string `mapstructure:"value"`
	}

	manager := NewManager(TestConfig{})

	// 并发设置选项
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			opts := NewOption()
			manager.SetOption(opts)
		}()
	}
	wg.Wait()
}

// TestDebounceRaceCondition 测试防抖机制的竞态条件
func TestDebounceRaceCondition(t *testing.T) {
	type TestConfig struct {
		Value string `mapstructure:"value"`
	}

	manager := NewManager(TestConfig{})
	manager.debounceDur = 100 * time.Millisecond

	// 模拟并发访问 lastChangeNano (现在使用 atomic 操作，应该是线程安全的)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			now := time.Now()
			// 使用 atomic 操作访问
			lastNano := manager.lastChangeNano.Load()
			lastTime := time.Unix(0, lastNano)
			_ = now.Sub(lastTime) < manager.debounceDur
			manager.lastChangeNano.Store(now.UnixNano())
		}()
	}
	wg.Wait()
}

// TestConcurrentHookExecution 测试并发执行钩子
func TestConcurrentHookExecution(t *testing.T) {
	type TestConfig struct {
		Value string `mapstructure:"value"`
	}

	manager := NewManager(TestConfig{})
	
	counter := 0
	var mu sync.Mutex

	// 设置钩子
	manager.SetHook(Info, func(ctx HookContext) {
		mu.Lock()
		counter++
		mu.Unlock()
	})

	// 并发执行钩子和设置钩子
	var wg sync.WaitGroup
	
	// 并发执行钩子（使用线程安全的 executeHook 方法）
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			manager.executeHook(Info, HookContext{
				Message: "test",
				Pattern: Info,
			})
		}()
	}

	// 同时并发修改钩子
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			manager.SetHook(Info, func(ctx HookContext) {
				mu.Lock()
				counter++
				mu.Unlock()
			})
		}()
	}

	wg.Wait()
}

// TestConcurrentConfigReload 测试并发配置重载
func TestConcurrentConfigReload(t *testing.T) {
	type TestConfig struct {
		Value string `mapstructure:"value"`
	}

	manager := NewManager(TestConfig{})
	opts := NewOption()
	opts.Filename.Set("test_config.yaml")
	opts.Filepath.Set("./testdata")
	manager.SetOption(opts)

	if err := manager.LoadConfig(); err != nil {
		t.Skipf("跳过测试，配置文件不存在: %v", err)
		return
	}

	// 并发重载和读取配置
	var wg sync.WaitGroup
	
	// 并发读取
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = manager.GetConfig()
		}()
	}

	// 并发重载
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = manager.LoadConfig()
		}()
	}

	wg.Wait()
}
