package configx

import (
	"testing"
)

// TestSetOptionChaining tests that SetOption supports method chaining
func TestSetOptionChaining(t *testing.T) {
	type TestConfig struct {
		Name string `mapstructure:"name"`
	}

	manager := NewManager(TestConfig{})
	opts := NewOption()

	// Test method chaining
	result := manager.SetOption(opts)
	if result != manager {
		t.Error("SetOption should return the manager instance for chaining")
	}
}

// TestSetHookChaining tests that SetHook supports method chaining
func TestSetHookChaining(t *testing.T) {
	type TestConfig struct {
		Name string `mapstructure:"name"`
	}

	manager := NewManager(TestConfig{})

	// Test method chaining
	result := manager.SetHook(Info, func(ctx HookContext) {
		// Test hook
	})
	if result != manager {
		t.Error("SetHook should return the manager instance for chaining")
	}
}

// TestSetHookCompatibility tests that SetHook works with generic Manager
func TestSetHookCompatibility(t *testing.T) {
	type TestConfig struct {
		Name string `mapstructure:"name"`
	}

	manager := NewManager(TestConfig{})
	called := false

	// Set hook
	manager.SetHook(Info, func(ctx HookContext) {
		called = true
	})

	// Trigger hook
	manager.hooks.Handles[Info].Exec(HookContext{
		Message: "test",
		Pattern: Info,
	})

	if !called {
		t.Error("Hook should have been called")
	}
}

// TestMultipleHookLevels tests that all hook levels work correctly
func TestMultipleHookLevels(t *testing.T) {
	type TestConfig struct {
		Name string `mapstructure:"name"`
	}

	manager := NewManager(TestConfig{})
	callCounts := make(map[HookPattern]int)

	// Set hooks for all levels
	manager.
		SetHook(InitHook, func(ctx HookContext) { callCounts[InitHook]++ }).
		SetHook(Debug, func(ctx HookContext) { callCounts[Debug]++ }).
		SetHook(Info, func(ctx HookContext) { callCounts[Info]++ }).
		SetHook(Warn, func(ctx HookContext) { callCounts[Warn]++ }).
		SetHook(Error, func(ctx HookContext) { callCounts[Error]++ })

	// Trigger each hook
	manager.hooks.Handles[InitHook].Exec(HookContext{Pattern: InitHook})
	manager.hooks.Handles[Debug].Exec(HookContext{Pattern: Debug})
	manager.hooks.Handles[Info].Exec(HookContext{Pattern: Info})
	manager.hooks.Handles[Warn].Exec(HookContext{Pattern: Warn})
	manager.hooks.Handles[Error].Exec(HookContext{Pattern: Error})

	// Verify all hooks were called
	for pattern, count := range callCounts {
		if count != 1 {
			t.Errorf("Hook %d should have been called once, got %d", pattern, count)
		}
	}
}
