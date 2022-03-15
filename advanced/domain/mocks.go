package domain

import (
	"context"
)

// LogProviderMock ...
type LogProviderMock struct {
	DebugFn func(ctx context.Context, title string, valueMaps ...LogBody)
	InfoFn  func(ctx context.Context, title string, valueMaps ...LogBody)
	WarnFn  func(ctx context.Context, title string, valueMaps ...LogBody)
	ErrorFn func(ctx context.Context, title string, valueMaps ...LogBody)
	FatalFn func(ctx context.Context, title string, valueMaps ...LogBody)
}

func (m LogProviderMock) Debug(ctx context.Context, title string, valueMaps ...LogBody) {
	if m.DebugFn != nil {
		m.DebugFn(ctx, title, valueMaps...)
	}
}

func (m LogProviderMock) Info(ctx context.Context, title string, valueMaps ...LogBody) {
	if m.InfoFn != nil {
		m.InfoFn(ctx, title, valueMaps...)
	}
}

func (m LogProviderMock) Warn(ctx context.Context, title string, valueMaps ...LogBody) {
	if m.WarnFn != nil {
		m.WarnFn(ctx, title, valueMaps...)
	}
}

func (m LogProviderMock) Error(ctx context.Context, title string, valueMaps ...LogBody) {
	if m.ErrorFn != nil {
		m.ErrorFn(ctx, title, valueMaps...)
	}
}

func (m LogProviderMock) Fatal(ctx context.Context, title string, valueMaps ...LogBody) {
	if m.FatalFn == nil {
		panic("calling domain.LogProviderMock.Fatal, if you are testing it you should specify the behavior")
	}

	m.FatalFn(ctx, title, valueMaps...)
}

type RestProviderMock struct {
	GetFn    func(ctx context.Context, url string, data RequestData) (resp Response, err error)
	PostFn   func(ctx context.Context, url string, data RequestData) (resp Response, err error)
	PutFn    func(ctx context.Context, url string, data RequestData) (resp Response, err error)
	PatchFn  func(ctx context.Context, url string, data RequestData) (resp Response, err error)
	DeleteFn func(ctx context.Context, url string, data RequestData) (resp Response, err error)
}

func (m RestProviderMock) Get(ctx context.Context, url string, data RequestData) (resp Response, err error) {
	if m.GetFn != nil {
		return m.GetFn(ctx, url, data)
	}
	return Response{}, nil
}

func (m RestProviderMock) Post(ctx context.Context, url string, data RequestData) (resp Response, err error) {
	if m.PostFn != nil {
		return m.PostFn(ctx, url, data)
	}
	return Response{}, nil
}

func (m RestProviderMock) Put(ctx context.Context, url string, data RequestData) (resp Response, err error) {
	if m.PutFn != nil {
		return m.PutFn(ctx, url, data)
	}
	return Response{}, nil
}

func (m RestProviderMock) Patch(ctx context.Context, url string, data RequestData) (resp Response, err error) {
	if m.PatchFn != nil {
		return m.PatchFn(ctx, url, data)
	}
	return Response{}, nil
}

func (m RestProviderMock) Delete(ctx context.Context, url string, data RequestData) (resp Response, err error) {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, url, data)
	}
	return Response{}, nil
}
