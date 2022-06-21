package log

import "context"

// Mock ...
type Mock struct {
	DebugFn func(ctx context.Context, title string, valueMaps ...Body)
	InfoFn  func(ctx context.Context, title string, valueMaps ...Body)
	WarnFn  func(ctx context.Context, title string, valueMaps ...Body)
	ErrorFn func(ctx context.Context, title string, valueMaps ...Body)
	FatalFn func(ctx context.Context, title string, valueMaps ...Body)
}

func (m Mock) Debug(ctx context.Context, title string, valueMaps ...Body) {
	if m.DebugFn != nil {
		m.DebugFn(ctx, title, valueMaps...)
	}
}

func (m Mock) Info(ctx context.Context, title string, valueMaps ...Body) {
	if m.InfoFn != nil {
		m.InfoFn(ctx, title, valueMaps...)
	}
}

func (m Mock) Warn(ctx context.Context, title string, valueMaps ...Body) {
	if m.WarnFn != nil {
		m.WarnFn(ctx, title, valueMaps...)
	}
}

func (m Mock) Error(ctx context.Context, title string, valueMaps ...Body) {
	if m.ErrorFn != nil {
		m.ErrorFn(ctx, title, valueMaps...)
	}
}

func (m Mock) Fatal(ctx context.Context, title string, valueMaps ...Body) {
	if m.FatalFn == nil {
		panic("calling log.Mock.Fatal, if you are testing it you should specify the behavior")
	}

	m.FatalFn(ctx, title, valueMaps...)
}
