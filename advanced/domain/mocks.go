package domain

import (
	"context"
)

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
