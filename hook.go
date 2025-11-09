package configx

type HookPattern int

const (
	InitHook HookPattern = iota
	Debug
	Info
	Warn
	Error
	HookIndex
)

type HookContext struct {
	Message string
	Pattern HookPattern
}

type HookHandlerFunc func(ctx HookContext)

func (h HookHandlerFunc) Exec(ctx HookContext) {
	if h == nil {
		return
	}
	h(ctx)
}

type Hook struct {
	Handles [HookIndex]HookHandlerFunc
}

func NewHook() *Hook {
	return &Hook{}
}

func (hooks *Hook) SetHook(index HookPattern, h HookHandlerFunc) *Hook {
	SetHook(index, h)
	return hooks
}

func SetHook(index HookPattern, h HookHandlerFunc) *Hook {
	hooks := Default().hooks
	hooks.Handles[index] = h
	return hooks
}
