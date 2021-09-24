package rest

type Mock struct {
	GetFn    func(url string, data RequestData) (resp Response, err error)
	PostFn   func(url string, data RequestData) (resp Response, err error)
	PutFn    func(url string, data RequestData) (resp Response, err error)
	PatchFn  func(url string, data RequestData) (resp Response, err error)
	DeleteFn func(url string, data RequestData) (resp Response, err error)
}

func (m Mock) Get(url string, data RequestData) (resp Response, err error) {
	if m.GetFn != nil {
		return m.GetFn(url, data)
	}
	return Response{}, nil
}

func (m Mock) Post(url string, data RequestData) (resp Response, err error) {
	if m.PostFn != nil {
		return m.PostFn(url, data)
	}
	return Response{}, nil
}

func (m Mock) Put(url string, data RequestData) (resp Response, err error) {
	if m.PutFn != nil {
		return m.PutFn(url, data)
	}
	return Response{}, nil
}

func (m Mock) Patch(url string, data RequestData) (resp Response, err error) {
	if m.PatchFn != nil {
		return m.PatchFn(url, data)
	}
	return Response{}, nil
}

func (m Mock) Delete(url string, data RequestData) (resp Response, err error) {
	if m.DeleteFn != nil {
		return m.DeleteFn(url, data)
	}
	return Response{}, nil
}
