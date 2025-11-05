package config

type LoggerHandlerFunc func(msg string)

type Logger struct {
	Debug LoggerHandlerFunc
	Info  LoggerHandlerFunc
	Error LoggerHandlerFunc
}

func (h LoggerHandlerFunc) Exec(msg string) {
	if h == nil {
		//fmt.Println("[config] logger is nil")
		return
	}

	h(msg)
}
