package ebui

import (
	"github.com/yanun0323/pkg/logs"
)

type RunOption func(*app)

// WithDebug setups option to show debug log
func WithDebug() RunOption {
	return func(a *app) {
		a.debug = true
		logs.SetDefaultLevel(logs.LevelDebug)
	}
}

// WithMemMonitor setups option to show memory analyzing
func WithMemMonitor() RunOption {
	return func(a *app) {
		a.memMonitor = true
	}
}
